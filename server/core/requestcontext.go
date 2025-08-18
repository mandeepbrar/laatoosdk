package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/utils"
)

type RequestContext interface {
	ctx.Context
	ServerContext() ServerContext
	EngineRequestContext() EngineContext
	EngineRequestParams() utils.StringMap
	SubContext(name string) RequestContext
	GetServerElement(elemType ServerElementType) ServerElement
	//NewContext(name string) RequestContext
	GetRequest() Request
	SetResponse(*Response)
	AddResponseInfo(utils.StringMap)
	GetSession() Session
	GetFromSession(key string) (interface{}, bool)
	SetInSession(key string, val interface{})
	GetResponse() *Response
	//GetBody() interface{}
	GetParam(string) (Param, bool)
	GetParams() map[string]Param
	GetParamValue(string) (interface{}, bool)
	GetIntParam(string) (int, bool)
	GetStringParam(string) (string, bool)
	GetConfigParam(string) (config.Config, bool)
	GetConfigArrParam(string) ([]config.Config, bool)
	GetStringMapParam(string) (utils.StringMap, bool)
	GetStringsMapParam(string) (utils.StringsMap, bool)
	Invoke(alias string, params utils.StringMap) (*Response, error)
	Forward(string, utils.StringMap) error
	ForwardToService(Service, utils.StringMap) error
	GetUser() auth.User
	GetTenant() auth.TenantInfo
	HasPermission(perm string) bool
	GetObjectFactory(name string) (ObjectFactory, bool)
	CreateCollection(objectName string, length int) (interface{}, error)
	CreateObject(objectName string) (interface{}, error)
	CreateObjectPointersCollection(objectName string, length int) (interface{}, error)
	PublishMessage(topic string, message interface{})
	SendSynchronousMessage(msgType string, data interface{}) error
	PutInCache(bucket string, key string, item interface{}) error
	PutMultiInCache(bucket string, vals utils.StringMap) error
	GetFromCache(bucket string, key string) (interface{}, bool)
	GetMultiFromCache(bucket string, keys []string) utils.StringMap
	GetObjectFromCache(bucket string, key string, objectType string) (interface{}, bool)
	IncrementInCache(bucket string, key string) error
	DecrementInCache(bucket string, key string) error
	GetObjectsFromCache(bucket string, keys []string, objectType string) utils.StringMap
	PushTask(queue string, taskdata interface{}) error
	SubscribeTaskCompletion(queue string, callback func(ctx RequestContext, invocationId string, result interface{})) error
	StartWorkflow(workflowName string, initData utils.StringMap, insconf utils.StringMap) (interface{}, error)
	InvalidateCache(bucket string, key string) error
	GetCodec(encoding string) (datatypes.Codec, bool)
	GetRegName(object interface{}) (string, bool, bool)
	GetExpressionValue(expression Expression, vars utils.StringMap) (interface{}, error)
	InvokeActivity(activity string, params utils.StringMap) (interface{}, error)
	InvokeScript(script string, params utils.StringMap) (interface{}, error)
	ExecuteAction(actiontype ActionType, params utils.StringMap) (interface{}, error)
	SendNotification(notification *Notification) error
	CompleteRequest()
}
