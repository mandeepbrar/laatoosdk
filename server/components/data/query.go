package data

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// QueryVersion identifies the AST revision a Query was built against.
//
// There is one revision, and adding a second was considered and rejected when traversal,
// projection and result-type navigation were added to this file — so the reasoning is recorded
// here rather than rediscovered.
//
// A version reconciles an old READER meeting new DATA. Neither half of that exists here. A
// Query is never serialized: it carries no codec tags, has no ReadAll/WriteAll, and never
// leaves the process. And its producer and consumer cannot disagree, because every data
// provider is a Go plugin — the runtime refuses to load a .so built against a different SDK,
// which is why an SDK change is a flag day in the first place. A provider compiled before a
// construct existed can therefore never be handed that construct.
//
// What negotiates instead is QueryCapability, per construct, and that is deliberately stronger
// than a version: a monotonic number says "newer", while a capability says WHICH shape a
// provider cannot compile — and refusing the specific shape is what stops a hop executor, which
// walks a tree, from accepting a cyclic pattern and answering it wrongly rather than declining.
// A single monotonic flag is the design this file's capability vocabulary exists to avoid.
//
// So the field earns its place for one reason only, and it is the reason QueryVersionInvalid is
// the zero value: it catches a Query built without a constructor. NOTHING ELSE MAY GATE ON IT.
// A provider deciding what it can execute asks SupportsQuery, never this field.
//
// Note that the guard is presently aspirational: QueryVersionInvalid is asserted by tests and
// cited in a comment, but no execution path rejects a query carrying it, so a hand-built
// Query{Filter: ...} runs everywhere. ValidateQuery is where the check belongs, and adding it
// needs an audit of the hand-built literals in the conformance and unit suites first.
type QueryVersion int

const (
	// QueryVersionInvalid is the zero value and is never a valid version. Every constant
	// group in this file starts one past zero deliberately: the condition framework this
	// replaces made its first constant the zero value, so a caller that forgot to set the
	// field silently selected a real behaviour.
	QueryVersionInvalid QueryVersion = iota
	// QueryV1 is the only revision: filters, parameters and extensions, and — added later
	// without a new constant, for the reasons above — traversal, projection and navigation.
	QueryV1
)

// PredicateKind reports which concrete node a Predicate is, for compilers that prefer a
// value switch over a type switch.
type PredicateKind string

const (
	// KindComparison is a field compared against an operand.
	KindComparison PredicateKind = "comparison"
	// KindLogical is a conjunction or disjunction of other predicates.
	KindLogical PredicateKind = "logical"
	// KindNot is the negation of a single predicate.
	KindNot PredicateKind = "not"
	// KindMembership tests a field against a set of operands.
	KindMembership PredicateKind = "membership"
	// KindNull tests whether a field is absent.
	KindNull PredicateKind = "null"
	// KindFunction applies a filter function to a field.
	KindFunction PredicateKind = "function"
	// KindExtension carries a predicate outside the OData grammar.
	KindExtension PredicateKind = "extension"
	// KindTraversal evaluates a predicate against the entities a navigation path reaches
	// rather than against the record being filtered.
	KindTraversal PredicateKind = "traversal"
)

// CompareOperator is a binary comparison. The values are the OData filter grammar's own
// tokens, so text lowering is a direct mapping rather than a translation table.
type CompareOperator string

const (
	OpEqual        CompareOperator = "eq"
	OpNotEqual     CompareOperator = "ne"
	OpGreater      CompareOperator = "gt"
	OpGreaterEqual CompareOperator = "ge"
	OpLess         CompareOperator = "lt"
	OpLessEqual    CompareOperator = "le"
)

// LogicalOperator combines predicates.
type LogicalOperator string

const (
	LogicalAnd LogicalOperator = "and"
	LogicalOr  LogicalOperator = "or"
)

// FilterFunction is a string function from the OData filter grammar.
type FilterFunction string

const (
	FuncContains   FilterFunction = "contains"
	FuncStartsWith FilterFunction = "startswith"
	FuncEndsWith   FilterFunction = "endswith"
	// FuncReachable tests whether the record can reach an entity by repeatedly following
	// the navigation property the call names — its Field is that property and its first
	// argument identifies the target, with an optional second argument bounding the depth.
	//
	// It is a filter function rather than an Extension payload because OData already has
	// function-call syntax for filters and this is a predicate over a SINGLE binding. The
	// test for the extension seam is whether a construct needs more than one free binding;
	// reachability does not, so it stays inside the grammar.
	//
	// Note for a provider: the existing string functions are gated as one capability, and
	// this one is deliberately NOT part of it. Declaring CapabilityStringFunctions says
	// nothing about reachability, which has its own CapabilityReachability.
	FuncReachable FilterFunction = "reachable"
)

// OperandKind distinguishes a fixed value from one supplied at execution.
type OperandKind string

const (
	// OperandLiteral is a value fixed when the query was built.
	OperandLiteral OperandKind = "literal"
	// OperandParameter names a value bound at execution. Parameters are never inlined at
	// build time — that is what allows a query to be compiled once and bound per request.
	OperandParameter OperandKind = "parameter"
	// OperandField names another field on the same entity.
	OperandField OperandKind = "field"
	// OperandPath names a value reached by following one or more navigation properties from
	// the entity currently in scope — OData's `Owner/Name`. It is distinct from OperandField,
	// which names a field on the entity itself and cannot cross a reference.
	OperandPath OperandKind = "path"
)

// Operand is one side of a comparison.
type Operand struct {
	Kind OperandKind
	// Value holds the operand when Kind is OperandLiteral.
	Value interface{}
	// Name holds the parameter or field name when Kind is OperandParameter or OperandField.
	Name string
	// Path holds the navigation segments when Kind is OperandPath, outermost first: every
	// segment but the last names a reference field, and the last names a scalar on the
	// entity they reach. It is a slice rather than a slash-joined string so that a compiler
	// never re-parses what the builder already knew, and so that a field name containing the
	// separator is not ambiguous.
	Path []string
}

