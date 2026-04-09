package ai

import "laatoo.io/sdk/server/core"

// Skill represents a modular expertise package that agents can discover and use
type Skill interface {
	core.UserInvokableService
	GetSkillDescriptor(ctx core.ServerContext) (*SkillDescriptor, error)
	GetSkillType() string
	GetExamples() []Example

	// WriteMessageToMemory records a message in the session memory bank.
	// Session ID is read from ctx. Silently no-ops if session is unavailable.
	WriteMessageToMemory(ctx core.RequestContext, role AgentStakeholder, content string)

	// GetMessagesFromMemory returns the full conversation history for the session.
	// Session ID is read from ctx. Returns nil if no messages exist.
	GetMessagesFromMemory(ctx core.RequestContext) []ConversationMessage
}
