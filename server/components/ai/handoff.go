package ai

import (
	"time"
	"laatoo.io/sdk/server/core"
)

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
	Strategy            string
	ExecutionMode       string
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