// LiteralOperand builds an operand from a fixed value.
func LiteralOperand(value interface{}) Operand {
	return Operand{Kind: OperandLiteral, Value: value}
}

// ParameterOperand builds an operand bound at execution from the named parameter.
func ParameterOperand(name string) Operand {
	return Operand{Kind: OperandParameter, Name: name}
}

// FieldOperand builds an operand referring to another field on the same entity.
func FieldOperand(name string) Operand {
	return Operand{Kind: OperandField, Name: name}
}

// PathOperand builds an operand naming a value reached across one or more navigation
// properties. Called with a single segment it is still a path and not a field: the caller has
// said it crosses a reference, and a provider that cannot cross one must refuse rather than
// read a same-entity field of that name.
func PathOperand(segments ...string) Operand {
	return Operand{Kind: OperandPath, Path: segments}
}

// parameterName returns the parameter this operand depends on, or the empty string when it
// does not depend on one.
func (o Operand) parameterName() string {
	if o.Kind == OperandParameter {
		return o.Name
	}
	return ""
}

// Predicate is any node in a query's filter tree. Each concrete node carries only the fields
// that are meaningful for it — there is no single struct whose valid fields depend on a kind
// discriminator, because that ambiguity is what makes a representation hard to compile against.
type Predicate interface {
	// Kind reports which concrete node this is.
	Kind() PredicateKind
	// IsOptional reports whether this predicate is dropped when a parameter it references is
	// not supplied at execution.
	IsOptional() bool
	// Parameters lists every parameter name this node and its children depend on.
	Parameters() []string
}

// Optionality is embedded by every predicate node to carry the elision flag.
type Optionality struct {
	// Optional marks a predicate that is removed from the query when a parameter it
	// references is unbound, rather than compared against a zero value. A dataset may
	// therefore declare more filters than a caller supplies parameters for.
	Optional bool
}

// IsOptional reports whether this predicate is elided when its parameters are unbound.
func (o Optionality) IsOptional() bool { return o.Optional }

// Comparison compares a field against an operand.
type Comparison struct {
	Optionality
	Field    string
	Operator CompareOperator
	Value    Operand
}

// Kind reports that this is a comparison node.
func (c *Comparison) Kind() PredicateKind { return KindComparison }

// Parameters lists the parameter this comparison's operand depends on, if any.
func (c *Comparison) Parameters() []string {
	if name := c.Value.parameterName(); name != "" {
		return []string{name}
	}
	return nil
}

// Logical combines two or more predicates with a single operator.
type Logical struct {
	Optionality
	Operator LogicalOperator
	Operands []Predicate
}

// Kind reports that this is a logical node.
func (l *Logical) Kind() PredicateKind { return KindLogical }

// Parameters lists every parameter its operands depend on.
func (l *Logical) Parameters() []string {
	var names []string
	for _, operand := range l.Operands {
		names = append(names, operand.Parameters()...)
	}
	return names
}

// Not negates a single predicate.
type Not struct {
	Optionality
	Operand Predicate
}

// Kind reports that this is a negation node.
func (n *Not) Kind() PredicateKind { return KindNot }

// Parameters lists every parameter the negated predicate depends on.
func (n *Not) Parameters() []string {
	if n.Operand == nil {
		return nil
	}
	return n.Operand.Parameters()
}

// Membership tests a field against a set of operands.
type Membership struct {
	Optionality
	Field string
	// Values is the set tested against. A single operand of a list-valued parameter is
	// permitted, so variable-arity membership does not change the query's shape.
	Values []Operand
	// Negated inverts the test, expressing "not in".
	Negated bool
}

// Kind reports that this is a membership node.
func (m *Membership) Kind() PredicateKind { return KindMembership }

// Parameters lists the parameters this membership test's operands depend on.
func (m *Membership) Parameters() []string {
	var names []string
	for _, value := range m.Values {
		if name := value.parameterName(); name != "" {
			names = append(names, name)
		}
	}
	return names
}

// NullTest tests whether a field is absent.
type NullTest struct {
	Optionality
	Field string
	// Negated inverts the test, expressing "is not null".
	Negated bool
}

// Kind reports that this is a null test node.
func (n *NullTest) Kind() PredicateKind { return KindNull }

// Parameters reports that a null test never depends on a parameter.
func (n *NullTest) Parameters() []string { return nil }

// FunctionCall applies a filter function to a field.
type FunctionCall struct {
	Optionality
	Function FilterFunction
	Field    string
	// Arguments are the function's operands beyond the field it applies to.
	Arguments []Operand
}

// Kind reports that this is a function node.
func (f *FunctionCall) Kind() PredicateKind { return KindFunction }

// Parameters lists the parameters this call's arguments depend on.
func (f *FunctionCall) Parameters() []string {
	var names []string
	for _, argument := range f.Arguments {
		if name := argument.parameterName(); name != "" {
			names = append(names, name)
		}
	}
	return names
}

// Extension carries a predicate the OData grammar cannot express — an ancestor lookup or a
// vector search, for example. An extension is explicitly not OData-compliant, and a provider
// that does not recognise its namespace must reject it rather than ignore it.
type Extension struct {
	Optionality
	// Namespace identifies the vocabulary the extension belongs to, so two providers cannot
	// collide on a name.
	Namespace string
	Name      string
	Payload   interface{}
	// Params names any parameters the payload depends on, since the payload is opaque here.
	Params []string
}

// Kind reports that this is an extension node.
func (e *Extension) Kind() PredicateKind { return KindExtension }

// Parameters lists the parameters the extension declares a dependency on.
func (e *Extension) Parameters() []string { return e.Params }

