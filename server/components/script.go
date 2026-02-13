package components

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type Script interface {
	GetName(ctx core.ServerContext) string
	GetParams(ctx core.ServerContext) map[string]core.Param
	GetModule() core.Module
	GetScriptManager() ScriptManager
}

type ScriptManager interface {
	Load(ctx core.ServerContext, dir string) error
	InvokeScript(ctx core.RequestContext, act Script, args utils.StringMap) error
}
