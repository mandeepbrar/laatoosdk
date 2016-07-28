package core

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
	ServerElementSecurityHandler
	ServerElementMessagingManager
	ServerElementTaskManager
	ServerElementCacheManager
	ServerElementOpen1
	ServerElementOpen2
	ServerElementOpen3
)

type ContextMap map[ServerElementType]ServerElement

type ServerElement interface {
	Context
}

type ServerContext interface {
	Context
	GetServerType() string
	GetElement() ServerElement
	GetServerElement(ServerElementType) ServerElement
	GetService(alias string) (Service, error)
	GetElementType() ServerElementType
	NewContext(name string) ServerContext
	NewContextWithElements(name string, elements ContextMap, primaryElement ServerElementType) ServerContext
	SubContext(name string) ServerContext
	SubContextWithElement(name string, primaryElement ServerElementType) ServerContext
	CreateNewRequest(name string, engineCtx interface{}) RequestContext
	CreateCollection(objectName string, length int) (interface{}, error)
	CreateObject(objectName string) (interface{}, error)
	GetMethod(methodName string) (ServiceFunc, error)
	GetObjectCollectionCreator(objectName string) (ObjectCollectionCreator, error)
	GetObjectCreator(objectName string) (ObjectCreator, error)
	CreateSystemRequest(name string) RequestContext
	SubscribeTopic(topics []string, lstnr ServiceFunc) error
}
