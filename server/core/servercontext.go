package core

import (
	"laatoo/sdk/common/config"
	"laatoo/sdk/server/ctx"
)

/*application and engine types*/
const (
	CONF_SERVERTYPE_STANDALONE = "STANDALONE"
	CONF_SERVERTYPE_GOOGLEAPP  = "GOOGLE_APP"
	CONF_ENGINE_HTTP           = "http"
	CONF_ENGINE_TCP            = "tcp"
)

type ServerElementType int

const (
	ServerElementEngine ServerElementType = iota
	ServerElementEnvironment
	ServerElementLoader
	ServerElementServiceFactory
	ServerElementServiceManager
	ServerElementChannel
	ServerElementChannelManager
	ServerElementFactoryManager
	ServerElementApplication
	ServerElementRulesManager
	ServerElementService
	ServerElementServiceResponseHandler
	ServerElementServer
	ServerElementSessionManager
	ServerElementSecurityHandler
	ServerElementMessagingManager
	ServerElementModuleManager
	ServerElementTaskManager
	ServerElementCacheManager
	ServerElementModule
	ServerElementLogger
	ServerElementCommunicator
	ServerElementSecretsManager
	ServerElementOpen1
	ServerElementOpen2
	ServerElementOpen3
)

type ContextMap map[ServerElementType]ServerElement

type ServerElement interface {
	Reference() ServerElement
	GetProperty(string) interface{}
	GetName() string
	GetContext() ServerContext
	GetType() ServerElementType
}

type ServerContext interface {
	ctx.Context
	GetServerElement(ServerElementType) ServerElement
	GetService(alias string) (Service, error)
	GetServiceContext(alias string) (ServerContext, error)
	//NewContext(name string) ServerContext
	SubContext(name string) ServerContext
	GetServerProperties() map[string]interface{}
	CreateNewRequest(name string, engine interface{}, engineCtx interface{}, sessionId string) (RequestContext, error)
	CreateCollection(objectName string, length int) (interface{}, error)
	CreateObjectPointersCollection(objectName string, length int) (interface{}, error)
	CreateObject(objectName string) (interface{}, error)
	GetObjectFactory(name string) (ObjectFactory, bool)
	GetObjectMetadata(objectName string) (Info, error)
	CreateSystemRequest(name string) RequestContext
	SubscribeTopic(topics []string, lstnr MessageListener, lsnrID string) error
	CreateConfig() config.Config
	GetCodec(encoding string) (Codec, bool)
	ReadConfigMap(cfg map[string]interface{}) (config.Config, error)
	ReadConfigData(data []byte, funcs map[string]interface{}) (config.Config, error)
	ReadConfig(file string, funcs map[string]interface{}) (config.Config, error)
	GetRegName(object interface{}) (string, bool, bool)
}
