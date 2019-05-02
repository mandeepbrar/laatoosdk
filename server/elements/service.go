package elements

import (
	"laatoo/sdk/common/config"
	"laatoo/sdk/server/core"
)

type Service interface {
	core.ServerElement
	Service() core.Service
	ServiceContext() core.ServerContext
	GetConfiguration() config.Config
	//Invoke(core.RequestContext, *core.Request) (*core.Response, error)
	//HandleEncodedRequest(ctx core.RequestContext, vals map[string]interface{}, body []byte) (*core.Response, error)
	HandleRequest(ctx core.RequestContext, vals map[string]interface{}) (*core.Response, error)
	//HandleStreamedRequest(ctx core.RequestContext, info map[string]interface{}, body io.ReadCloser) (*core.Response, error)
}
