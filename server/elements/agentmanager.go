package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type AgentManager interface {
	core.ServerElement
	GetAgent(ctx core.ServerContext, alias string) (core.Agent, error)
	GetEngine(ctx core.ServerContext, name string) (AgentEngine, error)
	List(ctx core.ServerContext) utils.StringsMap
	RegisterAgentType(ctx core.ServerContext, agenttype string, agentengine string, factory core.ServiceFactory) error

	// MCP Support
	GetMCPServer(ctx core.ServerContext, rootpath string) (components.Mcp, error)
	RegisterMCPServer(ctx core.ServerContext, rootpath string, engine components.Mcp) error
}

type AgentEngine interface {
	core.ServerElement
	GetTools(ctx core.ServerContext, toolNames []string) (utils.StringMap, error)
	RegisterTool(ctx core.ServerContext, toolName string, svc Service) error
	GetModel(ctx core.ServerContext, modelName string) (interface{}, error)
	RegisterModel(ctx core.ServerContext, modelName string, model interface{}) error
}
