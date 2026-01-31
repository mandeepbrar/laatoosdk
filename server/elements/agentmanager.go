package elements

import (
	"laatoo.io/sdk/server/components/ai"
	"laatoo.io/sdk/server/core"
)

// ============================================================
// HANDOFF MANAGER - Coordinates agent handoffs across the system
// ============================================================

// HandoffManager coordinates agent handoffs across the system
type HandoffManager interface {
	// Agent registration
	RegisterAgent(agent ai.HandoffCapableAgent) error
	UnregisterAgent(agentID string) error
	GetAgent(agentID string) (ai.HandoffCapableAgent, error)
	
	// Handoff execution
	ExecuteHandoff(ctx core.RequestContext, req *ai.HandoffRequest) (*ai.HandoffResult, error)
	
	// Agent discovery
	FindAgentByCapability(ctx core.RequestContext, capabilities []string) (ai.HandoffCapableAgent, error)
	FindAgentByID(agentID string) (ai.HandoffCapableAgent, error)
	
	// Lifecycle
	Start(ctx core.ServerContext) error
	Stop(ctx core.ServerContext) error
	
	// Monitoring
	GetStatistics() map[string]interface{}
}

// ============================================================
// AGENT MANAGER - Manages AI agents and their capabilities
// ============================================================

type AgentManager interface {
	core.ServerElement

	GetAgent(ctx core.ServerContext, alias string) (ai.Agent, error)
//	GetEngine(ctx core.ServerContext, name string) (AgentEngine, error)

	ListAgents(ctx core.ServerContext) map[string]ai.Agent
	RegisterAgentType(ctx core.ServerContext, agenttype ai.AgentType, factory core.ServiceFactory) error


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
	
	// ============================================================
	// HANDOFF MANAGEMENT
	// ============================================================
	
	// GetHandoffManager returns the centralized handoff manager
	GetHandoffManager() HandoffManager
	
	// FindHandoffAgent discovers agents by capability for handoff targeting
	FindHandoffAgent(ctx core.ServerContext, capabilities []string) (ai.HandoffCapableAgent, error)
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