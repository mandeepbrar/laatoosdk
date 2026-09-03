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

// Values are ASSIGNED EXPLICITLY, never by iota.
//
// Under iota, removing a constant renumbers every one after it, and a plugin .so compiled
// against the old numbering then resolves a different element than the server means — an unset
// extension slot, whose nil interface panics at the caller's type assertion far from the cause.
// Explicit values make a future removal cost nothing.
//
// ServerElementNamespace replaces the former Solution, Application and Isolation constants:
// levels are namespaces now, and the server holds no level vocabulary.
const (
	ServerElementEngine              ServerElementType = 0
	ServerElementNamespace           ServerElementType = 1
	ServerElementLoader              ServerElementType = 2
	ServerElementServiceFactory      ServerElementType = 3
	ServerElementServiceManager      ServerElementType = 4
	ServerElementChannel             ServerElementType = 5
	ServerElementChannelManager      ServerElementType = 6
	ServerElementFactoryManager      ServerElementType = 7
	ServerElementRulesManager        ServerElementType = 8
	ServerElementService             ServerElementType = 9
	ServerElementSessionManager      ServerElementType = 10
	ServerElementSecurityHandler     ServerElementType = 11
	ServerElementMessagingManager    ServerElementType = 12
	ServerElementModuleManager       ServerElementType = 13
	ServerElementTaskManager         ServerElementType = 14
	ServerElementCacheManager        ServerElementType = 15
	ServerElementModule              ServerElementType = 16
	ServerElementLogger              ServerElementType = 17
	ServerElementNotificationManager ServerElementType = 18
	ServerElementSecretsManager      ServerElementType = 19
	ServerElementWorkflowManager     ServerElementType = 20
	ServerElementActivityManager     ServerElementType = 21
	ServerElementScriptManager       ServerElementType = 22
	ServerElementExpressionManager   ServerElementType = 23
	ServerElementActionsManager      ServerElementType = 24
	ServerElementDataManager         ServerElementType = 25
	ServerElementAgentManager        ServerElementType = 26
	ServerElementOpen1               ServerElementType = 27
	ServerElementOpen2               ServerElementType = 28
	ServerElementOpen3               ServerElementType = 29
	ServerElementKnowledgeManager    ServerElementType = 30

	// ServerElementNamespaceManager owns the namespace tree: discovery from configuration, holding,
	// resolution and lifecycle. Added at 31, the next free value -- NOT by repurposing one of the
	// Open slots above, because a spare renamed is a slot whose meaning changed under any plugin
	// still compiled against the old name, which is the hazard explicit numbering exists to prevent.
	ServerElementNamespaceManager ServerElementType = 31
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

	// GetContext returns the ServerContext of the NAMESPACE this handle is scoped to.
	//
	// It formerly returned the context of the level the element was CREATED at, which is the
	// captured-context defect this replaced: an element built at one scope and invoked from
	// another resolved through the scope that built it. Scope now arrives at the call.
	GetContext() ServerContext

	// GetType returns the discriminator identifying which manager this is.
	GetType() ServerElementType

	// GetAddress returns this element's fully qualified address in `::` form —
	// "myapp::scriptmanager::validate" for a script declared in myapp. The root's address is the
	// separator itself, "::", so a child of the root reads "::orders" rather than "orders".
	//
	// The element reports its OWN address rather than having one reconstructed from outside. That
	// is what makes resolution diagnosable: asking which of two same-named declarations a bare
	// reference bound to is answered by asking the bound element where it lives, rather than by
	// re-deriving the walk and hoping it matches.
	//
	// An address is stable across relocation. A script declared in a namespace that owns no
	// ScriptManager is held by whichever manager resolves, under its own address; when that
	// namespace later declares a manager of its own, the holder changes and the address does not.
	GetAddress() string

	// GetParent returns the element enclosing this one, or a NIL INTERFACE at the root.
	//
	// `el.GetParent() == nil` is the correct root test. Stated in the contract rather than left to
	// implementations because a typed nil reaching a type assertion has already caused a
	// production panic in this platform (laatooserver/src/core/servicemanager.go records one on
	// the remote-service path).
	//
	// This is a stored pointer and is deliberately NOT a checked lookup, unlike GetChild. Walking
	// up is unchecked because it needs no check: ancestors are already reachable from any element,
	// and coming back down requires GetChild, which is checked. Do not "fix" this into a lookup
	// for symmetry — it would put a lock acquisition on a hot path and buy nothing.
	//
	// The returned element may be narrowed by assertion where the relationship is total: a
	// namespace's parent is always a namespace, and a channel's is always a channel or an engine.
	GetParent() ServerElement

	// GetChild returns this element's named child, or nil when there is none.
	//
	// Descent is uniform: given an element and a name, the next element is that element's named
	// child, whatever kind either is. A caller walking "myapp::scriptmanager::validate" asks three
	// times and never has to know which segments are namespaces and which are managers.
	//
	// A LEAF ANSWERS EXACTLY AS A MISS DOES — both return nil. This is deliberate: a caller
	// descending a path does not care whether it stopped because the name was absent or because
	// the element it reached holds no children, and collapsing the two keeps every call site from
	// having to distinguish a case it has no use for.
	//
	// The ctx argument supplies the CALLER'S scope, which is what makes this checkable: a
	// namespace is a containment boundary, and descent that would leave the caller's own subtree
	// is refused. A refusal also returns nil; it is reported through the server's diagnostic
	// surface rather than through this signature, so that the common path stays a single nilable
	// return.
	GetChild(ctx ServerContext, name string) ServerElement
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
