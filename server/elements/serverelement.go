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
// in a declared order. That order is data, not code: each phase declares its own traversal
// direction and its own ordering of element kinds, and the two differ -- initialize and start do
// NOT share one sequence. Nothing guards against a second call.
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
	// component's Start hook. It runs after EVERY element in the tree has been initialized, which
	// is what lets one element resolve another during start (a channel resolving its service, a
	// task queue resolving its processor). An error returned here aborts the boot.
	//
	// INITIALIZE DECLARES; START RESOLVES. Resolution scope widens between the two phases, and
	// this is the rule that makes it safe for a parent to be created and initialized before its
	// children exist:
	//
	//   - During Initialize, an element may resolve its ANCESTORS ONLY. They are guaranteed to be
	//     initialized already, because a parent is always initialized before its children.
	//     Anything else -- a sibling namespace, a descendant -- is not yet guaranteed to exist,
	//     and reaching for it is refused rather than left to boot order.
	//   - During Start, an element may resolve ANYTHING CONTAINMENT PERMITS, because by then every
	//     element in the tree has been initialized.
	//
	// Containment and phase are INDEPENDENT axes: containment decides what may be reached at all,
	// phase decides when a permitted reference is guaranteed to resolve. A reference that is
	// legal by containment can still be too early.
	//
	// This sentence formerly read "every element at the level", which was the same rule stated at
	// the scope of a three-level hierarchy that no longer exists. The rule did not change; the
	// tree it describes did.
	Start(ctx core.ServerContext) error
}

/*
type Server interface {
	core.ServerElement
}*/
