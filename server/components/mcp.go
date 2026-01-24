package components

import (
	"laatoo.io/sdk/server/core"
)

// Mcp defines the interface for an MCP Engine wrapper.
// Arguments are defined as interface{} to avoid hard dependency on mcp sdk in this package.
// Implementations should cast them to the appropriate mcp types (e.g. *mcp.Tool, mcp.ToolHandler).
type Mcp interface {
	// Tools
	AddTool(ctx core.ServerContext, tool interface{}, handler interface{}) error
	RemoveTools(ctx core.ServerContext, names ...string) error

	// Prompts
	AddPrompt(ctx core.ServerContext, prompt interface{}, handler interface{}) error
	RemovePrompts(ctx core.ServerContext, names ...string) error

	// Resources
	AddResource(ctx core.ServerContext, resource interface{}, handler interface{}) error
	RemoveResources(ctx core.ServerContext, uris ...string) error
	AddResourceTemplate(ctx core.ServerContext, template interface{}, handler interface{}) error
	RemoveResourceTemplates(ctx core.ServerContext, templates ...string) error
}
