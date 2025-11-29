package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type AgentManager interface {
	core.ServerElement
	GetAgent(ctx core.ServerContext, alias string) (core.Agent, error)
	List(ctx core.ServerContext) utils.StringsMap
	GetModels(ctx core.ServerContext, models []string) (utils.StringMap, error)
	GetTools(ctx core.ServerContext, tools []string) (utils.StringMap, error)
}

type AgentEngine interface {
	GetAgentType(ctx core.ServerContext) string
}
