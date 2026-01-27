package ai

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// Mcp defines the interface for an MCP Engine wrapper.
// Arguments are defined as interface{} to avoid hard dependency on mcp sdk in this package.
// Implementations should cast them to the appropriate mcp types (e.g. *mcp.Tool, mcp.ToolHandler).
type Mcp interface {
	// Tools
	RemoveTools(ctx core.ServerContext, names ...string) error
	GetTool(ctx core.ServerContext, name string) (Tool, error)
	ListTools(ctx core.ServerContext) (map[string]Tool, error)
	CallTool(ctx core.RequestContext, name string, args utils.StringMap) (interface{}, error)

	// Prompts
	RemovePrompts(ctx core.ServerContext, names ...string) error
	GetPrompt(ctx core.ServerContext, name string) (Prompt, error)
	ListPrompts(ctx core.ServerContext) (map[string]Prompt, error)

	// Resources
	//AddResource(ctx core.ServerContext, name string, resource Resource, handler core.Service) error
	RemoveResources(ctx core.ServerContext, uris ...string) error
	GetResource(ctx core.ServerContext, uri string) (Resource, error)
	ListResources(ctx core.ServerContext) (map[string]Resource, error)
/*
	AddResourceTemplate(ctx core.ServerContext, name string, template ResourceTemplate, handler core.Service) error
	RemoveResourceTemplates(ctx core.ServerContext, templates ...string) error
	GetResourceTemplate(ctx core.ServerContext, template string) (ResourceTemplate, error)
	ListResourceTemplates(ctx core.ServerContext) ([]interface{}, error)*/
}

type Resource interface {
	GetUri() string
	GetName() string
	GetDescription() string
	GetContent() []byte
	GetMimeType() string
	GetMetadata() utils.StringMap
}