package elements

import (
	"laatoo/sdk/server/core"
)

type FactoryManager interface {
	core.ServerElement
	GetFactory(ctx core.ServerContext, factoryName string) (Factory, error)
}
