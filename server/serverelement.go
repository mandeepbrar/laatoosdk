package server

import (
	"laatoo/sdk/config"
	"laatoo/sdk/core"
)

//writable
type ServerElementHandle interface {
	Initialize(ctx core.ServerContext, conf config.Config) error
	Start(ctx core.ServerContext) error
}

type Server interface {
	core.ServerElement
}
