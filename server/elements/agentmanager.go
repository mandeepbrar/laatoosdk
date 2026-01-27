package elements

import (
	"laatoo.io/sdk/server/components/ai"
	"laatoo.io/sdk/server/core"
)



type AgentManager interface {
	core.ServerElement

	GetAgent(ctx core.ServerContext, alias string) (ai.Agent, error)
//	GetEngine(ctx core.ServerContext, name string) (AgentEngine, error)

	ListAgents(ctx core.ServerContext) map[string]ai.Agent
	RegisterAgentType(ctx core.ServerContext, agenttype string, factory core.ServiceFactory) error


	// Complete sends a prompt and gets a response
	LLMRequest(ctx core.RequestContext, req *ai.CompletionRequest) (*ai.CompletionResponse, error)

	// Stream sends a prompt and streams back responses
	// Returns a channel of StreamEvent
	LLMStreamingRequest(ctx core.RequestContext, req *ai.CompletionRequest) (<-chan ai.StreamEvent, error)


	// MCP Support
	GetMCPServer(ctx core.ServerContext, rootpath string) (ai.Mcp, error)
	RegisterMCPServer(ctx core.ServerContext, rootpath string, mcpsvr ai.Mcp) error

	//LLM Providers Support
	GetLLMProvider(ctx core.ServerContext, name string) (ai.LLMProvider, error)
	RegisterLLMProvider(ctx core.ServerContext, name string, llmprovider ai.LLMProvider) error

	// Skill Support
	ListSkills(ctx core.ServerContext) map[string]ai.Skill
	RegisterSkillType(ctx core.ServerContext, skillType string, factory core.ServiceFactory) error
	GetSkill(ctx core.ServerContext, name string) (ai.Skill, error)
	GetSkillsByTag(ctx core.ServerContext, tag *core.Tag) []ai.Skill

	CreateMemory(ctx core.RequestContext, memorytype ai.MemoryType, id string, config map[string]interface{}) (ai.MemoryBank, error)

	RegisterAgentMemoryManager(ctx core.ServerContext, memorytype ai.MemoryType,mgr ai.AgentMemoryManager) error
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