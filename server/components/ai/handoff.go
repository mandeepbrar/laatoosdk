package ai

import (
	"time"
	"laatoo.io/sdk/server/core"
)

// HandoffStrategy defines how handoff should be executed
type HandoffStrategy string

const (
	HandoffStrategyRuleBased       HandoffStrategy = "rule_based"
	HandoffStrategyCapabilityBased HandoffStrategy = "capability_based"
	HandoffStrategyLLMDecided      HandoffStrategy = "llm_decided"
	HandoffStrategyHybrid          HandoffStrategy = "hybrid"
)

// HandoffExecutionMode determines direct vs async execution
type HandoffExecutionMode string

const (
	HandoffExecutionModeDirect HandoffExecutionMode = "direct"
	HandoffExecutionModeAsync  HandoffExecutionMode = "async"
	HandoffExecutionModeAuto   HandoffExecutionMode = "auto"
)

// HandoffConfig defines handoff settings for a specific task or agent
type HandoffConfig struct {
	Enabled       bool               `json:"enabled" yaml:"enabled"`
	Strategy      string             `json:"strategy" yaml:"strategy"`
	Conditions    []HandoffCondition `json:"conditions" yaml:"conditions"`
	TargetAgents  []string           `json:"target_agents" yaml:"target_agents"`
	ContextFields []string           `json:"context_fields" yaml:"context_fields"`
	ReturnControl bool               `json:"return_control" yaml:"return_control"`
}

// HandoffCondition defines a trigger for handoff
type HandoffCondition struct {
	Type     string                 `json:"type" yaml:"type"`
	Field    string                 `json:"field,omitempty" yaml:"field,omitempty"`
	Operator string                 `json:"operator,omitempty" yaml:"operator,omitempty"`
	Value    interface{}            `json:"value,omitempty" yaml:"value,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// HandoffRule defines an agent-level rule for handoffs
type HandoffRule struct {
	ID               string           `json:"id" yaml:"id"`
	Name             string           `json:"name" yaml:"name"`
	Description      string           `json:"description" yaml:"description"`
	SourceAgent      string           `json:"source_agent" yaml:"source_agent"`
	TargetAgent      string           `json:"target_agent" yaml:"target_agent"`
	TriggerCondition HandoffCondition `json:"trigger_condition" yaml:"trigger_condition"`
	Priority         int              `json:"priority" yaml:"priority"`
}

// HandoffCapableAgent defines agents that can participate in handoffs
// This interface should be implemented by agent services that support handoff functionality
type HandoffCapableAgent interface {
	// Agent identification
	GetAgentID() string
	GetCapabilities() []string
	GetAgentType() AgentType
	
	// Handoff capabilities
	CanHandoff() bool
	
	// Handoff execution
	RequestHandoff(ctx core.RequestContext, req *HandoffRequest) (*HandoffResult, error)
	AcceptHandoff(ctx core.RequestContext, req *HandoffRequest) (*HandoffResult, error)
}

// HandoffRequest represents a request to transition execution to another agent
type HandoffRequest struct {
	RequestID           string
	SourceAgentID       string
	SourceAgentType     AgentType
	TargetAgentID       string
	TargetCapabilities  []string
	Strategy            HandoffStrategy
	ExecutionMode       HandoffExecutionMode
	Reason              string
	Context             map[string]interface{}
	ConversationHistory []interface{}
	ReturnControl       bool
	Priority            int
	Timeout             time.Duration
	CreatedAt           time.Time
}

// HandoffResult represents the outcome of a handoff execution
type HandoffResult struct {
	Success         bool
	RequestID       string
	TargetAgentID   string
	Output          map[string]interface{}
	Error           string
	Duration        time.Duration
	ControlReturned bool
	CostIncurred    float64
}
