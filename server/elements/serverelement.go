package elements

import (
	"laatoo/sdk/common/config"
	"laatoo/sdk/server/core"
)

//writable
type ServerElementHandle interface {
	Initialize(ctx core.ServerContext, conf config.Config) error
	Start(ctx core.ServerContext) error
}

type Server interface {
	core.ServerElement
}
