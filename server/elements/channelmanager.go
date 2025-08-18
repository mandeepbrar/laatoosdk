package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type ChannelManager interface {
	core.ServerElement
	GetChannel(ctx core.ServerContext, name string) (Channel, bool)
	List(ctx core.ServerContext) utils.StringsMap
	ChangeLogger(ctx core.ServerContext, chanName string, logLevel string, logFormat string) error
	Describe(ctx core.ServerContext, chanName string) (utils.StringMap, error)
}
