package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/utils"
)

// RequestContext is the context passed while executing a request through all the layers of execution.
// It carries information about the current request, user, session, and provides access to server elements.
type RequestContext interface {
	ctx.Context
	// ServerContext returns the server context that generated this request.
	ServerContext() ServerContext
	// EngineRequestContext returns the context from the engine (e.g., HTTP engine) that received the request.
	EngineRequestContext() EngineContext
	// EngineRequestParams returns the parameters extracted from the engine request.
	EngineRequestParams() utils.StringMap
	// SubContext creates a sub-context of the request.
	// It retains the ID and tracks flow along with variables.
	SubContext(name string) RequestContext
	// GetServerElement returns a server element applicable to the context by its type.
	GetServerElement(elemType ServerElementType) ServerElement
	//NewContext(name string) RequestContext
	// GetRequest returns the request object containing parameters.
	GetRequest() Request
	// SetResponse sets the response for the request.
	SetResponse(*Response)
	// AddResponseInfo adds additional metadata to the response.
	AddResponseInfo(utils.StringMap)
	// GetSession returns the session associated with the request.
	GetSession() Session
	// GetFromSession retrieves a value from the session.
	GetFromSession(key string) (interface{}, bool)
	// SetInSession sets a value in the session.
	SetInSession(key string, val interface{})
	// GetResponse returns the current response object.
	GetResponse() *Response
	//GetBody() interface{}
	// GetParam retrieves a parameter by name.
	GetParam(string) (Param, bool)
	// GetParams retrieves all parameters.
	GetParams() map[string]Param
	// GetParamValue retrieves the value of a parameter.
	GetParamValue(string) (interface{}, bool)
	// GetIntParam retrieves an integer parameter.
	GetIntParam(string) (int, bool)
	// GetStringParam retrieves a string parameter.
	GetStringParam(string) (string, bool)
	// GetConfigParam retrieves a configuration parameter.
	GetConfigParam(string) (config.Config, bool)
	// GetConfigArrParam retrieves an array of configuration parameters.
	GetConfigArrParam(string) ([]config.Config, bool)
	// GetStringMapParam retrieves a string map parameter.
	GetStringMapParam(string) (utils.StringMap, bool)
	// GetStringsMapParam retrieves a strings map parameter.
	GetStringsMapParam(string) (utils.StringsMap, bool)
	// Invoke invokes a service by alias with parameters.
	Invoke(alias string, params utils.StringMap) (*Response, error)
	// Forward forwards the request to another service.
	Forward(string, utils.StringMap) error
	// ForwardToService forwards the request to a specific service instance.
	ForwardToService(Service, utils.StringMap) error
	// GetUser returns the user executing the request.
	GetUser() auth.User
	// GetTenant returns the tenant info for the request.
	GetTenant() auth.TenantInfo
	// HasPermission checks if the user has a specific permission.
	HasPermission(perm string) bool
	// GetObjectFactory retrieves an object factory by name.
	GetObjectFactory(name string) (ObjectFactory, bool)
	// CreateCollection creates a collection of objects.
	CreateCollection(objectName string, length int) (interface{}, error)
	// CreateObject creates a new object instance.
	CreateObject(objectName string) (interface{}, error)
	// CreateObjectPointersCollection creates a collection of object pointers.
	CreateObjectPointersCollection(objectName string, length int) (interface{}, error)
	// PublishMessage publishes a message to a topic.
	PublishMessage(topic string, message *Message)
	// SendSynchronousMessage sends a synchronous message.
	SendSynchronousMessage(msgType string, data interface{}) error
	// PutInCache puts an item in the cache.
	PutInCache(bucket string, key string, item interface{}) error
	// PutMultiInCache puts multiple items in the cache.
	PutMultiInCache(bucket string, vals utils.StringMap) error
	// GetFromCache retrieves an item from the cache.
	GetFromCache(bucket string, key string) (interface{}, bool)
	// GetMultiFromCache retrieves multiple items from the cache.
	GetMultiFromCache(bucket string, keys []string) utils.StringMap
	// GetObjectFromCache retrieves an object from the cache.
	GetObjectFromCache(bucket string, key string, objectType string) (interface{}, bool)
	// IncrementInCache increments a value in the cache.
	IncrementInCache(bucket string, key string) error
	// DecrementInCache decrements a value in the cache.
	DecrementInCache(bucket string, key string) error
	// GetObjectsFromCache retrieves multiple objects from the cache.
	GetObjectsFromCache(bucket string, keys []string, objectType string) utils.StringMap
	// PushTask pushes a task to a queue. It returns the task ID on success.
	PushTask(queue string, taskdata interface{}, metadata utils.StringMap) (string, error)
	// SubscribeTaskCompletion subscribes to task completion events.
	SubscribeTaskCompletion(topic string, handler MessageListener, subscriberId string) error
	// StartWorkflow starts a workflow.
	StartWorkflow(workflowName string, initData utils.StringMap, insconf utils.StringMap) (interface{}, error)
	// InvalidateCache invalidates a cache entry.
	InvalidateCache(bucket string, key string) error
	// GetCodec retrieves a codec by encoding.
	GetCodec(encoding string) (datatypes.Codec, bool)
	// GetRegName retrieves the registered name of an object.
	GetRegName(object interface{}) (string, bool, bool)
	// GetExpressionValue evaluates an expression.
	GetExpressionValue(expression Expression, vars utils.StringMap) (interface{}, error)
	// InvokeActivity invokes an activity.
	InvokeActivity(activity string, params utils.StringMap) (interface{}, error)

	// SendNotification sends a notification.
	SendNotification(notification *Notification) error
	// CompleteRequest marks the request as complete.
	CompleteRequest()
}
