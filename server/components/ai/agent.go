package ai

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type Agent interface {
	core.Service
	Invoke(ctx core.RequestContext) error
	GetAgentType() AgentType
	GetAgentPreferences() *AgentPreferences

	// WriteMessageToMemory records a conversation turn in the session memory bank.
	// Session ID is read from ctx. Silently no-ops if session is unavailable.
	WriteMessageToMemory(ctx core.RequestContext, role AgentStakeholder, content string)

	// GetMessagesFromMemory returns all conversation messages for the session in
	// chronological order. Session ID is read from ctx. Returns nil if no messages exist.
	GetMessagesFromMemory(ctx core.RequestContext) []ConversationMessage
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
