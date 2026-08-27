package elements

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
)

// ServerElementHandle is the WRITABLE side of a server element.
//
// Every manager is built as a pair: this handle, which the server keeps private, and a read-only
// proxy (ServiceManager, ChannelManager, ...) which is what gets published into contexts and handed
// to plugins. Holding the proxy therefore cannot start or reconfigure an element -- only the server
// that created it holds the handle.
//
// Both methods are ONE-SHOT and are not idempotent: the server calls each exactly once per element,
// in a fixed order across all elements (see laatooserver/src/core/abstractserver_initialize.go and
// abstractserver_start.go). Nothing guards against a second call.
type ServerElementHandle interface {
	// Initialize configures the element from its own section of the level's configuration -- the
	// server passes an empty config rather than skipping the call when the section is absent, so
	// an element must tolerate a config with nothing in it.
	//
	// This is where an element loads its declarations: services, channels, factories, tasks and so
	// on are created here, from modules and from the config directory, but not yet started. An
	// error returned here aborts the boot of the level.
	Initialize(ctx core.ServerContext, conf config.Config) error

	// Start activates what Initialize created -- binds routes, subscribes queues, runs each
	// component's Start hook. It runs after EVERY element at the level has been initialized, which
	// is what lets one element resolve another during start (a channel resolving its service, a
	// task queue resolving its processor). An error returned here aborts the boot of the level.
	Start(ctx core.ServerContext) error
}

/*
type Server interface {
	core.ServerElement
}*/
