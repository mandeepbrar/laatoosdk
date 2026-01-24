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
	ServerElementAgentManager
	ServerElementOpen1
	ServerElementOpen2
	ServerElementOpen3
)

type ContextMap map[ServerElementType]ServerElement

// ServerElement is the base interface for all server components.
type ServerElement interface {
	Reference() ServerElement
	GetProperty(string) interface{}
	GetName() string
	GetContext() ServerContext
	GetType() ServerElementType
}

// ServerContext is the context passed during initialization of factories and services.
// It acts as a proxy to the server and provides access to server elements and configuration.
type ServerContext interface {
	ctx.Context
	// GetServerElement returns a server element applicable to the context by its type.
	GetServerElement(ServerElementType) ServerElement
	// GetService retrieves a service by its alias.
	GetService(alias string) (Service, error)
	// GetServiceContext retrieves the server context for a specific service.
	GetServiceContext(alias string) (ServerContext, error)
	//NewContext(name string) ServerContext
	// SubContext creates a child context with the same underlying context.
	// Changes made to context parameters will be visible on the parent.
	SubContext(name string) ServerContext
	// GetServerProperties returns the properties of the server.
	GetServerProperties() utils.StringMap
	// CreateNewRequest creates a new request context.
	CreateNewRequest(name string, tenant auth.TenantInfo, engine interface{}, engineCtx EngineContext, sessionId string) (RequestContext, error)
	// CreateCollection creates a collection of objects.
	CreateCollection(objectName string, length int) (interface{}, error)
	// CreateObjectPointersCollection creates a collection of object pointers.
	CreateObjectPointersCollection(objectName string, length int) (interface{}, error)
	// CreateObject creates a new object instance.
	CreateObject(objectName string) (interface{}, error)
	// GetObjectFactory retrieves an object factory by name.
	GetObjectFactory(name string) (ObjectFactory, bool)
	// GetObjectMetadata retrieves metadata for an object.
	GetObjectMetadata(objectName string) (Info, error)
	// CreateSystemRequest creates a system request context (e.g., for background tasks).
	CreateSystemRequest(name string, tenant auth.TenantInfo, behalfOf interface{}) RequestContext
	// SubscribeTopic subscribes to a message topic.
	SubscribeTopic(topics []string, lstnr MessageListener, lsnrID string) error
	// CreateConfig creates a new configuration object.
	CreateConfig() config.Config
	// GetCodec retrieves a codec by encoding.
	GetCodec(encoding string) (datatypes.Codec, bool)
	// RegisterExpression registers a new expression type.
	RegisterExpression(expression Expression, dtype datatypes.DataType) error
	// ReadConfigMap reads configuration from a map.
	ReadConfigMap(cfg map[string]interface{}) (config.Config, error)
	// ReadConfigData reads configuration from byte data.
	ReadConfigData(data []byte, funcs map[string]interface{}) (config.Config, error)
	// ReadConfig reads configuration from a file.
	ReadConfig(file string, funcs map[string]interface{}) (config.Config, error)
	// GetRegName retrieves the registered name of an object.
	GetRegName(object interface{}) (string, bool, bool)
	// GetRegisteredComponent retrieves a registered component by name.
	GetRegisteredComponent(obj string) (interface{}, error)
	// GetLogLevel returns the current log level.
	GetLogLevel() int
	// GetLogFormat returns the current log format.
	GetLogFormat() string
	// GetTenant returns the tenant info for the context.
	GetTenant() auth.TenantInfo
}

type EngineContext interface {
	GetRequest() interface{}
	GetRequestStream() (interface{}, error)
	GetResponseStream() (interface{}, error)
	GetConnection() interface{}
	GetUnderlyingContext() interface{}
}
