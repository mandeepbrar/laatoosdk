package core

import (
	"laatoo/sdk/server/auth"
	"laatoo/sdk/server/ctx"
)

type RequestContext interface {
	ctx.Context
	ServerContext() ServerContext
	EngineRequestContext() interface{}
	EngineRequestParams() map[string]interface{}
	SubContext(name string) RequestContext
	GetServerElement(elemType ServerElementType) ServerElement
	//NewContext(name string) RequestContext
	GetRequest() Request
	SetResponse(*Response)
	GetSession() Session
	GetResponse() *Response
	//GetBody() interface{}
	GetParam(string) (Param, bool)
	GetParams() map[string]Param
	GetParamValue(string) (interface{}, bool)
	GetIntParam(string) (int, bool)
	GetStringParam(string) (string, bool)
	GetStringMapParam(string) (map[string]interface{}, bool)
	GetStringsMapParam(string) (map[string]string, bool)
	Forward(string, map[string]interface{}) error
	ForwardToService(Service, map[string]interface{}) error
	GetUser() auth.User
	HasPermission(perm string) bool
	GetObjectFactory(name string) (ObjectFactory, bool)
	CreateCollection(objectName string, length int) (interface{}, error)
	CreateObject(objectName string) (interface{}, error)
	CreateObjectPointersCollection(objectName string, length int) (interface{}, error)
	PublishMessage(topic string, message interface{})
	SendSynchronousMessage(msgType string, data interface{}) error
	PutInCache(bucket string, key string, item interface{}) error
	PutMultiInCache(bucket string, vals map[string]interface{}) error
	GetFromCache(bucket string, key string) (interface{}, bool)
	GetMultiFromCache(bucket string, keys []string) map[string]interface{}
	GetObjectFromCache(bucket string, key string, objectType string) (interface{}, bool)
	IncrementInCache(bucket string, key string) error
	DecrementInCache(bucket string, key string) error
	GetObjectsFromCache(bucket string, keys []string, objectType string) map[string]interface{}
	PushTask(queue string, task interface{}) error
	InvalidateCache(bucket string, key string) error
	GetCodec(encoding string) (Codec, bool)
	SendCommunication(communication interface{}) error
	GetRegName(object interface{}) (string, bool, bool)
	IsAdmin() bool
	CompleteRequest()
}
