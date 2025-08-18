package elements

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
)

type Channel interface {
	core.ServerElement
	GetServiceName() string
	Serve(ctx core.ServerContext) error
	Child(ctx core.ServerContext, name string, channelConfig config.Config) (Channel, error)
	Destruct(ctx core.ServerContext, parentChannel Channel) error
	GetEngine(ctx core.ServerContext) Engine
}
