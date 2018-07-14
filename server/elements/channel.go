package elements

import (
	"laatoo/sdk/common/config"
	"laatoo/sdk/server/core"
)

type Channel interface {
	core.ServerElement
	GetServiceName() string
	Serve(ctx core.ServerContext) error
	Child(ctx core.ServerContext, name string, channelConfig config.Config) (Channel, error)
}
