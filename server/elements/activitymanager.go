package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type ActivityManager interface {
	core.ServerElement
	RegisterActivity(ctx core.ServerContext, activityName string, executor core.ActivityExecutor) error
	ExecuteActivity(ctx core.RequestContext, activityName string, params utils.StringMap) (interface{}, error)
}
