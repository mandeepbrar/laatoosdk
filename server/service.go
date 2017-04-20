package server

import (
	"laatoo/sdk/config"
	"laatoo/sdk/core"
)

type Service interface {
	core.ServerElement
	Service() core.Service
	ParamsConfig() config.Config
	Invoke(ctx core.RequestContext) error
}
