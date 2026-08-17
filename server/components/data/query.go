package data

import (
	"laatoo.io/sdk/utils"
)

// QueryVersion identifies the AST revision a Query was built against. It exists so that
// later stages can add shaping, paging and aggregation to the container without changing
// the shape of what already exists.
type QueryVersion int

const (
	// QueryVersionInvalid is the zero value and is never a valid version. Every constant
	// group in this file starts one past zero deliberately: the condition framework this
	// replaces made its first constant the zero value, so a caller that forgot to set the
	// field silently selected a real behaviour.
	QueryVersionInvalid QueryVersion = iota
	// QueryV1 is the predicate-only revision: filters, parameters and extensions.
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
)

// Operand is one side of a comparison.
type Operand struct {
	Kind OperandKind
	// Value holds the operand when Kind is OperandLiteral.
	Value interface{}
	// Name holds the parameter or field name when Kind is OperandParameter or OperandField.
	Name string
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
)

// Query is the root of a data-layer query. It is the single representation every front-end
// lowers into — OData filter text, declarative dataset filters, and the map shorthand — and
// the only one a data provider compiles.
type Query struct {
	Version QueryVersion
	// Filter is the predicate tree. A nil filter means the query is unconstrained; note that
	// this is distinct from a nil condition passed to a data service, which returns nothing.
	Filter Predicate
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

// Parameters lists every parameter name the query depends on, so a caller can tell what a
// compiled query still needs before binding it.
func (q *Query) Parameters() []string {
	if q == nil || q.Filter == nil {
		return nil
	}
	return q.Filter.Parameters()
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
	resolved := &Query{Version: q.Version}
	if q.Filter != nil {
		resolved.Filter = resolvePredicate(q.Filter, params)
	}
	return resolved
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
