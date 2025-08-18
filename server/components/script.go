package components

import (
	"laatoo.io/sdk/server/core"
)

type Script interface {
	GetName(ctx core.ServerContext) string
	GetParams(ctx core.ServerContext) map[string]core.Param
}

type ScriptManager interface {
	Load(ctx core.ServerContext, dir string) error
	GetScript(ctx core.ServerContext, alias string) (Script, error)
	RegisterScript(ctx core.ServerContext, alias string, act Script) error
	InvokeScript(ctx core.RequestContext, act Script, args ...interface{}) (interface{}, error)
}
