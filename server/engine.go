package server

import (
	//	"laatoo/sdk/config"
	"laatoo/sdk/core"
)

type Engine interface {
	core.ServerElement
	GetRootChannel(ctx core.ServerContext) Channel
}
