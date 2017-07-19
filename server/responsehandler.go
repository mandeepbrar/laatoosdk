package server

import (
	"laatoo/sdk/core"
)

type ServiceResponseHandler interface {
	core.ServerElement
	HandleResponse(ctx core.RequestContext, resp *core.ServiceResponse) error
}
