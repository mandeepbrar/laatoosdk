package elements

import (
	"laatoo/sdk/server/core"
	"laatoo/sdk/common/config"
)

type ServiceResponseHandler interface {
	core.ServerElement
	Initialize(ctx core.ServerContext, conf config.Config) error
	HandleResponse(ctx core.RequestContext, resp *core.Response) error
}
