package ai

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ConversationMessage is one turn of an LLM conversation. It is what CompletionRequest.Messages
// carries into every provider (openai, anthropic, gemini, ollama).
//
// *AIMessage in this package is the only implementation; use it, and prefer the NewAgentMessage
// / NewUserAIMessage / NewToolMessage / NewAgentToolMessage constructors. Note that none of
// those constructors sets a timestamp, an actor or an ID, and neither does
// AgentManager.WriteMessageToMemory when it builds a turn — so messages read back out of
// session memory routinely have GetTimestamp() == "" and GetActorId() == "".
//
// SEVERAL ACCESSORS BELOW ARE DEAD. GetToolCallId, GetFunctionCalls, GetActorId and
// GetActorName have zero call sites in the platform: the providers read the equivalent values
// out of GetMetadata() instead (see below). Populating the struct fields those accessors
// return has no effect on what is sent to any LLM.
type ConversationMessage interface {
	core.Storable

	// GetActorId returns the id of the participant who produced the message.
	// Nothing in the platform reads it, and nothing populates it.
	GetActorId() string

	// GetActorName returns the display name of the participant who produced the message.
	// Nothing in the platform reads it, and nothing populates it.
	GetActorName() string

	// GetAttachmentData returns raw attachment bytes, or an empty slice for a text-only
	// message. Every provider checks len()>0 here to decide whether to build a multi-part
	// message, and only an image/* GetAttachmentMimeType is actually transmitted — the bytes
	// are base64'd into a data: URI. A non-image attachment is silently dropped and only the
	// text is sent. The bytes also count toward the providers' token estimates.
	GetAttachmentData() []byte

	// GetAttachmentMimeType returns the attachment's MIME type. Must start with "image/" for
	// the attachment to reach the model; anything else is ignored without warning.
	GetAttachmentMimeType() string

	// GetRole returns who is speaking, and it is the whole of the provider's routing decision.
	//
	// A StakeholderSystem MESSAGE IS SILENTLY DISCARDED on openai, anthropic and gemini — each
	// `continue`s past it (openai.go:162, anthropic.go:129, gemini.go:456). It is NOT folded
	// into the system prompt: that comes from the separate CompletionRequest.SystemPrompt
	// field, so system instructions passed as a message simply never reach the model, with no
	// error. Ollama is the lone exception and forwards it as role "system". Put system
	// instructions in req.SystemPrompt.
	//
	// Of the rest: StakeholderTool (or any message carrying metadata["tool_call_id"]) becomes
	// a tool-result message; StakeholderAgent becomes an assistant message; EVERYTHING ELSE
	// falls through to "user", including StakeholderUnknown and any unrecognised string. A
	// mis-set role is therefore never an error, just a message attributed to the user.
	// See openaillmprovider openai.go:161-224 and the equivalent in each provider.
	GetRole() AgentStakeholder

	// GetMessageContent returns the message text. This is the only field guaranteed to reach
	// the model on every provider and every role.
	GetMessageContent() string

	// GetTimestamp returns an ISO8601/RFC3339 STRING, not a time.Time — and it is empty on
	// every message the platform constructs, because none of the AIMessage constructors nor
	// AgentManager.WriteMessageToMemory sets it. No provider sends it to the model.
	GetTimestamp() string

	// GetToolCallId returns the id of the tool call this message answers.
	//
	// NOT READ BY ANY PROVIDER. The providers take the tool call id from
	// GetMetadata()["tool_call_id"] instead, so a message built with NewToolMessage(content,
	// id) — which sets only the struct field — is sent upstream with an EMPTY ToolCallID and
	// the model cannot match the result to its call. Set metadata["tool_call_id"] as well.
	GetToolCallId() string

	// GetFunctionCalls returns the tool calls an assistant message is making.
	//
	// NOT READ BY ANY PROVIDER. The providers look for GetMetadata()["functionCalls"] holding
	// a JSON-encoded []FunctionCall STRING, and ignore this accessor entirely, so calls set
	// only via NewAgentToolMessage never reach the model.
	GetFunctionCalls() []FunctionCall

	// GetMetadata carries the values the providers actually act on — "tool_call_id" (string)
	// and "functionCalls" (JSON-encoded []FunctionCall as a string) — plus whatever else a
	// caller stores. Merely having a non-empty "tool_call_id" reclassifies the message as a
	// tool result regardless of its role.
	GetMetadata() utils.StringMap
}

// AgentConversation is a storable conversation transcript.
//
// NOTHING IMPLEMENTS OR USES THIS INTERFACE. As of 2026-08-27 there is no type in laatoo,
// laatoosdk or the solutions with this method set, and no call site anywhere names the type.
// Conversation history is instead kept in a session MemoryBank and read back through
// AgentManager.GetMessagesFromMemory. The method contracts below are the interface's intent
// only — they are unverified, and no nil, empty or error convention should be assumed.
type AgentConversation interface {
	core.Storable

	// GetMessages returns the conversation's turns. No implementation exists, so ordering,
	// paging and the empty-conversation result are all undefined.
	GetMessages(ctx core.RequestContext) ([]ConversationMessage, error)

	// AddMessage appends a turn to the conversation. No implementation exists, so whether it
	// persists immediately and how duplicates are treated are both undefined.
	AddMessage(ctx core.RequestContext, input ConversationMessage) error
}
