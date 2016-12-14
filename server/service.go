package server

import (
	"laatoo/sdk/config"
	"laatoo/sdk/core"
)

type Service interface {
	core.ServerElement
	Service() core.Service
	Config() config.Config
	Invoke(ctx core.RequestContext) error
}
