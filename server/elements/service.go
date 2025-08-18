package elements

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type Service interface {
	core.ServerElement
	Service() core.Service
	ServiceContext() core.ServerContext
	GetConfiguration() config.Config
	HandleRequest(ctx core.RequestContext, vals utils.StringMap, encoding utils.StringsMap) (*core.Response, error)
}
