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
	StakeholderUnknown AgentStakeholder = "Unknown"
)


type AgentType string

const (
	AgentTypeSimple       AgentType = "SimpleAgent"
	AgentTypeDeepResearch AgentType = "DeepResearchAgent"
	AgentTypeFlow         AgentType = "FlowAgent"
	AgentTypeGoal         AgentType = "GoalAgent"
	AgentTypeOthers       AgentType = "Others"
)