// Quantifier says how a Traversal applies its sub-predicate to what the path reaches.
type Quantifier string

const (
	// QuantifierInvalid is the zero value and is never valid, following this file's rule
	// that a caller who forgot to set a discriminator must not silently select a behaviour.
	QuantifierInvalid Quantifier = ""
	// QuantifierSingle applies the predicate to the one entity a single-valued navigation
	// reaches — OData's plain path filter, `Category/Name eq 'Food'`. A record whose
	// reference is absent does not match, which is what separates it from QuantifierAll
	// over an empty collection.
	QuantifierSingle Quantifier = "single"
	// QuantifierAny matches when at least one entity the path reaches satisfies the
	// predicate — OData's `any` lambda, and the closest thing its grammar has to a join.
	QuantifierAny Quantifier = "any"
	// QuantifierAll matches when every entity the path reaches satisfies the predicate —
	// OData's `all` lambda. It is vacuously true for an empty collection, which is the one
	// place the two quantifiers disagree about a record with no related entities.
	QuantifierAll Quantifier = "all"
)

// DepthUnbounded is the MaxDepth of a variable-length path with no upper bound: repeat the
// last segment until the frontier is exhausted. Only a provider declaring
// CapabilityUnboundedPath may compile it — every other one refuses, rather than choosing a
// depth of its own, because a silently chosen depth returns a partial answer that looks whole.
const DepthUnbounded = -1

// Traversal evaluates a predicate against the entities a navigation path reaches, rather than
// against fields of the record being filtered.
//
// It is one node covering the plain path filter, both lambdas and the variable-length path,
// because OData spells all of them the same way — a path, optionally quantified — and because
// a provider compiling one needs the path and the quantifier together to choose between a
// join and a semi-join. The capabilities that gate it are still separate, so a store that can
// express `any` but not `all` refuses on the quantifier rather than on the node.
//
// Comparison.Field never crosses a reference. To filter on a field of a related entity —
// OData's `Category/Name eq 'Food'` — wrap the comparison in a QuantifierSingle Traversal
// whose Path is `["Category"]`; the comparison's Field is then resolved on Category.
//
// The result type is unchanged: a Traversal filters the records the query already targets. To
// return the related records instead, use Query.Navigate.
type Traversal struct {
	Optionality
	// Path is the navigation segments from the entity in scope, outermost first. Every
	// segment names a reference field, and Predicate is evaluated in the scope of the entity
	// the last one reaches.
	Path []string
	// Quantifier says how Predicate is applied to what Path reaches.
	Quantifier Quantifier
	// Predicate is evaluated against the entity at the end of Path, so its field names are
	// resolved there and not on the entity the query started from. It may itself be a
	// Traversal, which is what gives arbitrary-depth filtering without needing a free
	// pattern: nesting stays tree-shaped, and a tree is what a hop executor can walk.
	Predicate Predicate
	// MinDepth and MaxDepth bound a variable-length path by repeating the LAST segment of
	// Path between them. Both left at zero means a fixed-length path — the ordinary case —
	// so a builder that ignores these fields gets exactly the traversal it wrote. A MaxDepth
	// of DepthUnbounded has no upper bound.
	MinDepth int
	MaxDepth int
	// Relationship narrows the traversal to edges of one declared TYPE, and is what
	// `relationshipname` in an entity's field declaration exists to be matched on.
	//
	// Empty means every edge the Path reaches, which is what every traversal built before this
	// field existed meant -- so an unset Relationship changes nothing.
	//
	// It is matched against data.StorableRef.Relationship, which codegen stamps onto the
	// reference at save. That value is STORED rather than derived, so a provider may answer this
	// from an index; a struct tag could not, which is the whole reason the SDK carries it as a
	// field. Until this existed, `relationshipname` could be declared and stamped and then matched
	// on by nothing: Relationship appeared only in storageinfo.go across the entire SDK, so the
	// AST had no way to say "every KNOWS edge".
	//
	// Gated by CapabilityRelationshipMatch: a provider that stores the ref but cannot filter on
	// its Relationship must refuse rather than return every edge, which is a wider answer than
	// the query asked for.
	Relationship string
	// TargetScope carries how the entity this path REACHES must be restricted -- see
	// ScopeRequirement. Nil means the caller has not said, and a provider must then refuse to
	// compile the hop natively rather than emit an unrestricted join.
	TargetScope *ScopeRequirement
	// MatchOptional makes a path that reaches nothing non-fatal to the match: the record is
	// kept and the sub-predicate contributes nothing. This is SPARQL's OPTIONAL and a left
	// outer join, and it is the construct a nested OPTIONAL currently has no home for.
	//
	// It is deliberately NOT Optionality.Optional, which sits alongside it and is easy to
	// mistake it for. Optional removes the whole predicate when a PARAMETER it references is
	// unbound — it is about binding. MatchOptional keeps the predicate and relaxes what a
	// failure to TRAVERSE means — it is about data. A node may set either, both or neither,
	// and the two are never substitutes.
	MatchOptional bool
}

