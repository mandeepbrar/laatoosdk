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
	ServerElementKnowledgeManager
)

type ContextMap map[ServerElementType]ServerElement

// ServerElement is the handle a plugin gets on one of the server's managers — the data manager,
// the service manager, the security handler and so on — obtained by asking a context for an
// element by ServerElementType and asserting the result to the manager interface you want.
//
// What is handed out is a proxy over the live manager rather than the manager itself, which is
// why the element type and the interface you assert to are separate things.
type ServerElement interface {
	// Reference returns another handle onto the SAME underlying manager. It is a new proxy, not
	// a copy of any state: both handles see the one manager.
	Reference() ServerElement

	// GetProperty reads a named property of the element, or nil when there is none.
	//
	// A nil result does not distinguish "no such property" from "this element does not support
	// property lookup at all" — several managers answer nil unconditionally. Do not infer
	// absence of a configured value from it.
	GetProperty(string) interface{}

	// GetName returns the element's registered name.
	GetName() string

	// GetContext returns the ServerContext this element belongs to, which is the level
	// (server, solution, application, isolation) it was created at.
	GetContext() ServerContext

	// GetType returns the discriminator identifying which manager this is.
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
	CreateNewRequest(name string, tenant auth.TenantInfo, engine interface{}, engineCtx EngineContext, sessionId string, responseHandler ResponseHandler) (RequestContext, error)
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
	CreateSystemRequest(name string, tenant auth.TenantInfo, behalfOf interface{}, responseHandler ResponseHandler) RequestContext
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
	// GetRegisteredDataComponent retrieves a registered data component by name.
	GetRegisteredDataComponent(obj string) (interface{}, error)
	// GetLogLevel returns the current log level.
	GetLogLevel() int
	// GetLogFormat returns the current log format.
	GetLogFormat() string
}

// EngineContext exposes the transport underneath a request, for the cases the platform's own
// parameter and response handling cannot serve — streaming a large body, upgrading a connection,
// or reading something engine-specific.
//
// Everything here returns interface{} and the concrete type is decided by the ENGINE serving the
// request, not by this interface. Code that type-asserts to one engine's types compiles fine and
// panics under another, so an assertion should be the comma-ok form and the engine it assumes
// should be stated where it is made.
//
// Reach for this only when the ordinary path will not do: parameters through the Get*Param family
// and replies through SetResponse work across every engine, and this does not.
type EngineContext interface {
	// GetRequest returns the engine's request object — an *http.Request on the HTTP engines.
	GetRequest() interface{}

	// GetRequestStream returns the request body as a live stream, an io.ReadCloser on the HTTP
	// engines, so a large payload can be consumed incrementally instead of being buffered.
	//
	// It is the reason a plugin never needs to read the whole body into memory to process it.
	// The body is consumed by reading: taking it here and also declaring a body-bearing channel
	// parameter means one of the two gets nothing.
	GetRequestStream() (interface{}, error)

	// GetResponseStream returns the response writer as a live stream, an http.ResponseWriter on
	// the HTTP engines, for writing a reply incrementally rather than as one value.
	//
	// Writing here bypasses the platform's response handling, so status and content type are
	// the caller's own responsibility.
	GetResponseStream() (interface{}, error)

	// GetConnection returns the underlying connection, which is what a protocol needing to hold
	// one open — a websocket upgrade, a long-lived stream — starts from.
	GetConnection() interface{}

	// GetUnderlyingContext returns the engine framework's own context object. On some engines
	// this is the same object GetConnection returns.
	GetUnderlyingContext() interface{}
}
