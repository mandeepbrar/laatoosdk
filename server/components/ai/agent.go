package ai

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// Agent is a service that the agent manager owns rather than the service manager alone.
//
// An agent is created from a YAML under a module's agents directory: the manager reads
// metadata.agentType, looks up the ServiceFactory registered for it via RegisterAgentType, and
// builds the service through that factory. If the produced service does not satisfy this
// interface, startup fails with Core_Bad_Conf ("Service is not an agent"). An agent whose
// metadata.agentType is absent silently defaults to "goal"; an agentType with no registered
// factory is a NotFound at startup.
// See laatooserver/src/core/agentmanager.go:137.
type Agent interface {
	core.Service

	// Invoke runs the agent for one request. Agents are always created with streaming enabled
	// and an AgentStreamingResponseHandler, so implementations report progress by streaming
	// AgentEventType events (AITHOUGHT / AIFINALRESPONSE / AIERROR) rather than only by
	// setting a response.
	Invoke(ctx core.RequestContext) error

	// GetAgentType returns the agent's kind. Every shipped implementation returns a hardcoded
	// constant for its plugin (goalagents returns AgentTypeGoal, workflowagents returns
	// AgentTypeWorkflow) — it does not echo back the metadata.agentType the agent was
	// configured with, so it cannot be used to recover the configured value.
	GetAgentType() AgentType

	// GetAgentPreferences returns the agent's preferred frontend experience and model, or nil
	// when it has no preference.
	//
	// nil is the normal case, not an error: both shipped agents return nil whenever
	// metadata.aspiredUserExperienceURL is unset, so callers must nil-check before
	// dereferencing. Neither shipped agent ever populates AgentPreferences.Model — it is
	// always the empty string.
	GetAgentPreferences() *AgentPreferences
}

// AspiredUserExperienceProvider is optionally implemented by agents that prefer
// a specific frontend experience URL. Returning nil means no preference.
type AgentPreferences struct {
	ExperienceURL string
	Model         string
}

type AgentStakeholder string

const (
	StakeholderUser    AgentStakeholder = "User"
	StakeholderAgent   AgentStakeholder = "Agent"
	StakeholderSystem  AgentStakeholder = "System"
	StakeholderTool    AgentStakeholder = "tool"
	StakeholderUnknown AgentStakeholder = "Unknown"
)

type AgentType string

const (
	AgentTypeWorkflow AgentType = "workflow"
	AgentTypeResearch AgentType = "research"
	AgentTypeGoal     AgentType = "goal"
	AgentTypeOthers   AgentType = "others"
	AgentTypeGolang   AgentType = "golangagent"
)

type AgentEventType string

const (
	THOUGHT       AgentEventType = "AITHOUGHT"
	FINALRESPONSE AgentEventType = "AIFINALRESPONSE"
	ERROR         AgentEventType = "AIERROR"
)

type AgentData struct {
	Content   string          `json:"content"`
	Metadata  utils.StringMap `json:"metadata,omitempty"`
	TotalCost float64         `json:"total_cost,omitempty"`
	Duration  string          `json:"duration,omitempty"`
}

// ToMap converts AgentData to a StringMap for notifications or generic payloads.
func (d AgentData) ToMap() utils.StringMap {
	m := utils.StringMap{
		"content": d.Content,
	}
	if d.Metadata != nil {
		m["metadata"] = d.Metadata
	}
	if d.TotalCost > 0 {
		m["total_cost"] = d.TotalCost
	}
	if d.Duration != "" {
		m["duration"] = d.Duration
	}
	return m
}
