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
	NewContext(name string) RequestContext
	GetUser() auth.User
	HasPermission(perm string) bool
	PublishMessage(topic string, message interface{})
	SendSynchronousMessage(msgType string, data interface{}) error
	PutInCache(key string, item interface{}) error
	GetFromCache(key string, val interface{}) bool
	GetMultiFromCache(keys []string, val map[string]interface{}) bool
	PushTask(queue string, task interface{}) error
	InvalidateCache(key string) error
	IsAdmin() bool
	SetRequest(interface{})
	GetRequest() interface{}
	SetResponse(*ServiceResponse)
	GetResponse() *ServiceResponse
	CompleteRequest()
}
