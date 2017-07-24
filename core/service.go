package core

import "laatoo/sdk/config"

type Service interface {
	Initialize(ctx ServerContext) error
	Info() ServiceInfo
	Start(ctx ServerContext) error
	Invoke(RequestContext, Request) (*Response, error)
	AddParams(map[string]string)
	AddStringParams(names []string, defaultValues []string)
	AddStringParam(name string)
	AddParam(name string, datatype string, collection bool)
	AddCollectionParams(map[string]string)
	AddStringConfigurations(names []string, defaultValues []string)
	AddStringConfiguration(name string)
	AddConfigurations(map[string]string)
	AddOptionalConfigurations(map[string]string, map[string]interface{})
	SetRequestType(datatype string, collection bool, stream bool)
	SetResponseType(stream bool)
	GetConfiguration(string) (interface{}, bool)
	GetStringConfiguration(string) (string, bool)
	GetBoolConfiguration(string) (bool, bool)
	GetMapConfiguration(string) (config.Config, bool)

	SetDescription(string)
	SetComponent(bool)
	ConfigureService(requestType string, collection bool, stream bool, params map[string]string, config map[string]string, description string)
}

type ServiceInfo interface {
	GetRequestInfo() RequestInfo
	GetResponseInfo() ResponseInfo
	GetConfigurations() map[string]interface{}
	GetDescription() string
	IsComponent() bool
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

type ServiceFunc func(ctx RequestContext, request Request) (*Response, error)

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
