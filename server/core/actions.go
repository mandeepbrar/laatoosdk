package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/utils"
)

type ActionType string

type ActionExecutor func(ctx RequestContext, action *Action, params utils.StringMap) (interface{}, error)

type Action struct {
	Type   ActionType
	Config config.GenericConfig
}
