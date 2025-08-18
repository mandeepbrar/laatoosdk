package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type FactoryManager interface {
	core.ServerElement
	GetFactory(ctx core.ServerContext, factoryName string) (Factory, error)
	List(ctx core.ServerContext) utils.StringsMap
	ChangeLogger(ctx core.ServerContext, chanName string, logLevel string, logFormat string) error
	Describe(ctx core.ServerContext, chanName string) (utils.StringMap, error)
}
