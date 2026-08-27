package ai

import "laatoo.io/sdk/server/core"

// ToolAnnotations provide hints about tool behavior
type ToolAnnotations struct {
	ReadOnly     bool `json:"read_only"`
	Destructive  bool `json:"destructive"`
	Idempotent   bool `json:"idempotent"`
	OpenWorld    bool `json:"open_world"` // accesses external resources
	RequiresAuth bool `json:"requires_auth"`
}

// Tool is a service exposed to an LLM as a callable function. It is a plain
// core.UserInvokableService plus behavioural hints; the platform builds the tool's JSON input
// schema from the service's own object-spec request params, not from anything declared here.
type Tool interface {
	core.UserInvokableService

	// Annotations returns behavioural hints about this tool.
	//
	// DO NOT GATE ON THIS. The only shipped implementation — laatooMcpTool, the wrapper the MCP
	// engine puts around every service exposed on an mcp channel — returns a zero-valued
	// ToolAnnotations{} constant, so ReadOnly, Destructive, Idempotent, OpenWorld and
	// RequiresAuth all read false for every tool in the platform regardless of what the tool
	// actually does. A caller that refuses to run when Destructive is true never refuses
	// anything, and a caller that trusts ReadOnly trusts every tool. There is no config key
	// that populates these fields.
	// See laatooserver/src/engine/mcp/laatoomcptool.go:28.
	Annotations() ToolAnnotations
}

// ToolCallMode controls tool calling behavior
type ToolCallMode string

const (
	ToolCallNone     ToolCallMode = "none"
	ToolCallAuto     ToolCallMode = "auto"
	ToolCallRequired ToolCallMode = "required"
)
