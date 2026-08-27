package elements

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
)

// Channel is one node in an engine's routing tree. A channel either binds a service to a path
// (a leaf that serves requests) or groups other channels under a path prefix (a group channel,
// declared by omitting `service:`).
//
// Every channel is created as a CHILD of an existing channel, from the root channel an engine
// exposes. There is one implementation per engine -- http, websocket, grpc, rpc, mcp -- and they
// are NOT interchangeable: Destruct type-asserts its parent to its own engine's concrete type.
type Channel interface {
	core.ServerElement

	// GetServiceName returns the service alias this channel invokes, or "" for a group channel.
	// The channel manager uses "" to decide a channel is a routing container and skips Serve
	// entirely for it.
	GetServiceName() string

	// Serve binds this channel's route to its service and starts accepting requests.
	//
	// Called by the channel manager once the service is resolved and a response handler has been
	// placed in ctx; it is not meant to be called by plugin code. On the http engine it returns nil
	// WITHOUT binding any route when the channel config sets `disableroute: true` -- a configured
	// endpoint then simply does not exist, with no error anywhere.
	Serve(ctx core.ServerContext) error

	// Child creates a channel under this one from channelConfig and returns it.
	//
	// On the http engine the presence of `service:` decides the shape: with a service the child
	// reuses this channel's router and registers a route; without one it creates a router GROUP at
	// `path` that further children hang off.
	//
	// A MISSING `path:` IS NOT REPORTED. The http implementation builds a BadConf error for it and
	// then discards the value (engine/http/httpchannel.go:199-202), so the child is created with an
	// empty path and inherits its parent's. Verify `path:` by reading the boot log, not by
	// expecting an error.
	Child(ctx core.ServerContext, name string, channelConfig config.Config, module core.Module) (Channel, error)

	// Destruct removes this channel's route from its parent, during module unload or hot-reload.
	//
	// parentChannel MUST be a channel of the SAME engine: every implementation type-asserts it to
	// its own concrete proxy type without checking, so passing a channel from another engine
	// panics. On the MCP engine Destruct is a no-op that returns nil -- unloading an MCP module
	// leaves its tools registered.
	Destruct(ctx core.ServerContext, parentChannel Channel) error

	// GetModule returns the module that declared this channel, or nil for a channel declared in
	// server or application configuration rather than inside a plugin. The nil is routine, not
	// exceptional, and callers must handle it.
	GetModule() core.Module

	// GetParent returns the channel this one was created under, or nil for an engine's root
	// channel.
	GetParent() Channel

	// GetEngine returns the engine hosting this channel. The ctx argument is not consulted by any
	// implementation.
	GetEngine(ctx core.ServerContext) Engine

	// GetDescription returns the channel's `description:` config value. When it is absent the
	// channel manager fills it in from the bound service's description at Serve time, so this can
	// read empty before the channel has been served.
	GetDescription(ctx core.ServerContext) string
}
