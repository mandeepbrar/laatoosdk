package core

import (
	"laatoo/sdk/config"
)

type ServiceParam struct {
	Name       string
	Value      interface{}
	Paramtype  string
	Collection bool
}

type ServiceParamsMap map[string]*ServiceParam

type RequestInfo struct {
	DataType     string
	IsCollection bool
	Streaming    bool
	Params       ServiceParamsMap
}

type ResponseInfo struct {
	DataType     string
	IsCollection bool
	Streaming    bool
	Params       ServiceParamsMap
}

type ServiceInfo struct {
	Request  RequestInfo
	Response ResponseInfo
}

type ServiceRequest interface {
	GetBody() interface{}
	SetBody(interface{})
	GetParams() ServiceParamsMap
	SetParams(ServiceParamsMap)
	GetParam(string) (*ServiceParam, bool)
	AddParam(name string, val interface{}, typ string, collection bool)
}

type ServiceResponse struct {
	Status int
	Data   interface{}
	Info   ServiceParamsMap
	Return bool
}

type ServiceFunc func(ctx RequestContext, request ServiceRequest) (*ServiceResponse, error)

type Service interface {
	Info() *ServiceInfo
	Initialize(ctx ServerContext, conf config.Config) error
	Start(ctx ServerContext) error
	Invoke(RequestContext, ServiceRequest) (*ServiceResponse, error)
}

func (paramsMap ServiceParamsMap) AddParam(name string, val interface{}, typ string, collection bool) {
	paramsMap[name] = &ServiceParam{name, val, typ, collection}
}
