package elements

import (
	"laatoo.io/sdk/server/core"
)

// Filter is the POLICY layer over element resolution, applied at the single point where every
// name is resolved.
//
// It sits ABOVE a structural floor that it does not supply. Every element in the server lives in
// one index under its fully qualified address, and every bare walk, every qualified bind and every
// GetChild descent is one lookup against it. That lookup enforces containment unconditionally: a
// caller reaches its own subtree and its ancestors, and nothing else. A Filter is consulted only
// for entries that already passed.
//
// THE FLOOR IS NOT DELEGATED TO POLICY, and that is the whole design decision here. If the index
// defaulted to open and relied on a configured Filter for the boundary, a namespace would not be a
// boundary until somebody configured one -- which is exactly how this interface came to be dead
// code the first time. Policy narrows further, or widens deliberately; it never supplies the
// boundary.
//
// This interface previously described something else entirely: which of a parent LEVEL's named
// things a child level inherited from a COPY of the parent's stores. There are no levels and no
// copies now -- one element is registered once, under one address, and resolution walks rather
// than inherits.
type Filter interface {
	// Allowed reports whether the resolving caller may bind this entry. objectName is the entry's
	// KEY in the index, which is its fully qualified address -- "myapp::scriptmanager::validate",
	// not a Go type name, despite the parameter's name. ctx carries the caller's own scope, so a
	// filter can decide on the relationship between the two rather than on the target alone.
	//
	// Returning false makes the entry unresolvable for THIS caller; it stays registered and other
	// callers are unaffected. A filter is expected to answer for names it knows nothing about, so
	// it must pick a deliberate default rather than assuming its own list is exhaustive -- and
	// since containment is already enforced beneath it, the safe default is true.
	Allowed(ctx core.ServerContext, objectName string) bool
}
