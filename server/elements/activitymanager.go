package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ActivityManager owns the activities available to workflows at one server level — both the ones
// loaded from src/server/registry/activities/ and the ones plugins register programmatically.
type ActivityManager interface {
	core.ServerElement

	// RegisterActivity binds a Go function as the executor for an activity name.
	//
	// CALL IT NO LATER THAN YOUR MODULE'S Initialize. The binding is consumed when the activity's
	// own service STARTS: ExecutorActivity.Start looks the executor up by activity alias, and a
	// non-manual activity with no executor by then fails startup with "Activity Executor not
	// found" (laatooserver/src/core/activities.go:93-104). Registering after that point leaves the
	// activity bound to nothing — Start has already run and is not repeated. Manual (HITL)
	// activities are the deliberate exception: they are allowed to carry no executor.
	//
	// A DUPLICATE NAME IS REFUSED, NOT REPLACED — a second registration for the same name logs a
	// warning and returns a Bad Conf error, leaving the first executor in place
	// (laatooserver/src/core/activitymanager_impl.go:224-244).
	//
	// Registering an executor does NOT create an ActivityDefinition; the code that would have done
	// so is commented out (activitymanager_impl.go:231-238). GetActivityDefinition therefore
	// returns nil for an activity that exists only as a registered executor.
	RegisterActivity(ctx core.ServerContext, activityName string, executor core.ActivityExecutor) error

	// ExecuteActivity runs an activity by name and returns its result.
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
