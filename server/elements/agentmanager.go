package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type ModelDetails struct {
	Name string
	Provider string
	APIKey string
}

type AgentManager interface {
	core.ServerElement
	GetModelDetails(ctx core.ServerContext, name string) (ModelDetails, error)
	GetAgent(ctx core.ServerContext, alias string) (core.Agent, error)
//	GetEngine(ctx core.ServerContext, name string) (AgentEngine, error)
	List(ctx core.ServerContext) utils.StringsMap
	RegisterAgentType(ctx core.ServerContext, agenttype string, factory core.ServiceFactory) error

	// MCP Support
	GetMCPServer(ctx core.ServerContext, rootpath string) (components.Mcp, error)
	RegisterMCPServer(ctx core.ServerContext, rootpath string, mcpsvr components.Mcp) error

	// Skill Support
	ListSkills(ctx core.ServerContext) []core.SkillMetadata
	GetSkill(ctx core.ServerContext, name string) (*core.Skill, error)
	GetSkillsByCategory(ctx core.ServerContext, category string) []*core.Skill
	GetSkillsByTag(ctx core.ServerContext, tag string) []*core.Skill
}
/*
type AgentEngine interface {
	core.ServerElement
	GetTools(ctx core.ServerContext, toolNames []string) (utils.StringMap, error)
	RegisterTool(ctx core.ServerContext, toolName string, svc Service) error
	GetModel(ctx core.ServerContext, modelName string) (interface{}, error)
	RegisterModel(ctx core.ServerContext, modelName string, model interface{}) error
}
*/