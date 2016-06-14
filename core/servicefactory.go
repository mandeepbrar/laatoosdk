package core

import (
	"laatoo/sdk/config"
)

type ServiceFactoryProvider func(ctx ServerContext, config config.Config) (ServiceFactory, error)

//Service interface that needs to be implemented by any service of a system
type ServiceFactory interface {
	Initialize(ctx ServerContext, conf config.Config) error
	Start(ctx ServerContext) error
	//Create the services configured for factory.
	CreateService(ctx ServerContext, name string, method string) (Service, error)
}
