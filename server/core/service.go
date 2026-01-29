package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/utils"
)

// Service interface that needs to be implemented by any service of a system.
// Services are the main units of logic in Laatoo.
type Service interface {
	ConfigurableObject
	// Metadata returns metadata about the service.
	Metadata() ServiceInfo
	// Describe describes the service to the server context.
	Describe(ServerContext) error
	// Initialize initializes the service with context and configuration.
	Initialize(ctx ServerContext, conf config.Config) error
	// Start starts the service.
	Start(ctx ServerContext) error
	// Stop stops the service.
	Stop(ctx ServerContext) error
	// Unload unloads the service.
	Unload(ctx ServerContext) error
	//Parameters for the service
	RequestParameters(ctx ServerContext) map[string]Param
	// AddParams adds parameters to the service definition.
	AddParams(ServerContext, map[string]datatypes.DataType, bool) error
	//	AddStringParams(ctx ServerContext, params utils.StringsMap, defaultValues []string)
	// AddStringParam adds a string parameter to the service definition.
	AddStringParam(ctx ServerContext, name string, desc string)
	// AddCustomObjectParam adds a custom object parameter to the service definition.
	AddCustomObjectParam(ctx ServerContext, name string, desc string, customObjectType string, collection, required, stream bool) error
	// AddParam adds a parameter of a specific data type to the service definition.
	AddParam(ctx ServerContext, name string, desc string, datatype datatypes.DataType, collection, required, stream bool) error
	// AddParamWithType adds a required parameter with a specific type.
	AddParamWithType(ctx ServerContext, name string, desc string, datatype datatypes.DataType) error
	// AddOptionalParamWithType adds an optional parameter with a specific type.
	AddOptionalParamWithType(ctx ServerContext, name string, desc string, datatype datatypes.DataType) error
	//AddCollectionParams(ctx ServerContext, params map[string]datatypes.DataType) error
	//	SetRequestType(ctx ServerContext, datatype string, collection bool, stream bool)
	//	SetResponseType(ctx ServerContext, stream bool)
	// SetDescription sets the description of the service.
	SetDescription(ServerContext, string)
	// SetComponent marks the service as a component.
	SetComponent(ServerContext, bool)
	// Get tags for a service
	GetTags(ServerContext) []*Tag

	ServerElement() ServerElement
	//ConfigureService(ctx ServerContext, requestType string, collection bool, stream bool, params []string, config []string, description string)
	//ConfigureService(ctx ServerContext, params []string, config []string, description string)
}
type Tag struct {
	Name string
	Description string
	ParentTag *Tag
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
	GetService(RequestContext) Service
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
