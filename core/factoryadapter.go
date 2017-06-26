package core

import (
	"laatoo/sdk/config"
)

//service method for doing various tasks
func NewFactory(creator ObjectCreator) ObjectCreator {
	return func() interface{} {
		return &FactoryImpl{creator: creator}
	}
}

type FactoryImpl struct {
	creator ObjectCreator
}

func (fac *FactoryImpl) Initialize(ctx ServerContext, conf config.Config) error {
	return nil
}

//Create the services configured for factory.
func (fac *FactoryImpl) CreateService(ctx ServerContext, name string, method string, conf config.Config) (Service, error) {
	obj := fac.creator()
	svc, _ := obj.(Service)
	return svc, nil
}

//The services start serving when this method is called
func (fac *FactoryImpl) Start(ctx ServerContext) error {
	return nil
}
