package core

type Service interface {
	ConfigurableObject
	Initialize(ctx ServerContext) error
	Info() ServiceInfo
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

type ServiceInfo interface {
	GetRequestInfo() RequestInfo
	GetResponseInfo() ResponseInfo
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
