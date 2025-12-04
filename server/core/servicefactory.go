package core

import (
	"laatoo.io/sdk/config"
)

type ServiceFactoryProvider func(ctx ServerContext, config config.Config) (ServiceFactory, error)

// Service interface that needs to be implemented by any service of a system
type ServiceFactory interface {
	ConfigurableObject
	Metadata() ServiceFactoryInfo
	Describe(ServerContext) error
	Initialize(ctx ServerContext, conf config.Config) error
	Start(ctx ServerContext) error
	Stop(ctx ServerContext) error
	Unload(ctx ServerContext) error
	//Create the services configured for factory.
	CreateService(ctx ServerContext, name string, method string, conf config.Config) (Service, error)
}
