package elements

import (
	"laatoo.io/sdk/server/components/ai"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ============================================================
// HANDOFF MANAGER - Coordinates agent handoffs across the system
// ============================================================

// HandoffManager coordinates agent handoffs across the system
type HandoffManager interface {
	// Agent registration
	RegisterAgent(ctx core.ServerContext, agent ai.HandoffCapableAgent) error
	UnregisterAgent(ctx core.ServerContext, agentID string) error
	GetAgent(ctx core.ServerContext, agentID string) (ai.HandoffCapableAgent, error)

	// Handoff execution
	ExecuteHandoff(ctx core.RequestContext, req *ai.HandoffRequest) (*ai.HandoffResult, error)

	// Agent discovery
	FindAgentByCapability(ctx core.ServerContext, capabilities []string) (ai.HandoffCapableAgent, error)
	FindAgentByID(ctx core.ServerContext, agentID string) (ai.HandoffCapableAgent, error)

	// Lifecycle
	Start(ctx core.ServerContext) error
	Stop(ctx core.ServerContext) error

	// Monitoring
	GetStatistics(ctx core.ServerContext) map[string]interface{}
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

	// RegisterAgent allows direct registration of pro-code agents (e.g. golangagent)
	RegisterAgent(ctx core.ServerContext, agent ai.Agent) error

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
	
	// HasModel checks if a model exists
	HasModel(ctx core.ServerContext, modelName string) (bool, error)

	// Skill Support
	ListSkills(ctx core.ServerContext) map[string]ai.Skill
	ListSkillSummaries(ctx core.ServerContext) map[string]ai.SkillSummary
	// RegisterSkill allows direct registration of pro-code skills (e.g. golangskill)
	RegisterSkill(ctx core.ServerContext, skill ai.Skill) error
	RegisterSkillType(ctx core.ServerContext, skillType string, factory core.ServiceFactory) error
	GetSkill(ctx core.ServerContext, name string) (ai.Skill, error)
	GetSkillDescriptor(ctx core.ServerContext, name string) (*ai.SkillDescriptor, error)
	GetSkillsByTag(ctx core.ServerContext, tag *core.Tag) []ai.Skill

	CreateMemory(ctx core.RequestContext, memorytype ai.MemoryType, id string, config map[string]interface{}) (ai.MemoryBank, error)
	GetMemory(ctx core.RequestContext, memorytype ai.MemoryType, id string) (ai.MemoryBank, error)

	RegisterAgentMemoryManager(ctx core.ServerContext, memorytype ai.MemoryType, mgr ai.AgentMemoryManager) error

	// WriteMessageToMemory records a conversation turn in the session memory bank.
	// Session ID is read from ctx.GetSession() — no explicit ID param needed.
	// Silently no-ops if session or memory is unavailable.
	WriteMessageToMemory(ctx core.RequestContext, role ai.AgentStakeholder, content string)

	// GetMessagesFromMemory returns all conversation messages for the session in
	// chronological order. Session ID is read from ctx. Returns nil if none exist.
	GetMessagesFromMemory(ctx core.RequestContext, opts utils.StringMap) []ai.ConversationMessage

	// ============================================================
	// HANDOFF MANAGEMENT
	// ============================================================

	// GetHandoffManager returns the centralized handoff manager
	GetHandoffManager() HandoffManager

	// FindHandoffAgent discovers agents by capability for handoff targeting
	FindHandoffAgent(ctx core.ServerContext, capabilities []string) (ai.HandoffCapableAgent, error)

	// ============================================================
	// HITL MANAGEMENT
	// ============================================================

	// GetHITLManager returns the server-level HITL coordinator.
	// Returns nil if no HITLManager has been registered.
	GetHITLManager() ai.HITLManager

	// RegisterHITLManager sets the server-level HITL coordinator.
	// Called once at startup by the server or a plugin that owns the implementation.
	RegisterHITLManager(ctx core.ServerContext, mgr ai.HITLManager) error

	// ============================================================
	// COMPLETION REQUEST FACTORIES
	// Returns pre-configured CompletionRequest instances for common use-cases.
	// ============================================================

	// DefaultCompletionRequest returns a CompletionRequest with sensible production defaults.
	DefaultCompletionRequest() *ai.CompletionRequest

	// DefaultCompletionRequestCostSensitive returns defaults optimized for minimum cost.
	DefaultCompletionRequestCostSensitive() *ai.CompletionRequest

	// DefaultCompletionRequestHighQuality returns defaults optimized for quality & reasoning.
	DefaultCompletionRequestHighQuality() *ai.CompletionRequest

	// DefaultCompletionRequestFast returns defaults optimized for lowest time-to-first-token.
	DefaultCompletionRequestFast() *ai.CompletionRequest

	// DefaultCompletionRequestBatching returns defaults optimized for batch processing.
	DefaultCompletionRequestBatching() *ai.CompletionRequest

	// DefaultCompletionRequestResearch returns defaults for research & analysis workloads.
	DefaultCompletionRequestResearch() *ai.CompletionRequest

	// DefaultCompletionRequestModerationSafe returns defaults for safe/moderated content.
	DefaultCompletionRequestModerationSafe() *ai.CompletionRequest

	// ============================================================
	// COMPLETION REQUEST BUILDERS
	// Immutable-style helpers that clone a base request with targeted overrides.
	// ============================================================

	// CloneWithOverrides creates a deep copy of a CompletionRequest and applies the
	// provided override function to it. Pointer-typed sub-structs are individually
	// cloned so the original is never mutated.
	CloneWithOverrides(base *ai.CompletionRequest, overrides func(*ai.CompletionRequest)) *ai.CompletionRequest

	// SetRequestMetadata sets common tracking fields (RequestID, UserID, AgentName)
	// on the request in-place and returns it for chaining.
	SetRequestMetadata(req *ai.CompletionRequest, requestID, userID, agentName string) *ai.CompletionRequest

	// WithModel returns a clone of req with the given model selected.
	WithModel(req *ai.CompletionRequest, model string) *ai.CompletionRequest

	// WithTemperature returns a clone of req with the given temperature.
	WithTemperature(req *ai.CompletionRequest, temp float32) *ai.CompletionRequest

	// WithMaxTokens returns a clone of req with the given max-token limit.
	WithMaxTokens(req *ai.CompletionRequest, tokens int) *ai.CompletionRequest

	// WithBudget returns a clone of req with the given USD cost cap.
	// The alert threshold is automatically set to 70 % of maxCostUSD.
	WithBudget(req *ai.CompletionRequest, maxCostUSD float64) *ai.CompletionRequest

	// WithStreaming returns a clone of req with streaming enabled.
	WithStreaming(req *ai.CompletionRequest) *ai.CompletionRequest
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
