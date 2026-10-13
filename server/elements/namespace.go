package elements

import (
	"laatoo.io/sdk/server/core"
)

// Namespace is a scope in the server's containment tree — the unit that replaces the Solution,
// Application and Isolation levels.
//
// Unlike the interfaces it replaces, it carries real structure. Application declared no methods of
// its own and existed for years providing nothing; a namespace IS the hierarchy, so a handle to
// one must be able to walk it. Nesting depth is unbounded: the server encodes no level vocabulary,
// and what a deployment calls each depth (application, environment, tenant) is a directory name it
// chooses, not a level the server implements.
//
// This is the READ side. Lifecycle lives on NamespaceHandle, mirroring the split
// ServerElementHandle already documents — a plugin that resolves a namespace can inspect it and
// walk it, but cannot start or stop it.
type Namespace interface {
	core.ServerElement

	// Walking UP is core.ServerElement.GetParent(), not a method here.
	//
	// This interface deliberately does NOT redeclare GetParent with a Namespace return. Go permits
	// an embedded and a locally declared method to share a name only when the signatures are
	// IDENTICAL, so `GetParent() Namespace` alongside the embedded `GetParent() ServerElement`
	// does not compile — and adding a second differently-named method for one relationship would
	// be worse than the assertion it saves.
	//
	// The assertion is total: a namespace's parent is always a namespace, because namespaces nest
	// directly within one another. So `ns.GetParent().(Namespace)` is safe wherever
	// `ns.GetParent() != nil`, and the nil-interface root test documented on ServerElement applies
	// here unchanged.

	// GetRoot returns the outermost namespace, which is the one the solution directory defines.
	// Called on the root itself it returns the root.
	GetRoot() Namespace

	// GetChildren returns the namespaces nested directly beneath this one, in no defined order.
	// A leaf returns an empty slice, never nil.
	//
	// This is NARROWER than core.ServerElement.GetChild, deliberately. A namespace's element
	// children include its managers -- GetChild(ctx, "scriptmanager") resolves one -- whereas this
	// returns only the children that are themselves namespaces. Use this to walk the namespace
	// tree, GetChild to descend an address.
	GetChildren() []Namespace

	// GetPath returns this namespace's fully qualified scope path in `::` form — "portal::uat"
	// for a namespace two deep. The root's path is "::".
	//
	// This is the form a declaration writes to bind a name exactly rather than let it resolve
	// nearest-first: `service: portal::orders` binds in portal, `service: orders` walks outward.
	GetPath() string

	// The three lifecycle contexts, retained rather than discarded.
	//
	// Every namespace runs the same cycle, and the server already derives a distinct context for
	// each phase — it simply drops them today. Keeping them makes the lifecycle inspectable after
	// the fact, which is what diagnostics need and what nothing can currently answer.
	//
	// Distinct from ServerElement.GetContext(), which returns the namespace's OWN context — the
	// one its elements resolve in. These three are snapshots of how it was brought up.

	// GetCreationContext returns the context the namespace was created under, before any of its
	// configuration was read.
	GetCreationContext() core.ServerContext

	// GetInitializationContext returns the context its Initialize ran under — the phase in which
	// its elements were created from configuration but not yet activated.
	GetInitializationContext() core.ServerContext

	// GetStartContext returns the context its Start ran under. Nil until Start has been called,
	// which is also how a caller can tell a namespace is created but not yet running.
	GetStartContext() core.ServerContext
}

// NamespaceHandle is the WRITABLE side of a namespace, held by the server and never handed to a
// plugin. Lifecycle is embedded from ServerElementHandle rather than redeclared, so a namespace
// obeys exactly the contract every other server element does — including the rule that Start runs
// only after every element at the scope has been initialized.
type NamespaceHandle interface {
	Namespace
	ServerElementHandle
}
