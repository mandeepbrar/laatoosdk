package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ActivityManager owns the activities available at one server level: the ones loaded from
// src/server/registry/activities/, and the scripts the script manager registers. An activity is a
// SERVICE-EQUIVALENT: it is exposed as the service activity.<name> and authorized like one.
type ActivityManager interface {
	core.ServerElement

	// RegisterActivityProvider registers the provider that serves one activity type -- a type the
	// server does not build in, such as core.ActivityTypeFlow.
	//
	// CALL IT FROM A FACTORY'S Initialize. Each activity of the type binds to its provider when
	// the activity's service starts (see components.ActivityProvider), which is after every
	// factory initializes and before any request.
	//
	// ONE PROVIDER PER TYPE ACROSS A BRANCH OF NAMESPACES: a type already served at this level or
	// any level above is refused, naming both. A built-in type (service, script, manual) cannot be
	// registered.
	//
	// There is no executor registration: a Go function is an ACTION, registered with
	// ActionManager.RegisterAction, and a workflow step names it directly.
	RegisterActivityProvider(ctx core.ServerContext, activityType core.ActivityType, provider components.ActivityProvider) error

	// ExecuteActivity runs an ACTIVITY by name and returns its result.
	//
	// IT NEVER RUNS AN ACTION. A workflow step names an action with `action:`, which reaches
	// ActionManager.ExecuteAction; an `activity:` step reaches this method. Neither falls back to
	// the other, so an action of the same name is not consulted when no activity has the name.
	//
	// IT DISPATCHES THROUGH THE SERVICE LAYER, NOT THROUGH THE EXECUTOR MAP: the name is resolved
	// as the service "activity."+activityName and invoked with params under the single argument
	// "activityparams" (activitymanager_impl.go:254-262). An activity name with no corresponding
	// service fails as a missing service, whatever the executor registry holds.
	//
	// Because it goes through a service, the activity is subject to ordinary service
	// authorization, and its result comes back on the request context: a success response's Data
	// is returned, a failure response's Error is wrapped and returned, and a service that set NO
	// response at all yields (nil, nil).
	//
	// Streaming activities are drained through the response handler before the return value is
	// read, and a drain error is LOGGED AND SWALLOWED rather than returned
	// (activitymanager_impl.go:268-275) — a partially delivered stream still looks successful.
	ExecuteActivity(ctx core.RequestContext, activityName string, params utils.StringMap) (interface{}, error)
	// GetActivityDefinition returns the ActivityDefinition for a registered activity by name,
	// or nil if the activity is not found. Used by workflow engines to resolve activity
	// metadata (e.g. ActivityType) without requiring it to be repeated in the workflow DSL.
	GetActivityDefinition(ctx core.ServerContext, activityName string) *core.ActivityDefinition
	// SetDefaultStreamingHandler registers the response handler used to drain streaming
	// ResponseStream channels after an activity completes with ctx.IsStreaming() == true.
	SetDefaultStreamingHandler(ctx core.ServerContext, handler core.ResponseHandler) error
}
