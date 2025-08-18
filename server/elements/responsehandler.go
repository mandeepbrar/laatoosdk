package elements

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
)

type ServiceResponseHandler interface {
	core.ServerElement
	Initialize(ctx core.ServerContext, conf config.Config) error
	HandleResponse(ctx core.RequestContext, resp *core.Response, err error) error
}