// ScopeRequirement states how the entity a traversal or expansion REACHES must be restricted, so
// a provider compiling the hop into its own store query applies what the target's own component
// would have applied.
//
// WHY THIS TRAVELS ON THE QUERY RATHER THAN BEING ASKED OF THE TARGET. The provider compiling a
// join holds the PARENT's component and never the child's -- resolving another entity's component
// is the DataManager's job, and a provider reaching for it on the query path is the coupling the
// hop executor exists to keep out. So the requirement is stamped onto the query by the caller,
// which can resolve components, and the provider reads it as data.
//
// WHY NOT READ THE TARGET'S StorableConfig INSTEAD. Because it is not the source of truth. A data
// service takes these settings from MODULE CONFIGURATION, falling back to the entity's own
// declaration only when the module supplies none -- which is exactly what the admin-portal
// pattern overrides. A provider reading the entity's flags would apply the default a deployment
// had deliberately changed, and would do it silently.
//
// AN UNSET REQUIREMENT MEANS UNKNOWN, NOT UNRESTRICTED. A provider handed a traversal whose
// TargetScope is nil has not been told the target needs no restriction -- it has not been told
// anything, and must refuse the native path rather than emit a join with none. Reading nil as "no
// restriction required" is a cross-tenant read that RETURNS ROWS, which is the failure this type
// exists to prevent.
type ScopeRequirement struct {
	// Multitenant says the target's rows are partitioned per tenant, so the hop must restrict to the
	// caller's tenant. The VALUE is not carried here: it belongs to the request, and a query is
	// compiled once and bound many times.
	Multitenant bool
	// TenantField names the column or property carrying the tenant, so a provider need not
	// hardcode the platform's own spelling of it.
	TenantField string
	// SoftDelete says the target marks deletion rather than removing the row, so the hop must
	// exclude the deleted ones.
	SoftDelete bool
	// DeletedField names the flag soft deletion sets.
	DeletedField string
}

// ScopedTraversal reports whether this traversal has been told how to restrict what it reaches.
// A provider compiling a hop natively must check this and refuse when it is false.
func (t *Traversal) ScopedTraversal() bool { return t.TargetScope != nil }

// Kind reports that this is a traversal node.
func (t *Traversal) Kind() PredicateKind { return KindTraversal }

// Parameters lists every parameter the traversed predicate depends on. The path itself is
// structural and never parameterised: which reference a query follows is part of its shape,
// not part of its binding.
func (t *Traversal) Parameters() []string {
	if t.Predicate == nil {
		return nil
	}
	return t.Predicate.Parameters()
}

// Expansion is one entry in a query's $expand projection: a reference field to resolve and
// attach to each returned record, and what to expand beneath it.
//
// Expansion is projection, not filtering. It changes what each returned record CARRIES and
// never which records are returned, which is why it is a sibling of Query.Filter rather than
// a Predicate — the one exception being Required, which says so explicitly.
type Expansion struct {
	// Field is the reference field, on the entity in scope, to resolve and attach.
	Field string
	// Levels repeats Field itself that many times, for a self-referencing hierarchy — this
	// is OData's $levels, and it applies to Field alone, not to Expand beneath it. Zero and
	// one both mean a single hop, so an unset Levels expands once. LevelsMax is the
	// unbounded form.
	Levels int
	// Expand nests further expansions beneath this one, resolved against the entity Field
	// reaches. This is how a multi-hop read is declared as one artefact: Service to
	// ObjectSpec to Config is two nested entries, not two queries.
	Expand []Expansion
	// Required drops records whose Field resolves to nothing, turning the outer join an
	// expansion performs by default into an inner one. The default — false — is the
	// behaviour expansion has always had, so an entry that ignores this field keeps it.
	Required bool

	// The options below are OData's expansion-scoped query options, and every one of them
	// applies to the EXPANDED SET — never to the records the query returns.
	//
	// That distinction is the whole reason they exist as separate fields, and getting it
	// backwards is a silent defect rather than an error. Given Region with a Nodes reference:
	//
	//	Expansion{Field: "Nodes", Filter: Status eq 'active'}
	//	    → EVERY region, each carrying only its active nodes. A region whose nodes are all
	//	      inactive is still returned, carrying an empty expansion.
	//	Query.Filter = Traversal{Path: ["Nodes"], Any, Status eq 'active'}
	//	    → only regions that HAVE an active node.
	//
	// The row counts differ and neither form errors, so a provider implementing one as the
	// other returns a plausible wrong answer. A provider that cannot scope an option to the
	// expanded set must decline the capability rather than apply it to the parents.

	// Filter constrains which related records are attached. It is a full predicate tree, so
	// it may itself contain a Traversal, and it may reference parameters — Query.Parameters
	// reports those alongside the main filter's, and Resolve prunes an optional one here
	// exactly as it does in Filter.
	Filter Predicate
	// Select limits which fields of the related records are read. Empty means all of them.
	Select []string
	// OrderBy orders the expanded set as FIELD/DIRECTION PAIRS — {"Name", "asc", "Id",
	// "desc"} — which is the shape the Get family's orderBy already uses and which providers
	// already walk two at a time. It is not a list of "Name asc" clauses; the declarative
	// front-end's single-string form is lowered into pairs before it reaches here.
	OrderBy []string
	// Top bounds how many related records are attached, and Skip offsets into them. Both are
	// zero when unset, and zero means unbounded rather than none — an expansion that attached
	// nothing by default would silently empty every existing query's projection.
	//
	// Top is the cheapest defence against expansion fan-out, which is the cost dimension a
	// bare expansion cannot bound at all.
	Top  int
	Skip int
	// Count asks for the total size of the expanded set alongside the attached page, so a
	// caller can tell a truncated expansion from a complete one. Without it, Top makes those
	// two indistinguishable.
	Count bool
}

// Parameters lists every parameter this expansion and its nested expansions depend on, so that
// an expansion filter participates in binding exactly as the query's own filter does. An
// expansion whose parameters were invisible here would be handed to a provider unresolved.
func (e Expansion) Parameters() []string {
	var names []string
	if e.Filter != nil {
		names = append(names, e.Filter.Parameters()...)
	}
	for _, nested := range e.Expand {
		names = append(names, nested.Parameters()...)
	}
	return names
}

// resolve prunes this expansion's own filter and its nested expansions against the supplied
// parameters, returning a copy. An expansion whose filter is elided keeps the expansion and
// drops only the constraint — unlike a Traversal, where the whole node goes: an expansion with
// no filter is a meaningful request (attach everything), whereas a traversal with no predicate
// asserts only that a path reaches something, which is not what the caller wrote.
func (e Expansion) resolve(params utils.StringsMap) Expansion {
	resolved := e
	resolved.Filter = nil
	if e.Filter != nil {
		resolved.Filter = resolvePredicate(e.Filter, params)
	}
	resolved.Expand = resolveExpansions(e.Expand, params)
	return resolved
}

