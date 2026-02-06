package ai

import (
	"laatoo.io/sdk/server/core"
)

type Agent interface {
	core.Service
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
	AgentTypeWorkflow       AgentType = "workflow"
	AgentTypeResearch       AgentType = "research"
	AgentTypeGoal         AgentType = "goal"
	AgentTypeOthers       AgentType = "others"
)

