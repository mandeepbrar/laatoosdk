package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/utils"
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
	ServerElementSolution
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
	ServerElementIsolation
	ServerElementSessionManager
	ServerElementSecurityHandler
	ServerElementMessagingManager
	ServerElementModuleManager
	ServerElementTaskManager
	ServerElementCacheManager
	ServerElementModule
	ServerElementLogger
	ServerElementNotificationManager
	ServerElementSecretsManager
	ServerElementWorkflowManager
	ServerElementActivityManager
	ServerElementScriptManager
	ServerElementExpressionManager
	ServerElementActionsManager
	ServerElementDataManager
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
	GetServerProperties() utils.StringMap
	CreateNewRequest(name string, tenant auth.TenantInfo, engine interface{}, engineCtx EngineContext, sessionId string) (RequestContext, error)
	CreateCollection(objectName string, length int) (interface{}, error)
	CreateObjectPointersCollection(objectName string, length int) (interface{}, error)
	CreateObject(objectName string) (interface{}, error)
	GetObjectFactory(name string) (ObjectFactory, bool)
	GetObjectMetadata(objectName string) (Info, error)
	CreateSystemRequest(name string, tenant auth.TenantInfo, behalfOf interface{}) RequestContext
	SubscribeTopic(topics []string, lstnr MessageListener, lsnrID string) error
	CreateConfig() config.Config
	GetCodec(encoding string) (datatypes.Codec, bool)
	RegisterExpression(expression Expression, dtype datatypes.DataType) error
	ReadConfigMap(cfg map[string]interface{}) (config.Config, error)
	ReadConfigData(data []byte, funcs map[string]interface{}) (config.Config, error)
	ReadConfig(file string, funcs map[string]interface{}) (config.Config, error)
	GetRegName(object interface{}) (string, bool, bool)
	GetRegisteredComponent(obj string) (interface{}, error)
	GetLogLevel() int
	GetLogFormat() string
	GetTenant() auth.TenantInfo
}

type EngineContext interface {
	GetRequest() interface{}
}
