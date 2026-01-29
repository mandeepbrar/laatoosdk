package components

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type Activity interface {
	GetName(ctx core.ServerContext) string
	GetParams(ctx core.ServerContext) map[string]core.Param
	GetModule() core.Module
}

type ActivityManager interface {
	Load(ctx core.ServerContext, dir string) error
	GetActivity(ctx core.ServerContext, alias string) (Activity, error)
	RegisterActivity(ctx core.ServerContext, alias string, act Activity) error
	InvokeActivity(ctx core.RequestContext, act Activity, args utils.StringMap) (interface{}, error)
}
