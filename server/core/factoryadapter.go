package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/ctx"
)

// service method for doing various tasks
func NewFactory(creator ObjectCreator) ObjectCreator {
	return func(cx ctx.Context) interface{} {
		return &FactoryImpl{creator: creator}
	}
}

type FactoryImpl struct {
	ServiceFactory
	creator ObjectCreator
}

// Create the services configured for factory.
func (fac *FactoryImpl) CreateService(ctx ServerContext, name string, method string, conf config.Config) (Service, error) {
	obj := fac.creator(ctx)
	svc, _ := obj.(Service)
	return svc, nil
}
