package ai

import "laatoo.io/sdk/utils"

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