// resolveExpansions prunes a list of expansions, returning nil for an empty input so that a
// query with no projection resolves to one with no projection rather than to an empty slice.
func resolveExpansions(expansions []Expansion, params utils.StringsMap) []Expansion {
	if len(expansions) == 0 {
		return nil
	}
	resolved := make([]Expansion, len(expansions))
	for i, expansion := range expansions {
		resolved[i] = expansion.resolve(params)
	}
	return resolved
}

// LevelsMax requests OData's `$levels=max`: expand until the graph is exhausted. Only a
// provider declaring CapabilityUnboundedExpand may compile it.
const LevelsMax = -1

// QueryCapability names something a provider may or may not be able to compile. It is
// separate from Feature, which describes datastore capabilities rather than query ones.
type QueryCapability string

const (
	// CapabilityComparison covers the ordered comparisons beyond equality.
	CapabilityComparison QueryCapability = "comparison"
	// CapabilityDisjunction covers logical or.
	CapabilityDisjunction QueryCapability = "disjunction"
	// CapabilityNegation covers not.
	CapabilityNegation QueryCapability = "negation"
	// CapabilityMembership covers set membership.
	CapabilityMembership QueryCapability = "membership"
	// CapabilityNullTest covers absence tests.
	CapabilityNullTest QueryCapability = "nulltest"
	// CapabilityStringFunctions covers contains, startswith and endswith.
	CapabilityStringFunctions QueryCapability = "stringfunctions"
	// CapabilityNesting covers arbitrarily nested logical grouping.
	CapabilityNesting QueryCapability = "nesting"

	// The capabilities below name SHAPE — traversal, projection and result type — and every
	// one of them is additive: none of the seven above is redefined, so a provider written
	// before they existed keeps meaning exactly what it meant.
	//
	// They are deliberately granular rather than one "supports graph" flag. A single flag is
	// unsafe here in a way it is not elsewhere: hop execution is a tree walk and cannot close
	// a cycle, so a document store accepting one boolean would accept a cyclic pattern and
	// answer it WRONGLY rather than decline it. Granularity is what keeps the failure mode
	// "refused" instead of "wrong".

	// CapabilityNavigationPath covers an operand or traversal that crosses a reference —
	// OData's `Category/Name`. What a provider declares here is whether it compiles the path
	// itself; one that does not can still be handed the traversal by a hop executor.
	CapabilityNavigationPath QueryCapability = "navigationpath"
	// CapabilityLambdaAny covers QuantifierAny over a navigation path.
	CapabilityLambdaAny QueryCapability = "lambdaany"
	// CapabilityLambdaAll covers QuantifierAll. It is separate from the `any` form because a
	// store that can express a semi-join need not be able to express its complement.
	CapabilityLambdaAll QueryCapability = "lambdaall"
	// CapabilityBoundedPath covers a variable-length path with a finite MaxDepth.
	CapabilityBoundedPath QueryCapability = "boundedpath"
	// CapabilityUnboundedPath covers a variable-length path with MaxDepth DepthUnbounded. It
	// is separate from the bounded form because a provider that unrolls a path to a fixed
	// depth serves one and not the other.
	CapabilityUnboundedPath QueryCapability = "unboundedpath"
	// CapabilityExpand covers a single level of the $expand projection.
	CapabilityExpand QueryCapability = "expand"
	// CapabilityNestedExpand covers expansions nested beneath other expansions.
	CapabilityNestedExpand QueryCapability = "nestedexpand"
	// CapabilityUnboundedExpand covers LevelsMax.
	CapabilityUnboundedExpand QueryCapability = "unboundedexpand"
	// CapabilityExpandFilter covers a filter scoped to the expanded set. It is separate from
	// CapabilityExpand, and separately dangerous: a provider that applies an expansion filter
	// to the PARENTS instead returns a plausible wrong answer with no error, since the two
	// forms differ only in row count. A provider that cannot scope the filter declines here
	// rather than approximating it.
	CapabilityExpandFilter QueryCapability = "expandfilter"
	// CapabilityExpandSelect covers limiting which fields of the expanded records are read.
	CapabilityExpandSelect QueryCapability = "expandselect"
	// CapabilityExpandOrderBy covers ordering the expanded set. It matters more in company
	// with paging than alone: a provider that honours Top while ignoring OrderBy returns an
	// arbitrary subset of the children rather than the requested one.
	CapabilityExpandOrderBy QueryCapability = "expandorderby"
	// CapabilityExpandPaging covers Top and Skip together, because every backend expresses
	// them in one clause and no provider has been observed to serve one without the other.
	// Ignoring Top merely over-fetches; ignoring Skip returns the wrong page.
	CapabilityExpandPaging QueryCapability = "expandpaging"
	// CapabilityExpandCount covers reporting the expanded set's full size alongside a page of
	// it, which is what distinguishes a truncated expansion from a complete one.
	CapabilityExpandCount QueryCapability = "expandcount"
	// CapabilityNavigate covers Query.Navigate — a query returning a related entity rather
	// than its own.
	CapabilityNavigate QueryCapability = "navigate"
	// CapabilityOptionalMatch covers Traversal.MatchOptional and Expansion.Required, the two
	// places a query says what reaching nothing should mean.
	CapabilityOptionalMatch QueryCapability = "optionalmatch"
	// CapabilityReachability covers FuncReachable. It is separate from
	// CapabilityStringFunctions, which gates the three OData string functions and must not be
	// read as covering every FunctionCall node.
	CapabilityReachability QueryCapability = "reachability"
	// CapabilityFreePattern covers a pattern with more than one free binding, carried as an
	// Extension under a graph namespace.
	CapabilityFreePattern QueryCapability = "freepattern"
	// CapabilityCyclicPattern covers a pattern that closes a cycle. It is separate from every
	// other capability here for the reason given above: an executor that wrongly accepts a
	// cycle does not fail, it answers.
	CapabilityCyclicPattern QueryCapability = "cyclicpattern"

	// CapabilityRelationshipMatch covers a Traversal that narrows to one edge TYPE via
	// Relationship. It is separate from CapabilityNavigationPath because following a reference
	// and filtering on the edge's declared type are different asks: every store that holds a
	// reference can follow it, while only one that stored and indexed StorableRef.Relationship
	// can filter on it.
	//
	// A provider must refuse rather than ignore the field. Ignoring it returns every edge the
	// path reaches instead of the ones of that type -- a WIDER answer than the query asked for,
	// which no caller can detect from the rows.
	CapabilityRelationshipMatch QueryCapability = "relationshipmatch"
)

