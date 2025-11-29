package core

import (
	"laatoo.io/sdk/utils"
)

type InformationBucket interface {
	ConfigurableObjectInfo
}

type Agent interface {
	Service
	GetAgentEngine() string
	GetModel() string
	GetVersion() string
	GetInstructions() string
	GetDescription() string
	Information() []InformationBucket
	Tools() []Service
}

type AgentConversation interface {
	GetId() string
	AssignAgent(Agent)
	GetPresentAgent() Agent
	GetHistory() utils.StringsMap
	AddHistory(ctx RequestContext, actor string, input utils.StringMap)
}
