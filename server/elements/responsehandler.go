package elements

import (
	"laatoo/sdk/server/core"
)

type ServiceResponseHandler interface {
	core.ServerElement
	HandleResponse(ctx core.RequestContext, resp *core.Response) error
}
