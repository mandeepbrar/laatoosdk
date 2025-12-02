package core

import (
	"laatoo.io/sdk/utils"
)

type InformationBucket interface {
	ConfigurableObjectInfo
}

type Agent interface {
	Service
	GetAgentType() string
	GetVersion() string
	GetDescription() string
	GetAgentProperties() utils.StringMap
}

type AgentConversation interface {
	GetId() string
	AssignAgent(Agent)
	GetPresentAgent() Agent
	GetHistory() utils.StringsMap
	AddHistory(ctx RequestContext, actor string, input utils.StringMap)
}