// ExpandingComponent is implemented by a DataComponent that compiles the $expand projection
// itself, rather than leaving it to the hop executor. It is an OPTIONAL interface, asserted for
// and never required: DataComponent is not widened, so no existing provider is affected and
// R31's silent-breakage hazard does not arise.
//
// It REPLACES CompileQuery for a provider that implements it, rather than supplementing it, and
// that is the whole point of the shape. Two weaker designs were considered and rejected:
//
//   - Gating on SupportsQuery(CapabilityExpand). That is a CLAIM the provider makes about
//     itself. A provider may answer true and still ignore Query.Expand, returning unexpanded
//     records with a nil error — the silently-wrong direction, and one the capability's own
//     contract forbids without anything enforcing it.
//   - A separate CompileExpansion alongside CompileQuery. No better: a provider could satisfy
//     it and still ignore the projection when compiling. It also leaves the provider holding
//     two opaque artifacts while Get takes one.
//
// Because expansion arrives ONLY through this method, a provider that does not implement it
// never receives Expand at all — the caller strips it first. "Claims support but ignores it"
// becomes unrepresentable rather than merely discouraged, and there is still exactly one
// compiled artifact, so Get is untouched.
//
// The caller's sequence:
//
//	if expander, ok := component.(ExpandingComponent); ok {
//	    compiled, err = expander.CompileWithExpansion(ctx, query)   // whole query, one artifact
//	} else {
//	    compiled, err = component.CompileQuery(ctx, query.WithoutExpand())
//	    // ... and expand the results with the hop executor
//	}
type ExpandingComponent interface {
	// CompileWithExpansion compiles a query INCLUDING its Expand projection into the single
	// opaque condition Get already accepts. A provider implementing it must honour every
	// expansion option the query carries or return an error naming the one it cannot —
	// silently dropping an option is the failure this interface exists to make impossible.
	CompileWithExpansion(ctx core.ServerContext, query *Query) (interface{}, error)
}

// Query is the root of a data-layer query. It is the single representation every front-end
// lowers into — OData filter text, declarative dataset filters, and the map shorthand — and
// the only one a data provider compiles.
type Query struct {
	Version QueryVersion
	// Filter is the predicate tree. A nil filter means the query is unconstrained; note that
	// this is distinct from a nil condition passed to a data service, which returns nothing.
	Filter Predicate
	// Expand is the $expand projection: which reference fields to resolve and attach to each
	// returned record, and what to expand beneath them. It is a sibling of Filter because
	// expansion changes what a record carries, never which records match.
	//
	// HOW IT REACHES A PROVIDER, because the contract does not make this visible and the
	// obvious reading of it is wrong.
	//
	// CompileQuery does NOT collapse the query — every shipped provider returns the *Query
	// itself as its compiled artifact, after validating it. The collapse happens one step
	// later, at BindQuery, and it is per-provider: sqldatabase binds to a condition that still
	// carries the whole *Query, so expansion is reachable when its statement is built, while
	// mongodatabase binds to a bson filter, which has no way to express a projection and so
	// drops it. Get then receives whatever that bind produced.
	//
	// So the transport already exists for some providers and not others, by accident rather
	// than by contract. That is what ExpandingComponent is for: it makes expansion arrive by
	// one declared route instead of by whatever each provider's condition happens to retain.
	// A provider implementing it receives the whole query through CompileWithExpansion; a
	// provider that does not never sees Expand at all, because the caller strips it with
	// WithoutExpand and expands the results with the hop executor instead.
	//
	// The strip is the load-bearing half. Leaving Expand on a query bound for a provider that
	// ignores it returns unexpanded records with a nil error — silently wrong, in the one
	// direction nothing downstream can detect.
	//
	// An empty Expand is every query written before this field existed, and returns records
	// with their references unresolved exactly as it always did.
	Expand []Expansion
	// Navigate changes what the query RETURNS: instead of the entity the data service holds,
	// it returns the entities reached by following these navigation segments from each record
	// that matched Filter. This is OData's URL path addressing,
	// /Groups('Github')/Admin/Friends/Pets.
	//
	// It is the third and last of the relationship constructs, and the three are distinct in
	// a way worth stating together because they are easy to confuse: Expand returns parents
	// with children attached, a Traversal returns parents filtered by their children, and
	// Navigate returns the children.
	//
	// THIS FIELD NEVER REACHES A PROVIDER, and that is stated here because a provider author
	// will otherwise search the condition for it and find nothing. It cannot reach one: it
	// changes the result ENTITY, and a DataComponent is bound to a single entity — every
	// execution method returns []core.Storable of the component's own object. Widening the
	// interface to carry it is forbidden for the usual reason, that a new interface method
	// breaks implementors at load rather than at compile.
	//
	// It is a directive to the SERVER. The hop executor resolves the target entity's component
	// through DataManager and executes there, which is the one use DataManager is reserved for
	// on the query path — a provider that executes traversal natively resolves nothing.
	//
	// An empty Navigate returns the query's own entity.
	Navigate []string
}

