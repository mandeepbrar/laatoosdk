package core

import (
	"laatoo/sdk/config"
)

//service method for doing various tasks
func NewFactory(creator ObjectCreator) ServiceFactory {
	return &factoryImpl{creator: creator}
}

type factoryImpl struct {
	creator ObjectCreator
}

func (fac *factoryImpl) Initialize(ctx ServerContext, conf config.Config) error {
	return nil
}

//Create the services configured for factory.
func (fac *factoryImpl) CreateService(ctx ServerContext, name string, method string, conf config.Config) (Service, error) {
	obj := fac.creator()
	svc, _ := obj.(Service)
	return svc, nil
}

//The services start serving when this method is called
func (fac *factoryImpl) Start(ctx ServerContext) error {
	return nil
}
