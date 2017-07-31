package core

import "laatoo/sdk/config"

type Service interface {
	Initialize(ctx ServerContext) error
	Info() ServiceInfo
	Start(ctx ServerContext) error
	Invoke(RequestContext) error
	AddParams(ServerContext, map[string]string)
	AddStringParams(ctx ServerContext, names []string, defaultValues []string)
	AddStringParam(ctx ServerContext, name string)
	AddParam(ctx ServerContext, name string, datatype string, collection bool)
	AddCollectionParams(ServerContext, map[string]string)
	AddStringConfigurations(ctx ServerContext, names []string, defaultValues []string)
	AddStringConfiguration(ctx ServerContext, name string)
	AddConfigurations(ServerContext, map[string]string)
	AddOptionalConfigurations(ServerContext, map[string]string, map[string]interface{})
	SetRequestType(ctx ServerContext, datatype string, collection bool, stream bool)
	SetResponseType(ctx ServerContext, stream bool)
	GetConfiguration(ServerContext, string) (interface{}, bool)
	GetStringConfiguration(ServerContext, string) (string, bool)
	GetBoolConfiguration(ServerContext, string) (bool, bool)
	GetMapConfiguration(ServerContext, string) (config.Config, bool)
	InjectServices(ServerContext, map[string]string)

	SetDescription(ServerContext, string)
	SetComponent(ServerContext, bool)
	ConfigureService(ctx ServerContext, requestType string, collection bool, stream bool, params map[string]string, config map[string]string, description string)
}

type ServiceInfo interface {
	GetRequestInfo() RequestInfo
	GetResponseInfo() ResponseInfo
	GetConfigurations() map[string]interface{}
	GetDescription() string
	IsComponent() bool
	GetRequiredServices() map[string]string
}

type RequestInfo interface {
	GetDataType() string
	IsCollection() bool
	IsStream() bool
	GetParams() map[string]Param
}

type ResponseInfo interface {
	IsStream() bool
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
