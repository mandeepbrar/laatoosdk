package elements

import (
	"laatoo.io/sdk/server/core"
)

// NamespaceManager is the server element that OWNS THE NAMESPACE TREE: it discovers namespaces from
// configuration, holds them, resolves them by name, and drives their lifecycle.
//
// It exists because the role was homeless. Every other element kind has a manager that creates its
// kind from configuration and hands instances back; namespaces had none, so the descent from a
// solution to its applications and from an application to its isolations was written out by hand,
// twice, at two fixed depths. That is not a style problem -- it IS the depth ceiling, because a
// fourth level needs a third hand-written descent. One manager owning one recursion removes the
// ceiling rather than raising it.
//
// Discovery reads a namespace's own children-directory attribute (childrendir, default "children"),
// walks it, and recurses. The server does not interpret that directory's NAME: whether a deployment
// calls its children applications, environments or tenants is a label for humans and for the
// designer, and no server code names any of them.
//
// ONE ASYMMETRY, DELIBERATE AND WORTH KNOWING. Every other manager's kind appears as a segment in
// the addresses of the elements it holds -- "::myapp::scriptmanager::validate" -- and that segment
// is load-bearing, because it is what keeps a script and a service both named "validate" apart.
// Namespace addresses carry no such segment: a namespace is addressed "::myapp::uat", directly
// beneath its parent. They must be, since namespaces nest into themselves and a segment per level
// would yield "::namespacemanager::myapp::namespacemanager::uat".
//
// This is consistent rather than special-cased: an element's address is independent of which object
// holds it, so a manager can hold a namespace without appearing in its address. The consequence to
// accept is that a namespace named "scriptmanager" collides with the script manager under the same
// parent -- a genuine collision, refused at registration naming both parties, which is the right
// outcome rather than a gap.
type NamespaceManager interface {
	core.ServerElement

	// GetCurrentNamespace returns the namespace the CALLER is executing in, taken from ctx rather
	// than from any state this manager holds.
	//
	// Scope arrives at the call. One manager serves every namespace, so a handle obtained in one
	// scope and invoked with another context answers for the caller's scope, not the one that
	// obtained it. A manager that remembered "its" namespace would reintroduce the captured-context
	// defect this design exists to remove.
	GetCurrentNamespace(ctx core.ServerContext) Namespace

	// GetNamespace resolves a namespace by reference from the caller's scope, or returns nil.
	//
	// A BARE name resolves nearest-first -- the calling namespace, then each enclosing one out to
	// the root -- so a nearer declaration shadows an enclosing one. A QUALIFIED address
	// ("::portal::uat") binds exactly where it says and never falls back, which is how a caller
	// says which one it means when two namespaces share a name.
	//
	// Resolution is bounded by containment: a caller reaches its own subtree and its ancestors, and
	// nothing else. A namespace in a sibling subtree returns nil, and returns nil identically to a
	// name that was never declared -- a caller is not told whether an unreachable namespace exists.
	GetNamespace(ctx core.ServerContext, reference string) Namespace
}
