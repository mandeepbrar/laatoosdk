package config

import (
	"laatoo/sdk/ctx"
	"regexp"
	"strings"
)

var varReplacer = regexp.MustCompile(`\[(.*?)\]`)

func fillVariables(ctx ctx.Context, val interface{}) interface{} {
	expr, ok := val.(string)
	if !ok {
		return val
	}
	if strings.HasPrefix(expr, "{") && strings.HasSuffix(expr, "}") {
		//expression needs to be evaluated
		expr = strings.TrimSuffix(strings.TrimPrefix(expr, "{"), "}")
	} else {
		return expr
	}

	if strings.HasPrefix(expr, "[") && strings.HasSuffix(expr, "]") {
		testVar := strings.TrimSuffix(strings.TrimPrefix(expr, "["), "]")
		val, ok := ctx.Get(testVar)
		if ok {
			return val
		}
	}

	return varReplacer.ReplaceAllStringFunc(expr, func(exp string) string {
		removebrackets := exp[1 : len(exp)-1]
		r := strings.NewReplacer(".", "", "/", "", "-", "", ":", "")
		varname := r.Replace(removebrackets)
		//varname := strings.Replace(removebrackets, ".", "", -1)
		val, ok := ctx.Get(varname)
		if ok {
			valStr, ok := val.(string)
			if ok {
				return strings.Replace(removebrackets, varname, valStr, -1)
			} else {
				return exp
			}
		} else {
			return exp
		}
	})
}
