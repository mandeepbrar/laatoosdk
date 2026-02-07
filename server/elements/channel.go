package elements

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
)

type Channel interface {
	core.ServerElement
	GetServiceName() string
	Serve(ctx core.ServerContext) error
	Child(ctx core.ServerContext, name string, channelConfig config.Config, module core.Module) (Channel, error)
	Destruct(ctx core.ServerContext, parentChannel Channel) error
	GetModule() core.Module
	GetEngine(ctx core.ServerContext) Engine
	GetDescription(ctx core.ServerContext) string
}
