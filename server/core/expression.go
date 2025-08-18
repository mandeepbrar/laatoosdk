package core

import (
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/utils"
)

type Expression interface {
	IsStaticValue() bool
	IsExpression() bool
	GetExpression() string
	SetManagerData(mgrData interface{})
	GetManagerData() interface{}
	GetDataType() datatypes.DataType
	SetDataType(datatypes.DataType)
	GetValue(ctx RequestContext, vars utils.StringMap) (interface{}, error)
}

type GenericExpression struct {
	Value      interface{}
	Expression string
	mgrData    interface{}
	DType      datatypes.DataType
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
