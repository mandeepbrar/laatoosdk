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
	ServerElementLogger
	ServerElementOpen1
	ServerElementOpen2
	ServerElementOpen3
)

type ContextMap map[ServerElementType]ServerElement

type ServerElement interface {
	Reference() ServerElement
	GetProperty(string) interface{}
	GetName() string
	GetType() ServerElementType
}

type ServerContext interface {
	Context
	GetServerElement(ServerElementType) ServerElement
	GetService(alias string) (Service, error)
	//NewContext(name string) ServerContext
	SubContext(name string) ServerContext
	CreateNewRequest(name string, params interface{}) RequestContext
	CreateCollection(objectName string, length int) (interface{}, error)
	CreateObject(objectName string) (interface{}, error)
	GetObjectCollectionCreator(objectName string) (ObjectCollectionCreator, error)
	GetObjectCreator(objectName string) (ObjectCreator, error)
	CreateSystemRequest(name string) RequestContext
	SubscribeTopic(topics []string, lstnr MessageListener) error
}
