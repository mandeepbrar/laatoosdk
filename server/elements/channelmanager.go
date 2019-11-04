package elements

import "laatoo/sdk/server/core"

type ChannelManager interface {
	core.ServerElement
	GetChannel(ctx core.ServerContext, name string) (Channel, bool)
}
