package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type ActivityManager interface {
	core.ServerElement
	RegisterActivity(ctx core.ServerContext, activityName string, executor core.ActivityExecutor) error
	ExecuteActivity(ctx core.RequestContext, activityName string, params utils.StringMap) (interface{}, error)
	// GetActivityDefinition returns the ActivityDefinition for a registered activity by name,
	// or nil if the activity is not found. Used by workflow engines to resolve activity
	// metadata (e.g. ActivityType) without requiring it to be repeated in the workflow DSL.
	GetActivityDefinition(ctx core.ServerContext, activityName string) *core.ActivityDefinition
	// SetDefaultStreamingHandler registers the response handler used to drain streaming
	// ResponseStream channels after an activity completes with ctx.IsStreaming() == true.
	SetDefaultStreamingHandler(ctx core.ServerContext, handler core.ResponseHandler) error
}
