package ai

import (
	"encoding/json"
	"time"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/log"
	"laatoo.io/sdk/utils"
)

// memorySessionFactory is the minimum AgentManager surface MemoryCapable needs.
// It is satisfied by elements.AgentManager without importing that package
// (which would create an import cycle, since elements imports ai).
type memorySessionFactory interface {
	CreateMemory(ctx core.RequestContext, memorytype MemoryType, id string, config map[string]interface{}) (MemoryBank, error)
	GetMemory(ctx core.RequestContext, memorytype MemoryType, id string) (MemoryBank, error)
}

// MemoryCapable provides default WriteMessageToMemory and GetMessagesFromMemory
// implementations backed by AgentManager's MemoryTypeSession bank.
// Embed this in any Agent or Skill struct to satisfy those interface methods.
// Session ID is always read from ctx — no explicit sessionID param needed.
type MemoryCapable struct{}

// WriteMessageToMemory records a conversation turn in the session memory bank.
// Session ID is read from ctx.GetSession(). Silently no-ops if session is unavailable.
func (m *MemoryCapable) WriteMessageToMemory(ctx core.RequestContext, role AgentStakeholder, content string) {
	if content == "" {
		return
	}
	sessionID := sessionIDFromCtx(ctx)
	if sessionID == "" {
		return
	}
	factory := sessionFactoryFromCtx(ctx)
	if factory == nil {
		return
	}
	bank, err := factory.GetMemory(ctx, MemoryTypeSession, sessionID)
	if err != nil || bank == nil {
		bank, err = factory.CreateMemory(ctx, MemoryTypeSession, sessionID, utils.StringMap{})
		if err != nil || bank == nil {
			log.Warn(ctx, "MemoryCapable.WriteMessageToMemory: could not get or create session memory")
			return
		}
	}
	msg := &AIMessage{Role: role, Content: content}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return
	}
	item := &AIMemoryItem{
		Type:      "message",
		Content:   msgBytes,
		Timestamp: time.Now().Format(time.RFC3339),
		Metadata:  utils.StringMap{"role": string(role)},
	}
	_ = bank.Add(ctx, item)
}

// GetMessagesFromMemory returns all conversation messages for the session in
// chronological order. Session ID is read from ctx. Returns nil if no messages exist.
func (m *MemoryCapable) GetMessagesFromMemory(ctx core.RequestContext) []ConversationMessage {
	sessionID := sessionIDFromCtx(ctx)
	if sessionID == "" {
		return nil
	}
	factory := sessionFactoryFromCtx(ctx)
	if factory == nil {
		return nil
	}
	bank, err := factory.GetMemory(ctx, MemoryTypeSession, sessionID)
	if err != nil || bank == nil {
		return nil
	}
	items, err := bank.Retrieve(ctx, "", nil)
	if err != nil || len(items) == 0 {
		return nil
	}
	var messages []ConversationMessage
	for _, item := range items {
		aiItem, ok := item.(*AIMemoryItem)
		if !ok || aiItem.Type != "message" {
			continue
		}
		var msg AIMessage
		if jsonErr := json.Unmarshal(aiItem.Content, &msg); jsonErr != nil {
			continue
		}
		messages = append(messages, &msg)
	}
	return messages
}

func sessionIDFromCtx(ctx core.RequestContext) string {
	session := ctx.GetSession()
	if session == nil {
		return ""
	}
	return session.GetId()
}

func sessionFactoryFromCtx(ctx core.RequestContext) memorySessionFactory {
	raw := ctx.GetServerElement(core.ServerElementAgentManager)
	if raw == nil {
		return nil
	}
	if f, ok := raw.(memorySessionFactory); ok {
		return f
	}
	return nil
}