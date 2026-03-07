package ai

import "laatoo.io/sdk/utils"

// SkillRequest is the canonical invocation envelope for both low-code and pro-code skills.
type SkillRequest struct {
	RequestID   string                `json:"request_id,omitempty" yaml:"request_id,omitempty"`
	SkillID     string                `json:"skill_id,omitempty" yaml:"skill_id,omitempty"`
	Input       string                `json:"input,omitempty" yaml:"input,omitempty"`
	Messages    []ConversationMessage `json:"messages,omitempty" yaml:"messages,omitempty"`
	Context     utils.StringMap       `json:"context,omitempty" yaml:"context,omitempty"`
	ToolName    string                `json:"tool_name,omitempty" yaml:"tool_name,omitempty"`
	ToolArgs    utils.StringMap       `json:"tool_args,omitempty" yaml:"tool_args,omitempty"`
	Model       string                `json:"model,omitempty" yaml:"model,omitempty"`
	Temperature *float32              `json:"temperature,omitempty" yaml:"temperature,omitempty"`
	MaxTokens   *int                  `json:"max_tokens,omitempty" yaml:"max_tokens,omitempty"`
	Metadata    utils.StringMap       `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	SessionID   string                `json:"session_id,omitempty" yaml:"session_id,omitempty"`
	UserID      string                `json:"user_id,omitempty" yaml:"user_id,omitempty"`
}

// SkillResponse is the canonical structured result returned by a skill invocation.
type SkillResponse struct {
	SkillID       string            `json:"skill_id,omitempty" yaml:"skill_id,omitempty"`
	Model         string            `json:"model,omitempty" yaml:"model,omitempty"`
	Content       string            `json:"content,omitempty" yaml:"content,omitempty"`
	FunctionCalls []FunctionCall    `json:"function_calls,omitempty" yaml:"function_calls,omitempty"`
	ToolResults   []utils.StringMap `json:"tool_results,omitempty" yaml:"tool_results,omitempty"`
	Output        utils.StringMap   `json:"output,omitempty" yaml:"output,omitempty"`
	Metadata      utils.StringMap   `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}
