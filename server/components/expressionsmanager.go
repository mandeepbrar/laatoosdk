package components

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type ExpressionsManager interface {
	RegisterExpression(ctx core.ServerContext, expression core.Expression, dtype datatypes.DataType) error
	GetExpressionValue(ctx ctx.Context, expression core.Expression, vars utils.StringMap) (interface{}, error)
}
