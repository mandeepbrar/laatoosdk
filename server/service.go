package server

import (
	"io"
	"laatoo/sdk/core"
)

type Service interface {
	core.ServerElement
	Service() core.Service
	//Invoke(core.RequestContext, *core.Request) (*core.ServiceResponse, error)
	HandleRequest(ctx core.RequestContext, info map[string]interface{}, body []byte) (*core.ServiceResponse, error)
	HandleStreamedRequest(ctx core.RequestContext, info map[string]interface{}, body io.ReadCloser) (*core.ServiceResponse, error)
}
