package elements

import (
	//	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type Engine interface {
	core.ServerElement
	GetRootChannel(ctx core.ServerContext) Channel
	GetRequestParams(ctx core.RequestContext) utils.StringMap
	GetDefaultResponseHandler(ctx core.ServerContext) ServiceResponseHandler
}
