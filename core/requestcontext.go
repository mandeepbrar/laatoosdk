package core

import (
	"laatoo/sdk/auth"
)

type RequestContext interface {
	Context
	ServerContext() ServerContext
	EngineRequestContext() interface{}
	SubContext(name string) RequestContext
	GetServerElement(elemType ServerElementType) ServerElement
	//NewContext(name string) RequestContext
	GetUser() auth.User
	HasPermission(perm string) bool
	PublishMessage(topic string, message interface{})
	SendSynchronousMessage(msgType string, data interface{}) error
	PutInCache(bucket string, key string, item interface{}) error
	PutMultiInCache(bucket string, vals map[string]interface{}) error
	GetFromCache(bucket string, key string) (interface{}, bool)
	GetMultiFromCache(bucket string, keys []string) map[string]interface{}
	GetObjectFromCache(bucket string, key string, objectType string) (interface{}, bool)
	CreateRequest() ServiceRequest
	IncrementInCache(bucket string, key string) error
	DecrementInCache(bucket string, key string) error
	GetObjectsFromCache(bucket string, keys []string, objectType string) map[string]interface{}
	PushTask(queue string, task interface{}) error
	InvalidateCache(bucket string, key string) error
	IsAdmin() bool
	CompleteRequest()
}
