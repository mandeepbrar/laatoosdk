package core

import (
	"laatoo.io/sdk/config"
)

type ServiceFactoryProvider func(ctx ServerContext, config config.Config) (ServiceFactory, error)

// Service interface that needs to be implemented by any service of a system
// ServiceFactory interface that needs to be implemented by any service factory.
// Service factories are responsible for creating service instances.
type ServiceFactory interface {
	ConfigurableObject
	// Metadata returns metadata about the service factory.
	Metadata() ServiceFactoryInfo
	// Describe describes the service factory to the server context.
	Describe(ServerContext) error
	// Initialize initializes the service factory with context and configuration.
	Initialize(ctx ServerContext, conf config.Config) error
	// Start starts the service factory.
	Start(ctx ServerContext) error
	// Stop stops the service factory.
	Stop(ctx ServerContext) error
	// Unload unloads the service factory.
	Unload(ctx ServerContext) error
	// CreateService creates a new service instance configured for the factory.
	CreateService(ctx ServerContext, name string, method string, conf config.Config) (Service, error)
	ServerElement() ServerElement
}
