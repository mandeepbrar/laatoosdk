package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type ActivityManager interface {
	core.ServerElement
	GetActivity(ctx core.ServerContext, alias string) (components.Activity, error)
	RegisterActivity(ctx core.ServerContext, alias string, act components.Activity) error
	InvokeActivity(ctx core.RequestContext, activity string, params utils.StringMap) (interface{}, error)
}
