package server

import (
	"laatoo/sdk/config"
	"laatoo/sdk/core"
)

type Channel interface {
	core.ServerElement
	GetServiceName() string
	Serve(ctx core.ServerContext) error
	Child(ctx core.ServerContext, name string, channelConfig config.Config) (Channel, error)
}
