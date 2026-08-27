// Package knowledge declares the RDF value types and the KnowledgeStore
// interface a plugin implements to act as storage for an RDF graph.
//
// The types here are deliberately the SDK's own rather than any RDF library's.
// Two independent reasons force that:
//
//   - A library's term type is typically a sealed interface — goRDFlib's
//     term.Term carries an unexported termType() marker — so no package outside
//     it can implement one. Reusing it is not merely undesirable, it is
//     impossible from here.
//   - Importing an RDF library into the SDK would put that dependency into
//     every plugin that pins laatoo.io/sdk (~130 go.mod files at the time of
//     writing) whether or not it touches RDF. The SDK has no external
//     dependencies at all, and this must not be the change that gives it one.
//
// A plugin consuming an RDF library therefore owns a thin adapter translating
// these types to and from the library's. That translation is the deliberate
// cost of keeping the dependency out of the SDK.
package knowledge

// IRI is an RDF resource identifier: an absolute IRI, unabbreviated.
//
// It is a distinct type rather than a bare string so that namespace bindings
// and the predicate position cannot silently accept a prefixed form such as
// "rdfs:label". Expansion belongs to the caller, before the value gets here.
type IRI string

// TermKind discriminates the three kinds of RDF term a Term can hold.
type TermKind uint8

const (
	// TermUnset is the zero value of TermKind and is not a valid term.
	//
	// It is what an uninitialised Term carries, and it exists so that such a
	// Term is detectably invalid rather than being mistaken for a valid one.
	// Where a graph name is expected, it carries the additional meaning
	// described on KnowledgeStore: the default graph.
	TermUnset TermKind = iota

	// TermIRI is a resource identifier. Value holds the IRI.
	TermIRI

	// TermBlank is a blank node. Value holds its local identifier, without the
	// "_:" prefix used by N-Triples.
	TermBlank

	// TermLiteral is a literal. Value holds the lexical form, and Datatype,
	// Language and Direction qualify it.
	TermLiteral
)

// Term is one RDF term: an IRI, a blank node, or a literal.
//
// It is a struct with a kind discriminator rather than an interface, because
// an interface here would be implementable by callers and every implementation
// would then have to be understood by every store — including stores that must
// serialize what they are given. A comparable struct is directly usable as a
// map key and has one obvious on-disk shape.
//
// The Datatype, Language and Direction fields are meaningful only when Kind is
// TermLiteral and must be empty otherwise.
type Term struct {
	// Kind says which of the three RDF term kinds this is.
	Kind TermKind

	// Value is the IRI, the blank node identifier, or the literal's lexical
	// form, according to Kind.
	Value string

	// Datatype is the literal's datatype IRI. Empty means the literal is
	// either language-tagged or a plain xsd:string; it is never inferred here.
	Datatype IRI

	// Language is the literal's BCP 47 language tag, lowercased. Mutually
	// exclusive with a Datatype other than rdf:langString.
	Language string

	// Direction is the base direction of a language-tagged literal: "", "ltr"
	// or "rtl" (RDF 1.2). Empty means unspecified.
	Direction string
}

// NewIRITerm returns a Term holding the given resource identifier.
func NewIRITerm(iri IRI) Term {
	return Term{Kind: TermIRI, Value: string(iri)}
}

// NewBlankTerm returns a Term holding a blank node with the given local
// identifier. The identifier is taken as-is, without the "_:" prefix.
func NewBlankTerm(id string) Term {
	return Term{Kind: TermBlank, Value: id}
}

// NewLiteralTerm returns a plain or typed literal with the given lexical form.
// Pass an empty datatype for a literal that carries none.
func NewLiteralTerm(lexical string, datatype IRI) Term {
	return Term{Kind: TermLiteral, Value: lexical, Datatype: datatype}
}

// NewLangLiteralTerm returns a language-tagged literal. Direction may be "",
// "ltr" or "rtl"; anything else is the caller's error and is not validated here.
func NewLangLiteralTerm(lexical, language, direction string) Term {
	return Term{Kind: TermLiteral, Value: lexical, Language: language, Direction: direction}
}

// IsZero reports whether the Term is the unset zero value.
//
// Where a graph name is expected this means the default graph; anywhere else it
// means the Term was never initialised.
func (t Term) IsZero() bool {
	return t.Kind == TermUnset
}

// IsDefaultGraph reports whether this Term, used as a graph name, denotes the
// default graph.
//
// Both the unset Term and a blank node do. The blank node case is not an
// oddity: a blank node cannot name a graph in any serialization or over any
// protocol, and an RDF library handling an unnamed graph typically passes its
// own blank identifier straight through, so folding it into the default graph
// is what makes an unnamed graph behave identically on every backend.
func (t Term) IsDefaultGraph() bool {
	return t.Kind == TermUnset || t.Kind == TermBlank
}

// Triple is an RDF statement.
//
// Subject must be an IRI or a blank node, Predicate must be an IRI, and Object
// may be any kind of term. The SDK does not enforce this — a store is not the
// right place to discover a malformed statement — but a store may reject one.
type Triple struct {
	Subject   Term
	Predicate Term
	Object    Term
}

// Quad is a Triple together with the graph it belongs to.
//
// A zero or blank-node Graph means the default graph, as described on
// Term.IsDefaultGraph.
type Quad struct {
	Triple
	Graph Term
}

// TriplePattern matches triples by position. A nil field is a wildcard matching
// any term in that position, so the zero TriplePattern matches every triple.
//
// The fields are pointers rather than values precisely so that a wildcard has
// to be written as one. With a zero Term standing in for "any", a Term that was
// accidentally left uninitialised would silently widen the pattern — which on
// KnowledgeStore.Remove is the difference between deleting one statement and
// deleting the graph.
type TriplePattern struct {
	Subject   *Term
	Predicate *Term
	Object    *Term
}

// NamespaceBinding associates a short prefix with the namespace IRI it
// abbreviates.
type NamespaceBinding struct {
	Prefix    string
	Namespace IRI
}
