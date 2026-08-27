package core

import (
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/utils"
)

// Expression is a value that may be either a literal or a govaluate expression evaluated per
// request -- the mechanism behind a workflow branch condition, a switch case and an agent step
// guard.
//
// An expression must be registered with ctx.RegisterExpression before it can be evaluated:
// registration is what compiles the source and stores the compiled form via SetManagerData
// (laatooserver core/expressionsmanager.go:41-63). An UNREGISTERED expression is not an error at
// evaluation time -- GetValue returns (nil, nil), so an `if err != nil` guard passes and the
// caller proceeds with a nil value. The same (nil, nil) is returned for a nil expression.
type Expression interface {
	// IsStaticValue reports whether no expression source was supplied. Beware: a static
	// expression is never compiled, so GetValue on one returns (nil, nil) -- see GetValue.
	IsStaticValue() bool
	// IsExpression reports whether an expression source was supplied. It is the exact negation of
	// IsStaticValue.
	IsExpression() bool
	// GetExpression returns the raw expression source, in govaluate syntax.
	GetExpression() string
	// SetManagerData stores the compiled expression. The expressions manager calls it during
	// registration; callers should not.
	SetManagerData(mgrData interface{})
	// GetManagerData returns the compiled expression, or nil when the expression was never
	// registered. A nil result is what makes GetValue return (nil, nil).
	GetManagerData() interface{}
	// GetDataType returns the type the expression was registered as.
	GetDataType() datatypes.DataType
	// SetDataType records the expected result type. The expressions manager calls it during
	// registration, from the dtype passed to RegisterExpression. Nothing coerces the evaluated
	// result to it -- govaluate's own result type is returned unchanged.
	SetDataType(datatypes.DataType)
	// GetValue evaluates the expression against vars and returns govaluate's result.
	//
	// It returns (nil, nil) -- no value AND no error -- in three cases: the expression is nil, it
	// was never registered, or it is a static value (a static value is never compiled, so it has
	// no manager data; see the GenericExpression.Value comment)
	// (laatooserver core/expressionsmanager.go:64-74). Only a request that reaches the manager at
	// all returns a real error, and only when govaluate itself fails.
	GetValue(ctx RequestContext, vars utils.StringMap) (interface{}, error)
}

// GenericExpression is the stock Expression implementation. Every caller in laatooserver and
// laatoomodules constructs it as &core.GenericExpression{Expression: src} and then registers it
// with ctx.RegisterExpression.
type GenericExpression struct {
	// Value is DEAD: nothing in the SDK, laatooserver or laatoomodules ever reads it. Setting it
	// and calling GetValue returns (nil, nil), because evaluation goes through the expressions
	// manager and the manager only ever returns the compiled expression's result. Do not use this
	// field to carry a literal.
	Value interface{}
	// Expression is the govaluate source. An empty string makes the expression static -- which,
	// given the Value field above, means it can never produce a value.
	Expression string
	mgrData    interface{}
	// DType is the declared result type, written by SetDataType during registration.
	DType datatypes.DataType
}

func (expr *GenericExpression) IsStaticValue() bool {
	return expr.Expression == ""
}

func (expr *GenericExpression) IsExpression() bool {
	return expr.Expression != ""
}

func (expr *GenericExpression) GetExpression() string {
	return expr.Expression
}

func (expr *GenericExpression) GetDataType() datatypes.DataType {
	return expr.DType
}
func (expr *GenericExpression) SetDataType(dt datatypes.DataType) {
	expr.DType = dt
}

func (expr *GenericExpression) SetManagerData(mgrData interface{}) {
	expr.mgrData = mgrData
}
func (expr *GenericExpression) GetManagerData() interface{} {
	return expr.mgrData
}

func (expr *GenericExpression) GetValue(ctx RequestContext, vars utils.StringMap) (interface{}, error) {
	return ctx.GetExpressionValue(expr, vars)
}
