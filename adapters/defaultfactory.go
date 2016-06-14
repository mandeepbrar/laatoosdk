package adapters

import (
	"laatoo/core/services"
	"laatoo/sdk/config"
	"laatoo/sdk/core"
)

type FactoryMethod func(ctx core.ServerContext, name string, method string, conf config.Config) (core.ServiceFunc, error)

type DefaultFactory struct {
	Conf      config.Config
	facMethod FactoryMethod
}

func NewDefaultFactory(facMethod FactoryMethod) *DefaultFactory {
	return &DefaultFactory{facMethod: facMethod}
}

//Create the services configured for factory.
func (df *DefaultFactory) CreateService(ctx core.ServerContext, name string, method string, conf config.Config) (core.Service, error) {
	df.Conf = conf
	svcFunc, err := df.facMethod(ctx, name, method, conf)
	if err != nil {
		return nil, err
	}
	return services.NewService(ctx, name, svcFunc, conf), nil
}

//The services start serving when this method is called
func (ds *DefaultFactory) StartServices(ctx core.ServerContext) error {
	return nil
}
