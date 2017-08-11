package core

import "laatoo/sdk/config"

type Service interface {
	ConfigurableObject
	Describe(ServerContext)
	Initialize(ctx ServerContext, conf config.Config) error
	Start(ctx ServerContext) error
	Invoke(RequestContext) error
	AddParams(ServerContext, map[string]string)
	AddStringParams(ctx ServerContext, names []string, defaultValues []string)
	AddStringParam(ctx ServerContext, name string)
	AddParam(ctx ServerContext, name string, datatype string, collection bool)
	AddCollectionParams(ServerContext, map[string]string)
	SetRequestType(ctx ServerContext, datatype string, collection bool, stream bool)
	SetResponseType(ctx ServerContext, stream bool)
	InjectServices(ServerContext, map[string]string)
	SetDescription(ServerContext, string)
	SetComponent(ServerContext, bool)
	ConfigureService(ctx ServerContext, requestType string, collection bool, stream bool, params []string, config []string, description string)
}

type Param interface {
	GetName() string
	IsCollection() bool
	GetDataType() string
	GetValue() interface{}
}

type ServiceFunc func(ctx RequestContext) error

type Request interface {
	GetBody() interface{}
	GetParam(string) (Param, bool)
	GetParams() map[string]Param
	GetIntParam(string) (int, bool)
	GetStringParam(string) (string, bool)
	GetStringMapValue(string) (map[string]interface{}, bool)
}

type Response struct {
	Status int
	Data   interface{}
	Info   map[string]interface{}
	Return bool
}
