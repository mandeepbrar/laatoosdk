package elements

import (
	"laatoo/sdk/server/core"
)

type ServiceManager interface {
	core.ServerElement
	GetService(ctx core.ServerContext, alias string) (Service, error)
	GetServiceContext(ctx core.ServerContext, alias string) (core.ServerContext, error)
}
