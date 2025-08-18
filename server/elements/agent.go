package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type InformationBucket interface {
	core.ConfigurableObjectInfo
}

type Agent interface {
	Service
	GetModel() string
	GetVersion() string
	GetInstructions() string
	GetDescription() string
	Information() []InformationBucket
	Tools() []Service
}

type Conversation interface {
	GetId() string
	AssignAgent(Agent)
	GetPresentAgent() Agent
	GetHistory() utils.StringsMap
	AddHistory(ctx core.RequestContext, actor string, input utils.StringMap)
}
