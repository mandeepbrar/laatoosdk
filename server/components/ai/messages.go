package ai

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/components/data"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// AIMessage is a generic, storable implementation of ai.ConversationMessage.
type AIMessage struct {
	data.StorageInfo
	data.TenantInfo
	Role               AgentStakeholder `json:"Role"`
	Content            string           `json:"Content"`
	ActorID            string           `json:"ActorID,omitempty"`
	ActorName          string           `json:"ActorName,omitempty"`
	AttachmentData     []byte           `json:"AttachmentData,omitempty"`
	AttachmentMimeType string           `json:"AttachmentMimeType,omitempty"`
	TimestampStr       string           `json:"Timestamp,omitempty"`
	ToolCallID         string           `json:"ToolCallID,omitempty"`
	FunctionCalls      []FunctionCall   `json:"FunctionCalls,omitempty"`
	Metadata           utils.StringMap  `json:"Metadata,omitempty"`
}

func (m *AIMessage) GetRole() AgentStakeholder  { return m.Role }
func (m *AIMessage) GetMessageContent() string     { return m.Content }
func (m *AIMessage) GetActorId() string            { return m.ActorID }
func (m *AIMessage) GetActorName() string          { return m.ActorName }
func (m *AIMessage) GetAttachmentData() []byte     { return m.AttachmentData }
func (m *AIMessage) GetAttachmentMimeType() string { return m.AttachmentMimeType }
func (m *AIMessage) GetTimestamp() string          { return m.TimestampStr }
func (m *AIMessage) GetToolCallId() string         { return m.ToolCallID }
func (m *AIMessage) GetFunctionCalls() []FunctionCall { return m.FunctionCalls }
func (m *AIMessage) GetMetadata() utils.StringMap  { return m.Metadata }

func (m *AIMessage) Constructor(c ctx.Context) {
	m.StorageInfo.Constructor(c)
	m.SetSelfReference(m)
}

func (m *AIMessage) Config() *core.StorableConfig {
	return &core.StorableConfig{
		ObjectType:  "ai.AIMessage",
		LabelField:  "Content",
		Multitenant: true,
		Cacheable:   true,
		Collection:  "AIMessages",
	}
}

func (ent *AIMessage) ReadAll(c ctx.Context, cdc datatypes.Codec, rdr datatypes.SerializableReader) error {
	var err error
	if err = rdr.ReadString(c, cdc, "Role", (*string)(&ent.Role)); err != nil {
		return err
	}
	if err = rdr.ReadString(c, cdc, "Content", &ent.Content); err != nil {
		return err
	}
	if err = rdr.ReadString(c, cdc, "ActorID", &ent.ActorID); err != nil {
		return err
	}
	if err = rdr.ReadString(c, cdc, "ActorName", &ent.ActorName); err != nil {
		return err
	}
	if ent.AttachmentData, err = rdr.ReadBytes(c, cdc, "AttachmentData"); err != nil {
		return err
	}
	if err = rdr.ReadString(c, cdc, "AttachmentMimeType", &ent.AttachmentMimeType); err != nil {
		return err
	}
	if err = rdr.ReadString(c, cdc, "Timestamp", &ent.TimestampStr); err != nil {
		return err
	}
	if err = rdr.ReadString(c, cdc, "ToolCallID", &ent.ToolCallID); err != nil {
		return err
	}
	if err = rdr.ReadArray(c, cdc, "FunctionCalls", &ent.FunctionCalls); err != nil {
		// Ignore error if not found
	}
	if err = rdr.ReadMap(c, cdc, "Metadata", &ent.Metadata); err != nil {
		return err
	}

	if err = ent.TenantInfo.ReadAll(c, cdc, rdr); err != nil {
		return err
	}
	return ent.StorageInfo.ReadAll(c, cdc, rdr)
}

func (ent *AIMessage) WriteAll(c ctx.Context, cdc datatypes.Codec, wtr datatypes.SerializableWriter) error {
	var err error
	roleStr := string(ent.Role)
	if err = wtr.WriteString(c, cdc, "Role", &roleStr); err != nil {
		return err
	}
	if err = wtr.WriteString(c, cdc, "Content", &ent.Content); err != nil {
		return err
	}
	if err = wtr.WriteString(c, cdc, "ActorID", &ent.ActorID); err != nil {
		return err
	}
	if err = wtr.WriteString(c, cdc, "ActorName", &ent.ActorName); err != nil {
		return err
	}
	if err = wtr.WriteBytes(c, cdc, "AttachmentData", &ent.AttachmentData); err != nil {
		return err
	}
	if err = wtr.WriteString(c, cdc, "AttachmentMimeType", &ent.AttachmentMimeType); err != nil {
		return err
	}
	if err = wtr.WriteString(c, cdc, "Timestamp", &ent.TimestampStr); err != nil {
		return err
	}
	if err = wtr.WriteString(c, cdc, "ToolCallID", &ent.ToolCallID); err != nil {
		return err
	}
	if len(ent.FunctionCalls) > 0 {
		if err = wtr.WriteArray(c, cdc, "FunctionCalls", &ent.FunctionCalls); err != nil {
			return err
		}
	}
	if err = wtr.WriteMap(c, cdc, "Metadata", &ent.Metadata); err != nil {
		return err
	}

	if err = ent.TenantInfo.WriteAll(c, cdc, wtr); err != nil {
		return err
	}
	return ent.StorageInfo.WriteAll(c, cdc, wtr)
}

func (m *AIMessage) GetObjectRef() interface{} { return m.StorageInfo.GetObjectRef() }

func NewAgentMessage(content string) *AIMessage {
	return &AIMessage{
		Role:    StakeholderAgent,
		Content: content,
	}
}

func NewAgentToolMessage(content string, calls []FunctionCall) *AIMessage {
	return &AIMessage{
		Role:          StakeholderAgent,
		Content:       content,
		FunctionCalls: calls,
	}
}

func NewUserAIMessage(content string) *AIMessage {
	return &AIMessage{
		Role:    StakeholderUser,
		Content: content,
	}
}

func NewToolMessage(content string, toolCallID string) *AIMessage {
	return &AIMessage{
		Role:       StakeholderTool,
		Content:    content,
		ToolCallID: toolCallID,
	}
}
