package elements

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
)

type ServiceResponseHandler interface {
	core.ServerElement
	core.ResponseHandler
	Initialize(ctx core.ServerContext, conf config.Config) error
}
