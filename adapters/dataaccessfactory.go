package adapters

import (
	"laatoo/core/services"
	"laatoo/sdk/config"
	"laatoo/sdk/core"
	"laatoo/sdk/data"
	"laatoo/sdk/errors"
)

const (
	CONF_DATA_SVC_NAME = "data_svc"
)

type DataAccessFactory struct {
	Conf      config.Config
	facMethod FactoryMethod
	DataStore data.DataService
}

func NewDataAccessFactory(facMethod FactoryMethod) *DataAccessFactory {
	return &DataAccessFactory{facMethod: facMethod}
}

//Create the services configured for factory.
func (df *DataAccessFactory) CreateService(ctx core.ServerContext, name string, method string, conf config.Config) (core.Service, error) {
	df.Conf = conf
	svcFunc, err := df.facMethod(ctx, name, method, conf)
	if err != nil {
		return nil, err
	}
	return services.NewService(ctx, name, svcFunc, conf), nil
}

//The services start serving when this method is called
func (df *DataAccessFactory) StartServices(ctx core.ServerContext) error {
	if df.Conf == nil {
		return errors.ThrowError(ctx, errors.CORE_ERROR_NOT_IMPLEMENTED, "Info", "Default data factory must set the config in create service")
	}
	datasvcName, ok := df.Conf.GetString(CONF_DATA_SVC_NAME)
	if !ok {
		return errors.ThrowError(ctx, errors.CORE_ERROR_MISSING_CONF, "Configuration", CONF_DATA_SVC_NAME)
	}
	dataSvc, err := ctx.GetService(datasvcName)
	if err != nil {
		return errors.RethrowError(ctx, errors.CORE_ERROR_MISSING_SERVICE, err, "Name", datasvcName)
	}
	df.DataStore = dataSvc.(data.DataService)
	return nil
}
