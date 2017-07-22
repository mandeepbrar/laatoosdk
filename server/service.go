package server

import (
	"io"
	"laatoo/sdk/config"
	"laatoo/sdk/core"
)

type Service interface {
	core.ServerElement
	Service() core.Service
	GetConfiguration() config.Config
	//Invoke(core.RequestContext, *core.Request) (*core.Response, error)
	HandleEncodedRequest(ctx core.RequestContext, vals map[string]interface{}, body []byte) (*core.Response, error)
	HandleRequest(ctx core.RequestContext, vals map[string]interface{}, body interface{}) (*core.Response, error)
	HandleStreamedRequest(ctx core.RequestContext, info map[string]interface{}, body io.ReadCloser) (*core.Response, error)
}
