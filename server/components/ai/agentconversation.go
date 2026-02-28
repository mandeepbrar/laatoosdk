package ai

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)
type ConversationMessage interface {
	core.Storable
	GetActorId() string
	GetActorName() string
	GetAttachmentData() []byte
	GetAttachmentMimeType() string
	GetRole() AgentStakeholder
	GetMessageContent() string
	GetTimestamp() string
	GetToolCallId() string
	GetFunctionCalls() []FunctionCall
	GetMetadata() utils.StringMap
}


type AgentConversation interface {
	core.Storable	
	GetMessages(ctx core.RequestContext) ([]ConversationMessage, error)
	AddMessage(ctx core.RequestContext, input ConversationMessage) error
}