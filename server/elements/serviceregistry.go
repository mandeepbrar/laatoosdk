package elements

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
)

type ServiceRegistry interface {
	RegisterService(ctx core.ServerContext, serviceAlias string, svc Service, conf config.Config) error
	GetService(ctx core.ServerContext, serviceName string) (Service, config.Config, error)
}
