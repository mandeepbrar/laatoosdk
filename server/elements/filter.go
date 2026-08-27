package elements

import (
	"laatoo.io/sdk/server/core"
)

// Filter decides which of a parent level's named things a child level inherits.
//
// A child level (solution under server, application under solution, isolation under application)
// is built with a COPY of its parent's stores, and each entry is offered to every supplied filter
// first -- one veto excludes it. The same interface gates services, channels, factories, agents,
// topics, caches, secrets, data services and registered objects.
//
// NO CALL SITE SUPPLIES FILTERS TODAY. Every child-manager constructor accepts them as variadic
// arguments and every one is called with none, so a child currently inherits its parent's entire
// store. This is an extension point that is wired but unused; treat "no filters" as the live
// behaviour rather than an edge case.
type Filter interface {
	// Allowed reports whether the named thing may exist at the child level. objectName is the
	// entry's KEY in the store being copied -- a service alias, a channel name, a topic name --
	// not a Go type name, despite the parameter's name.
	//
	// Returning false excludes the entry from the child; the parent keeps it either way. A filter
	// is expected to answer for names it knows nothing about, so it must pick a deliberate default
	// rather than assuming its own list is exhaustive.
	Allowed(ctx core.ServerContext, objectName string) bool
}