// NewQuery builds an empty query at the current version.
func NewQuery() *Query {
	return &Query{Version: QueryV1}
}

// NewEqualityQuery lowers the map shorthand into a query: every entry becomes an equality
// comparison and the entries are combined with and. This is the form the large majority of
// callers use, and it exists here so that each consuming module does not rewrite it.
func NewEqualityQuery(args utils.StringMap) *Query {
	query := NewQuery()
	if len(args) == 0 {
		return query
	}
	operands := make([]Predicate, 0, len(args))
	for field, value := range args {
		operands = append(operands, &Comparison{
			Field:    field,
			Operator: OpEqual,
			Value:    LiteralOperand(value),
		})
	}
	if len(operands) == 1 {
		query.Filter = operands[0]
		return query
	}
	query.Filter = &Logical{Operator: LogicalAnd, Operands: operands}
	return query
}

// The chaining API below exists because the AST's expressive constructs were reachable only by
// hand-written struct literals — verbose, undiscoverable without reading this file, and offering
// no help at the call site. A capability reachable only that way goes unused.
//
// These are methods on a struct, so they break no implementor and cost no flag day to add. What
// they deliberately do NOT do is make field names type-safe: Where("Nmae", ...) still compiles.
// Only per-entity generated code can catch that, and this API is designed to be what such code
// builds on rather than a substitute for it.
//
// Two operations people expect in a chain and will not find here, with the reason:
//
//   - Count. Counting is a terminal operation on a DATA SERVICE, not a property of a query —
//     DataComponent.Count(ctx, condition) returns one, and the ordinary Get already returns
//     totalrecs alongside the page. A .Count() here would either duplicate that or silently
//     mean something different from OData's $count.
//   - Grouping and aggregation. Not in this AST at all, deliberately: DataComponent.CountGroups
//     is the only grouping the platform has, it counts and nothing else, and general aggregation
//     is out of scope for this work. A .GroupBy() that lowered to nothing would be worse than
//     its absence.

// Where adds a predicate to the query's filter, combining with anything already there under and.
// Calling it repeatedly is therefore additive rather than replacing, which is what makes a chain
// built across several statements behave the way it reads.
func (q *Query) Where(predicate Predicate) *Query {
	if q == nil || predicate == nil {
		return q
	}
	switch {
	case q.Filter == nil:
		q.Filter = predicate
	default:
		// flatten into an existing top-level conjunction rather than nesting one inside
		// another, so that three Where calls produce one and-node and not three
		if logical, ok := q.Filter.(*Logical); ok && logical.Operator == LogicalAnd && !logical.Optional {
			logical.Operands = append(logical.Operands, predicate)
			return q
		}
		q.Filter = &Logical{Operator: LogicalAnd, Operands: []Predicate{q.Filter, predicate}}
	}
	return q
}

// Through filters the records this query returns by a predicate on a RELATED entity, which is the
// traversal without its struct literal. The predicate's field names resolve on the entity the
// path reaches, never on the one the query started from.
//
// It returns parents filtered by their children. To return the parents WITH their children, use
// Expanding; to return the children themselves, use NavigatingTo.
func (q *Query) Through(path []string, quantifier Quantifier, predicate Predicate) *Query {
	return q.Where(&Traversal{Path: path, Quantifier: quantifier, Predicate: predicate})
}

// Expanding appends expansions to the query's projection, so each returned record carries the
// related records named. Call it more than once and the expansions accumulate.
func (q *Query) Expanding(expansions ...Expansion) *Query {
	if q == nil || len(expansions) == 0 {
		return q
	}
	q.Expand = append(q.Expand, expansions...)
	return q
}

// NavigatingTo makes the query return records of a RELATED entity instead of its own, following
// the named navigation segments. It replaces rather than appends: a query has one result type.
func (q *Query) NavigatingTo(segments ...string) *Query {
	if q == nil {
		return q
	}
	q.Navigate = segments
	return q
}

// NewExpansion builds an expansion of one reference field, with any nested expansions beneath it.
// The option setters below chain onto it, and each returns a copy, so an expansion built once can
// be varied without the variants sharing state.
func NewExpansion(field string, nested ...Expansion) Expansion {
	return Expansion{Field: field, Expand: nested}
}

// Where constrains which RELATED records are attached. It filters the children and never the
// parents — a parent whose children all fail the predicate is still returned, carrying an empty
// expansion. Query.Through is the one that drops such a parent.
func (e Expansion) Where(predicate Predicate) Expansion {
	e.Filter = predicate
	return e
}

// Selecting limits which fields of the related records are read. No call means all of them.
func (e Expansion) Selecting(fields ...string) Expansion {
	e.Select = fields
	return e
}

// OrderedBy orders the expanded set, taking FIELD/DIRECTION PAIRS — OrderedBy("Name", "asc") —
// which is the shape the Get family's orderBy already uses and which providers walk two at a
// time. It is not a list of "Name asc" clauses.
func (e Expansion) OrderedBy(fieldsAndDirections ...string) Expansion {
	e.OrderBy = fieldsAndDirections
	return e
}

// Limit bounds the expanded set to top records starting at skip. It is the cheapest defence
// against expansion fan-out, and is worth pairing with OrderedBy: a bound without an order
// returns an arbitrary subset rather than a chosen one.
func (e Expansion) Limit(top, skip int) Expansion {
	e.Top, e.Skip = top, skip
	return e
}

// Deep expands this field repeatedly, for a self-referencing hierarchy — OData's $levels. It
// repeats THIS field, and has nothing to do with expansions nested beneath it, which is the
// distinction most easily got wrong. LevelsMax is the unbounded form.
func (e Expansion) Deep(levels int) Expansion {
	e.Levels = levels
	return e
}

