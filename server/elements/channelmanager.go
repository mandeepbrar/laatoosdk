package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ChannelManager owns the channels at one level of the server hierarchy: it collects channel
// configuration from server config, from module registry/channels directories and from the config
// directory, then builds the routing tree by resolving each channel's `parent:` to an already-built
// channel.
//
// PARENT RESOLUTION IS ITERATIVE, NOT ORDERED. Channels whose parent is not yet built are retried
// on the next pass, so declaration order does not matter. When a pass resolves nothing new, the
// remaining channels are reported as Core_Bad_Conf "Missing Parents for Channels" and the level
// fails to start -- a `parent:` naming a channel that does not exist is a boot failure, not a
// silently dropped endpoint. The same check runs on the hot-load path.
//
// Note that a channel CAN still vanish quietly for other reasons: `disableroute: true` makes Serve
// bind nothing, and on the http engine a missing `path:` is not reported at all -- see
// Channel.Child.
type ChannelManager interface {
	core.ServerElement

	// GetChannel returns a built channel by name, and false when this level has none. The name is
	// the channel YAML's filename minus ".yml", unless the file declares `name:`. Only this level's
	// store is consulted.
	GetChannel(ctx core.ServerContext, name string) (Channel, bool)

	// List returns every channel at this level mapped to the name of the module that declared it.
	// A channel declared in server or application configuration rather than a plugin is reported
	// with the literal string "<no module>" -- named rather than omitted, since omitting it would
	// read as "not routed".
	List(ctx core.ServerContext) utils.StringsMap

	// ChangeLogger changes one channel's log level and format at runtime, with no restart. Returns
	// a not-found error for an unknown channel name.
	ChangeLogger(ctx core.ServerContext, chanName string, logLevel string, logFormat string) error

	// Describe returns an admin-facing snapshot of one channel: Name, Service, Config, LogLevel,
	// Parent and Module. Returns a not-found error for an unknown channel name.
	Describe(ctx core.ServerContext, chanName string) (utils.StringMap, error)
}
