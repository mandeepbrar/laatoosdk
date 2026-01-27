package ai 

import "laatoo.io/sdk/server/core"

// ToolAnnotations provide hints about tool behavior
type ToolAnnotations struct {
	ReadOnly     bool `json:"read_only"`
	Destructive  bool `json:"destructive"`
	Idempotent   bool `json:"idempotent"`
	OpenWorld    bool `json:"open_world"`     // accesses external resources
	RequiresAuth bool `json:"requires_auth"`
}

type Tool interface {
	core.UserInvokableService
	Annotations() ToolAnnotations
}

// ToolCallMode controls tool calling behavior
type ToolCallMode string

const (
	ToolCallNone     ToolCallMode = "none"
	ToolCallAuto     ToolCallMode = "auto"
	ToolCallRequired ToolCallMode = "required"
)
