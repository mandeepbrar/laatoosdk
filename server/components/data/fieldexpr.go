package data

// A field-level predicate DSL, so a Go service can express a condition without hand-building AST
// structs.
//
// THIS IS NOT A SECOND QUERY BUILDER. QueryBuilder (querybuilder.go) already provides the fluent
// query surface — Where, Through, Expanding, NavigatingTo, Page, OrderBy, All, One, Count — and a
// parallel builder would fork an API that works. What was missing is one rung below it: Where
// takes a fully-constructed Predicate, so every caller wanting `Score > 80` wrote
//
//	&data.Comparison{Field: "Score", Operator: data.OpGreater, Value: data.LiteralOperand(80)}
//
// which is the expressive cliff imperative callers actually hit. This closes that, and composes
// with the existing builder rather than replacing it:
//
//	dm.CreateQuery(ctx, "myplugin.Course").
//	    Where(data.Field("Score").Gt(80)).
//	    Where(data.Field("Status").In("active", "pending")).
//	    All()
//
// Every method returns a Predicate from this package's existing types, so what reaches a provider
// is exactly what a dataset or an OData filter produces. There is no second lowering path to keep
// in step, which is the property that makes this cheap to keep correct.

// FieldExpr names a field partway through building a predicate about it.
type FieldExpr struct {
	field string
}

// Field starts a predicate about one field. The returned value is inert until a comparison is
// called on it.
func Field(name string) FieldExpr {
	return FieldExpr{field: name}
}

// compare builds a comparison against a literal value.
func (f FieldExpr) compare(op CompareOperator, value interface{}) *Comparison {
	return &Comparison{Field: f.field, Operator: op, Value: LiteralOperand(value)}
}

// Eq matches records whose field equals value.
func (f FieldExpr) Eq(value interface{}) *Comparison { return f.compare(OpEqual, value) }

// Ne matches records whose field does not equal value.
func (f FieldExpr) Ne(value interface{}) *Comparison { return f.compare(OpNotEqual, value) }

// Gt matches records whose field is greater than value.
func (f FieldExpr) Gt(value interface{}) *Comparison { return f.compare(OpGreater, value) }

// Gte matches records whose field is greater than or equal to value.
func (f FieldExpr) Gte(value interface{}) *Comparison { return f.compare(OpGreaterEqual, value) }

// Lt matches records whose field is less than value.
func (f FieldExpr) Lt(value interface{}) *Comparison { return f.compare(OpLess, value) }

// Lte matches records whose field is less than or equal to value.
func (f FieldExpr) Lte(value interface{}) *Comparison { return f.compare(OpLessEqual, value) }

// EqParam matches the field against a named request parameter rather than a literal.
//
// A parameter operand is what lets a query be compiled once and bound per request, which is the
// difference between CompileQuery/BindQuery and rebuilding the AST on every call.
func (f FieldExpr) EqParam(name string) *Comparison {
	return &Comparison{Field: f.field, Operator: OpEqual, Value: ParameterOperand(name)}
}

// EqField matches the field against another field on the same record.
func (f FieldExpr) EqField(other string) *Comparison {
	return &Comparison{Field: f.field, Operator: OpEqual, Value: FieldOperand(other)}
}

// In matches records whose field is any of values.
//
// Values are spread rather than taken as a slice: a single slice argument becomes ONE operand
// holding the slice, and every comparison against it then fails. That is a real defect this
// platform has already paid for once, in the KV batch read, and the variadic signature is what
// makes it unrepresentable here.
func (f FieldExpr) In(values ...interface{}) *Membership {
	return &Membership{Field: f.field, Values: literalOperands(values)}
}

// NotIn matches records whose field is none of values.
func (f FieldExpr) NotIn(values ...interface{}) *Membership {
	return &Membership{Field: f.field, Values: literalOperands(values), Negated: true}
}

// IsNull matches records whose field is absent.
func (f FieldExpr) IsNull() *NullTest {
	return &NullTest{Field: f.field}
}

// literalOperands wraps each value as a literal operand.
func literalOperands(values []interface{}) []Operand {
	out := make([]Operand, len(values))
	for i, v := range values {
		out[i] = LiteralOperand(v)
	}
	return out
}

// And combines predicates conjunctively, skipping nils.
//
// Returns the single predicate when given one and nil when given none, so a caller assembling
// conditions in a loop does not have to special-case the empty and one-element cases — which is
// where hand-rolled combinators usually go wrong.
func And(predicates ...Predicate) Predicate { return combine(LogicalAnd, predicates) }

// Or combines predicates disjunctively, skipping nils.
//
// A provider must declare CapabilityDisjunction to compile the result; one that does not will
// refuse the query rather than approximate it.
func Or(predicates ...Predicate) Predicate { return combine(LogicalOr, predicates) }

// combine builds a logical node, collapsing the degenerate cases.
func combine(op LogicalOperator, predicates []Predicate) Predicate {
	kept := make([]Predicate, 0, len(predicates))
	for _, p := range predicates {
		if p != nil {
			kept = append(kept, p)
		}
	}
	switch len(kept) {
	case 0:
		return nil
	case 1:
		// wrapping a single operand in a logical node changes nothing semantically and costs a
		// provider an extra level to walk
		return kept[0]
	default:
		return &Logical{Operator: op, Operands: kept}
	}
}
