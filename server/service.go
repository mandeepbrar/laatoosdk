package server

import (
	"laatoo/sdk/core"
)

type Service interface {
	core.ServerElement
	Service() core.Service
	Invoke(ctx core.RequestContext) error
}
