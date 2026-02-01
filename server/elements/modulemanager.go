package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type ModuleManager interface {
	core.ServerElement
	List(ctx core.ServerContext) utils.StringsMap
	Describe(ctx core.ServerContext, mod string) (utils.StringMap, error)
	ChangeLogger(ctx core.ServerContext, mod string, logLevel string, logFormat string) error
	GetModule(ctx core.ServerContext, modName string) (core.Module, error)
}
