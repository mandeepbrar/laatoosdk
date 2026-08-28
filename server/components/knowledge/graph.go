package knowledge

import (
	"laatoo.io/sdk/server/core"
)

// Graph is the contract a plugin satisfies to be resolved through the KnowledgeManager server
// element, in place of a per-consumer alias string plus a locally-declared structural interface.
// Two Go plugins are separate package main binaries and cannot share a locally-declared interface
// type across the boundary; a structural assertion against one holds only while every type in the
// signature comes from the shared SDK, so both methods here use only core and builtin types.
type Graph interface {
	// RunSPARQL executes query against the graph and returns the result set, serialized in
	// whatever form the implementation's SPARQL engine produces.
	RunSPARQL(ctx core.RequestContext, query string) (string, error)

	// Reload rebuilds the graph from its underlying corpus. An implementation is expected to do
	// this atomically with respect to concurrent RunSPARQL calls, so a reader never observes a
	// partially-reloaded graph and a failed reload leaves the previous graph serving.
	Reload(ctx core.RequestContext) error
}

// GraphType discriminates the storage backend a Graph is built on -- for example "memory" or
// "laatoo" -- so a new backend can register under its own type instead of a consumer switching on
// a fixed set of hardcoded names.
type GraphType string

// GraphProvider builds a Graph for one registered GraphType. A plugin implementing a new storage
// backend registers a GraphProvider with the KnowledgeManager under that backend's GraphType,
// which is what lets a graph declaring the type be built without editing a closed switch.
type GraphProvider interface {
	// NewGraph builds a Graph backed by this provider, from settings taken from the caller's own
	// module configuration.
	NewGraph(ctx core.ServerContext, settings map[string]interface{}) (Graph, error)
}