// Counted asks for the expanded set's full size alongside the attached page, which is what
// distinguishes a truncated expansion from a complete one when Limit is in play.
func (e Expansion) Counted() Expansion {
	e.Count = true
	return e
}

// Inner drops parent records whose expansion resolves to nothing, turning the outer join an
// expansion performs by default into an inner one.
func (e Expansion) Inner() Expansion {
	e.Required = true
	return e
}

// Parameters lists every parameter name the query depends on, so a caller can tell what a
// compiled query still needs before binding it.
func (q *Query) Parameters() []string {
	if q == nil {
		return nil
	}
	var names []string
	if q.Filter != nil {
		names = append(names, q.Filter.Parameters()...)
	}
	// an expansion carries its own filter, which may be parameterised — omitting these would
	// hand a provider an expansion filter whose parameter was never bound, and the caller
	// would have had no way to learn it was needed
	for _, expansion := range q.Expand {
		names = append(names, expansion.Parameters()...)
	}
	return names
}

// Resolve returns a copy of the query with optional predicates whose parameters are not
// supplied removed from the tree. This is what allows a caller to supply a subset of the
// parameters a query declares: an unsupplied optional filter drops out, rather than being
// compared against an empty value and matching nothing.
//
// A predicate that is not marked optional is left in place even when its parameter is
// missing, so that a genuinely required parameter fails at execution rather than silently
// widening the result set.
func (q *Query) Resolve(params utils.StringsMap) *Query {
	if q == nil {
		return nil
	}
	// projection and result type are structural: they are not parameterised and so are never
	// pruned. clone carries them, so a resolved query cannot silently lose its expansion and
	// return the right records with the wrong shape.
	resolved := q.clone()
	resolved.Filter = nil
	if q.Filter != nil {
		resolved.Filter = resolvePredicate(q.Filter, params)
	}
	// expansion filters bind on the same terms as the query's own filter. clone copied the
	// slice header, so this replaces it rather than writing through it — pruning in place
	// would mutate the caller's query, and a compiled query is reused across requests.
	resolved.Expand = resolveExpansions(q.Expand, params)
	return resolved
}

// clone copies every field of the query.
//
// It exists so that the field list is enumerated in ONE place. Both callers that need a copy —
// Resolve, which prunes the filter, and WithoutExpand, which drops the projection — would
// otherwise hand-enumerate it, and a field added to Query and forgotten in either is dropped
// with no compile error: the query still executes, and returns the wrong shape rather than an
// error. That defect has already been written once in this file's history.
func (q *Query) clone() *Query {
	return &Query{Version: q.Version, Filter: q.Filter, Expand: q.Expand, Navigate: q.Navigate}
}

// WithoutExpand returns a copy of the query with the projection removed, leaving everything else
// intact. A nil receiver returns nil, and a query with no projection is copied unchanged rather
// than special-cased, so the caller never has to branch.
//
// This is the strip half of expansion negotiation. A caller asks SupportsQuery(CapabilityExpand);
// when the provider declines, it compiles THIS query instead of the original and expands the
// results itself with the hop executor. Compiling the original would hand a provider a projection
// it ignores, and unexpanded records would come back with a nil error — wrong in the one
// direction nothing downstream can detect.
//
// It is a method here rather than a copy at each call site for the reason given on clone: the
// copy must enumerate every other field of Query, and this package is where that list lives.
func (q *Query) WithoutExpand() *Query {
	if q == nil {
		return nil
	}
	stripped := q.clone()
	stripped.Expand = nil
	return stripped
}

// resolvePredicate prunes a predicate tree against the supplied parameters, returning nil
// when the node itself is elided.
func resolvePredicate(predicate Predicate, params utils.StringsMap) Predicate {
	if predicate == nil {
		return nil
	}
	// a logical node is pruned child-first, and disappears when nothing survives beneath it
	if logical, ok := predicate.(*Logical); ok {
		survivors := make([]Predicate, 0, len(logical.Operands))
		for _, operand := range logical.Operands {
			if resolvedOperand := resolvePredicate(operand, params); resolvedOperand != nil {
				survivors = append(survivors, resolvedOperand)
			}
		}
		switch len(survivors) {
		case 0:
			return nil
		case 1:
			return survivors[0]
		}
		return &Logical{Optionality: logical.Optionality, Operator: logical.Operator, Operands: survivors}
	}
	if not, ok := predicate.(*Not); ok {
		resolvedOperand := resolvePredicate(not.Operand, params)
		if resolvedOperand == nil {
			return nil
		}
		return &Not{Optionality: not.Optionality, Operand: resolvedOperand}
	}
	// a traversal is pruned through its sub-predicate, and disappears with it: a traversal
	// whose predicate was elided asserts only that the path reaches something, which is not
	// what the caller wrote and would silently narrow the result set
	if traversal, ok := predicate.(*Traversal); ok {
		resolvedPredicate := resolvePredicate(traversal.Predicate, params)
		if resolvedPredicate == nil {
			return nil
		}
		return &Traversal{
			Optionality:   traversal.Optionality,
			Path:          traversal.Path,
			Quantifier:    traversal.Quantifier,
			Predicate:     resolvedPredicate,
			MinDepth:      traversal.MinDepth,
			MaxDepth:      traversal.MaxDepth,
			MatchOptional: traversal.MatchOptional,
		}
	}
	if predicate.IsOptional() && !parametersBound(predicate.Parameters(), params) {
		return nil
	}
	return predicate
}

// parametersBound reports whether every named parameter was supplied.
func parametersBound(names []string, params utils.StringsMap) bool {
	for _, name := range names {
		if _, ok := params[name]; !ok {
			return false
		}
	}
	return true
}
