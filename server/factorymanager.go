package server

import (
	"laatoo/sdk/core"
)

type FactoryManager interface {
	core.ServerElement
	GetFactory(ctx core.ServerContext, factoryName string) (Factory, error)
}
