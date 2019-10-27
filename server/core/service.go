package core

import "laatoo/sdk/common/config"

type Service interface {
	ConfigurableObject
	Describe(ServerContext) error
	Initialize(ctx ServerContext, conf config.Config) error
	Start(ctx ServerContext) error
	Stop(ctx ServerContext) error
	Unload(ctx ServerContext) error
	Invoke(RequestContext) error
	AddParams(ServerContext, map[string]string, bool) error
	AddStringParams(ctx ServerContext, names []string, defaultValues []string)
	AddStringParam(ctx ServerContext, name string)
	AddParam(ctx ServerContext, name string, datatype string, collection, required, stream bool) error
	AddParamWithType(ctx ServerContext, name string, datatype string) error
	AddOptionalParamWithType(ctx ServerContext, name string, datatype string) error
	AddCollectionParams(ServerContext, map[string]string) error
	//	SetRequestType(ctx ServerContext, datatype string, collection bool, stream bool)
	//	SetResponseType(ctx ServerContext, stream bool)
	SetDescription(ServerContext, string)
	SetComponent(ServerContext, bool)
	//ConfigureService(ctx ServerContext, requestType string, collection bool, stream bool, params []string, config []string, description string)
	ConfigureService(ctx ServerContext, params []string, config []string, description string)
}

type Param interface {
	GetName() string
	IsCollection() bool
	IsStream() bool
	IsRequired() bool
	GetDataType() string
	GetValue() interface{}
}

type ServiceFunc func(ctx RequestContext) error

type Request interface {
	//GetBody() interface{}
	GetParam(string) (Param, bool)
	GetParams() map[string]Param
	GetParamValue(string) (interface{}, bool)
	GetIntParam(string) (int, bool)
	GetStringParam(string) (string, bool)
	GetStringMapParam(string) (map[string]interface{}, bool)
	GetStringsMapParam(string) (map[string]string, bool)
}

type Response struct {
	Status   int
	Data     interface{}
	MetaInfo map[string]interface{}
	Error    error
	Return   bool
}
