package components

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ExpressionsManager compiles and evaluates the small govaluate expressions that back computed
// fields and conditional configuration. An expression is compiled once at registration and the
// compiled form is stashed on the core.Expression itself, so evaluation never re-parses.
type ExpressionsManager interface {
	// RegisterExpression stamps dtype onto the expression and, if the expression actually is one
	// (IsExpression()), compiles it and stores the compiled form via SetManagerData.
	//
	// A nil expression — including a typed nil, which is the usual way one arrives — is a SILENT
	// NO-OP returning nil, not an error (laatooserver/src/core/expressionsmanager.go:405-409).
	// A registration that quietly did nothing is indistinguishable from a successful one, and the
	// only symptom appears later as GetExpressionValue returning nil.
	//
	// A syntactically invalid expression IS reported: govaluate's parse error is wrapped and
	// returned.
	RegisterExpression(ctx core.ServerContext, expression core.Expression, dtype datatypes.DataType) error

	// GetExpressionValue evaluates a previously registered expression against vars.
	//
	// RETURNS (nil, nil) IN TWO CASES THAT ARE NOT THE SAME THING: the expression is nil, or it
	// carries no compiled form — which is what a literal (non-expression) value looks like, and
	// also what an expression that was never registered looks like
	// (expressionsmanager.go:421-431). There is no way to tell "evaluated to nothing" from "was
	// never compiled" through this method; check registration succeeded rather than inferring it
	// from a nil result.
	//
	// Evaluation errors from govaluate — an unbound variable, a type mismatch — ARE returned.
	GetExpressionValue(ctx ctx.Context, expression core.Expression, vars utils.StringMap) (interface{}, error)
}
