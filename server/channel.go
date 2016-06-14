package server

import (
	"laatoo/sdk/config"
	"laatoo/sdk/core"
)

type Channel interface {
	core.ServerElement
	Serve(ctx core.ServerContext, service Service, channelConfig config.Config) error
	Child(ctx core.ServerContext, name string, channelConfig config.Config) (Channel, error)
}
