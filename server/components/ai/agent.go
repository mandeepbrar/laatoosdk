package ai

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type Agent interface {
	core.Service
	Invoke(ctx core.RequestContext) error
	GetAgentType() AgentType
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
