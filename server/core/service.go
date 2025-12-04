package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/utils"
)

type Service interface {
	ConfigurableObject
	Metadata() ServiceInfo
	Describe(ServerContext) error
	Initialize(ctx ServerContext, conf config.Config) error
	Start(ctx ServerContext) error
	Stop(ctx ServerContext) error
	Unload(ctx ServerContext) error
	AddParams(ServerContext, map[string]datatypes.DataType, bool) error
	//	AddStringParams(ctx ServerContext, params utils.StringsMap, defaultValues []string)
	AddStringParam(ctx ServerContext, name string, desc string)
	AddCustomObjectParam(ctx ServerContext, name string, desc string, customObjectType string, collection, required, stream bool) error
	AddParam(ctx ServerContext, name string, desc string, datatype datatypes.DataType, collection, required, stream bool) error
	AddParamWithType(ctx ServerContext, name string, desc string, datatype datatypes.DataType) error
	AddOptionalParamWithType(ctx ServerContext, name string, desc string, datatype datatypes.DataType) error
	//AddCollectionParams(ctx ServerContext, params map[string]datatypes.DataType) error
	//	SetRequestType(ctx ServerContext, datatype string, collection bool, stream bool)
	//	SetResponseType(ctx ServerContext, stream bool)
	SetDescription(ServerContext, string)
	SetComponent(ServerContext, bool)
	//ConfigureService(ctx ServerContext, requestType string, collection bool, stream bool, params []string, config []string, description string)
	//ConfigureService(ctx ServerContext, params []string, config []string, description string)
}

type UserInvokableService interface {
	Service
	Invoke(RequestContext) error
}

type Param interface {
	GetName() string
	GetDescription() string
	IsCollection() bool
	IsStream() bool
	IsRequired() bool
	GetDataType() datatypes.DataType
	GetValue() interface{}
}

type ServiceFunc func(ctx RequestContext) error

type Request interface {
	//GetBody() interface{}
	GetParam(RequestContext, string) (Param, bool)
	GetParams(RequestContext) map[string]Param
	GetParamValue(RequestContext, string) (interface{}, bool)
	GetIntParam(RequestContext, string) (int, bool)
	GetStringParam(RequestContext, string) (string, bool)
	GetStringMapParam(RequestContext, string) (utils.StringMap, bool)
	GetStringsMapParam(RequestContext, string) (utils.StringsMap, bool)
}

type Response struct {
	Status   int
	Data     interface{}
	MetaInfo map[string]interface{}
	Error    error
	Return   bool
}
