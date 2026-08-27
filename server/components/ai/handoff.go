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

// HandoffCapableAgent defines agents that can participate in handoffs.
//
// Membership is decided by a bare type assertion, never by a declaration: the agent manager
// does `agt.(ai.HandoffCapableAgent)` after building each agent and, on success, registers it
// with the HandoffManager (laatooserver/src/core/agentmanager.go:191 and :226). There is no
// else-branch and no log line, so an agent that misses this interface by one method signature
// is silently excluded from every handoff — it starts, serves requests, and is simply never
// routed to. This has already happened once in the platform: WorkflowAgentService declares
// `CanHandoff() bool` with no ctx parameter (workflowagents agent.go:300), so it does not
// satisfy this interface and workflow agents never reach the handoff registry despite
// implementing the other three methods. Add `var _ ai.HandoffCapableAgent = (*YourAgent)(nil)`
// to your plugin so the compiler catches this instead of production.
type HandoffCapableAgent interface {
	Agent

	// GetCapabilities returns the capability strings this agent advertises. The HandoffManager
	// reads it once at registration and once at unregistration to build and unwind its
	// capability index, so a value that changes afterwards leaves the index stale.
	// Capabilities are matched literally — FindAgentByCapability does no normalization, so
	// "billing" and "Billing" are different capabilities.
	GetCapabilities(ctx core.ServerContext) []string

	// CanHandoff reports whether the agent will initiate a handoff.
	//
	// Nothing in the platform calls this — the manager registers on the type assertion alone,
	// so returning false does not opt an agent out of handoff targeting. The one shipped
	// implementation (GoalAgentService) returns
	// `len(HandoffRules) > 0 || s.handoffManager != nil`, and the manager always injects a
	// handoff manager, so in practice it is a constant true.
	CanHandoff(ctx core.ServerContext) bool

	// RequestHandoff initiates a handoff from this agent to another and returns the outcome.
	// Both shipped implementations simply forward to HandoffManager.ExecuteHandoff, which
	// resolves req.TargetAgentID by capability routing when it is empty.
	RequestHandoff(ctx core.RequestContext, req *HandoffRequest) (*HandoffResult, error)

	// AcceptHandoff receives a handoff and runs the work described by req.
	//
	// The implementation is expected to return a populated *HandoffResult with Success=false
	// ALONGSIDE a non-nil error on rejection rather than a bare nil — GoalAgentService does
	// this when req.Context["sessionId"] is missing (goalagents/handoffs.go:60). Callers must
	// therefore check the error rather than assuming a non-nil result means success, and must
	// still nil-check the result because that convention is not enforced anywhere.
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
