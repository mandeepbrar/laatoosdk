package ai

import "laatoo.io/sdk/server/core"

// Example shows concrete application patterns
type Example struct {
	Description    string                 `json:"description"`
	Input          map[string]interface{} `json:"input,omitempty"`
	ExpectedOutput string                 `json:"expected_output,omitempty"`
	Notes          string                 `json:"notes,omitempty"`
}

// Instruction is a storable block of prompt text, optionally parameterised as a template.
//
// NOTHING IN THE PLATFORM IMPLEMENTS THIS INTERFACE. As of 2026-08-27 there is no type
// anywhere in laatoo, laatoosdk or the solutions with a GetText/IsTemplate/Params method set,
// and the only code that names the type is the MCP engine's prompt registry — whose AddPrompt
// is commented out (laatooserver/src/engine/mcp/mcpimpl.go:232-349). The semantics below are
// therefore the interface's stated intent, not observed behaviour; treat them as unverified
// and do not rely on any particular nil, empty or error convention until a real implementation
// exists.
type Instruction interface {
	core.Storable

	// GetText returns the instruction body. When IsTemplate reports true this is the
	// un-rendered template source, with placeholders for the names in Params.
	GetText() string

	// GetDescription returns a human-readable description of what the instruction is for.
	GetDescription() string

	// IsTemplate reports whether GetText must be rendered with parameter values before use.
	IsTemplate() bool

	// Params returns the parameters the template accepts, keyed by parameter name.
	Params() map[string]core.Param
}

// Prompt is an Instruction used as an MCP prompt. It is a type alias in intent only — Go
// treats it as a distinct named interface type with identical method set.
//
// The MCP registry side of this is dead: mcpEngine.prompts is allocated but never written,
// because the only line that would populate it lives inside the commented-out AddPrompt
// (mcpimpl.go:347). Prompts registered from an mcp channel go straight to the underlying MCP
// SDK server (mcpchannel.go:457) and never appear here. Consequently Mcp.GetPrompt always
// reports "prompt not found" and Mcp.ListPrompts always returns an empty map, no matter how
// many prompt channels the solution declares.
type Prompt Instruction
