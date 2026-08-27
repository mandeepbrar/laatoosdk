package knowledge

import (
	"iter"

	"laatoo.io/sdk/server/core"
)

// KnowledgeStore is implemented by a plugin that wants to act as storage for an
// RDF graph. A consuming plugin resolves one through the context and hands it to
// whatever RDF engine it embeds, so the engine's triples persist wherever this
// implementation puts them.
//
// # Graph names
//
// Every method that takes a graph Term follows one convention, and an
// implementation that departs from it will appear to work and lose data:
//
//   - The unset Term and a blank-node Term both mean the default graph, per
//     Term.IsDefaultGraph. Neither means "every graph": Len counts the default
//     graph, Triples reads it, and Remove deletes from it.
//   - Contexts reports named graphs only. The default graph is not one of them.
//
// # Errors
//
// Every operation returns an error, including the reads. This is a deliberate
// divergence from the in-memory RDF libraries these types are shaped against,
// whose equivalents return nothing at all: an in-memory map cannot fail, but a
// store backed by a disk, a database or a network genuinely can. An interface
// that omitted the error would force every such implementation to swallow its
// own failures silently, which is the one outcome worth designing against.
//
// A caller bridging this interface to a library that has nowhere to return an
// error is the party that decides what to do with one, and should do it in a
// single place rather than at each call site.
//
// # Concurrency
//
// A Laatoo service is a singleton whose Invoke is entered concurrently, so an
// implementation must assume concurrent calls from multiple goroutines and
// serialize its own state. Set exists because that guarantee cannot be
// assembled from Remove followed by Add by a caller.
//
// # Iteration
//
// Iteration order is unspecified and need not be stable between calls. The
// iterator-returning methods yield a non-nil error to report failure and stop;
// a caller must check it on every step rather than only at the end, because a
// failure part-way through otherwise reads as a shorter result set.
type KnowledgeStore interface {
	// Add inserts one statement into the named graph.
	Add(ctx core.RequestContext, triple Triple, graph Term) error

	// AddN inserts many statements, each carrying its own graph. It exists so
	// that an implementation can commit a bulk load in one transaction; a
	// caller loading a corpus should prefer it over repeated Add.
	AddN(ctx core.RequestContext, quads []Quad) error

	// Remove deletes every statement in the named graph matching the pattern.
	//
	// The zero TriplePattern matches everything, so it empties the graph.
	Remove(ctx core.RequestContext, pattern TriplePattern, graph Term) error

	// Set atomically replaces the object of a statement: it removes everything
	// matching (subject, predicate, *) and adds the given triple, without any
	// concurrent reader observing the interval in between where the old value
	// is gone and the new one is not yet present.
	Set(ctx core.RequestContext, triple Triple, graph Term) error

	// Triples iterates the statements in the named graph matching the pattern.
	Triples(ctx core.RequestContext, pattern TriplePattern, graph Term) iter.Seq2[Triple, error]

	// Len returns the number of statements in the named graph.
	Len(ctx core.RequestContext, graph Term) (int, error)

	// Contexts iterates the named graphs, optionally only those containing the
	// given statement. Pass nil to iterate them all. The default graph is never
	// reported.
	Contexts(ctx core.RequestContext, triple *Triple) iter.Seq2[Term, error]

	// Bind associates a prefix with the namespace IRI it abbreviates.
	Bind(ctx core.RequestContext, prefix string, namespace IRI) error

	// Namespace returns the namespace bound to a prefix. The boolean reports
	// whether a binding exists, and is distinct from the error, which reports
	// that the store could not be consulted at all.
	Namespace(ctx core.RequestContext, prefix string) (IRI, bool, error)

	// Prefix returns the prefix bound to a namespace — the inverse of
	// Namespace, with the same distinction between the boolean and the error.
	Prefix(ctx core.RequestContext, namespace IRI) (string, bool, error)

	// Namespaces iterates every prefix-to-namespace binding in the store.
	Namespaces(ctx core.RequestContext) iter.Seq2[NamespaceBinding, error]

	// ContextAware reports whether this store can hold named graphs. An
	// implementation returning false is expected to treat every graph name as
	// the default graph rather than to fail.
	//
	// It takes no context and returns no error because it is a static property
	// of the implementation, not an operation on it: a caller asks it while
	// wiring the store up, when there may be no request in flight to pass.
	ContextAware() bool

	// TransactionAware reports whether this store implements its operations
	// transactionally. See ContextAware for why it takes no context.
	TransactionAware() bool
}
