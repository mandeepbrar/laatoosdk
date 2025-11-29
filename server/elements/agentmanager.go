package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type AgentManager interface {
	core.ServerElement
	GetAgent(ctx core.ServerContext, alias string) (core.Agent, error)
	List(ctx core.ServerContext) utils.StringsMap
	RegisterModel(ctx core.ServerContext, modelName string, model interface{}) error

	GetModel(ctx core.ServerContext, modelsName string) (interface{}, error)
}

type AgentEngine interface {
	GetAgentType(ctx core.ServerContext) string
}
