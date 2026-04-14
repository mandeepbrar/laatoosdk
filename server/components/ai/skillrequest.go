package ai

import "laatoo.io/sdk/utils"

// SkillRequest is the canonical invocation envelope for both low-code and pro-code skills.
type SkillRequest struct {
	RequestID   string                `json:"requestId,omitempty" yaml:"requestId,omitempty"`
	SkillID     string                `json:"skillId,omitempty" yaml:"skillId,omitempty"`
	Input       string                `json:"input,omitempty" yaml:"input,omitempty"`
	Messages    []ConversationMessage `json:"messages,omitempty" yaml:"messages,omitempty"`
	Context     utils.StringMap       `json:"context,omitempty" yaml:"context,omitempty"`
	ToolName    string                `json:"toolName,omitempty" yaml:"toolName,omitempty"`
	ToolArgs    utils.StringMap       `json:"toolArgs,omitempty" yaml:"toolArgs,omitempty"`
	Model       string                `json:"model,omitempty" yaml:"model,omitempty"`
	Temperature *float32              `json:"temperature,omitempty" yaml:"temperature,omitempty"`
	MaxTokens   *int                  `json:"maxTokens,omitempty" yaml:"maxTokens,omitempty"`
	Metadata    utils.StringMap       `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	SessionID   string                `json:"sessionId,omitempty" yaml:"sessionId,omitempty"`
	UserID      string                `json:"userId,omitempty" yaml:"userId,omitempty"`
}

// SkillResponse is the canonical structured result returned by a skill invocation.
type SkillResponse struct {
	SkillID       string            `json:"skillId,omitempty" yaml:"skillId,omitempty"`
	Model         string            `json:"model,omitempty" yaml:"model,omitempty"`
	Content       string            `json:"content,omitempty" yaml:"content,omitempty"`
	FunctionCalls []FunctionCall    `json:"functionCalls,omitempty" yaml:"functionCalls,omitempty"`
	ToolResults   []utils.StringMap `json:"toolResults,omitempty" yaml:"toolResults,omitempty"`
	Output        utils.StringMap   `json:"output,omitempty" yaml:"output,omitempty"`
	Metadata      utils.StringMap   `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}
