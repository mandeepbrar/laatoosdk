package core

import (
	"laatoo/sdk/common/config"
)

//service method for doing various tasks
func NewFactory(creator ObjectCreator) ObjectCreator {
	return func() interface{} {
		return &FactoryImpl{creator: creator}
	}
}

type FactoryImpl struct {
	ServiceFactory
	creator ObjectCreator
}

//Create the services configured for factory.
func (fac *FactoryImpl) CreateService(ctx ServerContext, name string, method string, conf config.Config) (Service, error) {
	obj := fac.creator()
	svc, _ := obj.(Service)
	return svc, nil
}
