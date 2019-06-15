package elements

import (
	//	"laatoo/sdk/common/config"
	"laatoo/sdk/server/core"
)

type Engine interface {
	core.ServerElement
	GetRootChannel(ctx core.ServerContext) Channel
	GetRequestParams(ctx core.RequestContext) map[string]interface{}
}
