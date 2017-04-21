package server

import (
	"laatoo/sdk/config"
	"laatoo/sdk/core"
)

type ChannelManager interface {
	core.ServerElement
	Serve(ctx core.ServerContext, channelName string, svc Service, channelConfig config.Config) error
}
