package ai

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// Mcp defines the interface for an MCP Engine wrapper.
// Arguments are defined as interface{} to avoid hard dependency on mcp sdk in this package.
// Implementations should cast them to the appropriate mcp types (e.g. *mcp.Tool, mcp.ToolHandler).
//
// THIS INTERFACE IS AN INSPECTION VIEW, NOT THE REGISTRATION PATH. There is no Add* method:
// tools, prompts and resources are declared as `mcp` channels in a plugin's
// registry/channels/, and the engine registers them with the underlying MCP SDK server at
// startup. Obtain the wrapper with AgentManager.GetMCPServer(rootpath).
//
// Only the TOOLS half of this view is actually maintained. Registering a prompt channel writes
// straight to the MCP SDK server (mcpchannel.go:457) and never touches the engine's own prompt
// map, and resource registration is commented out entirely (mcpimpl.go:384-466) — so the
// prompt and resource accessors below report an empty world even on a server serving prompts
// over the wire. See the per-method notes.
type Mcp interface {
	// RemoveTools unregisters tools by name and drops them from the inspection map. Names that
	// were never registered are ignored. A panic inside the underlying MCP SDK is recovered
	// and returned as an error, so a non-nil error here does not mean nothing was removed.
	RemoveTools(ctx core.ServerContext, names ...string) error

	// GetTool returns a registered tool by its MCP tool name, or an error if none matches.
	//
	// The name is derived from the CHANNEL PATH, not the service: "/math/add" registers as
	// "math_add" (mcpchannel.go:148). Passing the service name finds nothing. Tools are
	// registered at channel-initialize time, so this works before the engine starts serving.
	GetTool(ctx core.ServerContext, name string) (Tool, error)

	// ListTools returns every registered tool keyed by MCP tool name. The error is always nil.
	//
	// It hands back the engine's LIVE internal map, not a copy — mutating the returned map
	// mutates the tool registry. Iterating it while channels are still initializing is a data
	// race.
	ListTools(ctx core.ServerContext) (map[string]Tool, error)

	// CallTool invokes a tool in-process and returns its response Data.
	//
	// THREE THINGS TO KNOW BEFORE USING IT. First, `name` is the SERVICE name — it goes
	// straight to ctx.Invoke(name, args) and never consults the tool map — so it is a
	// different namespace from GetTool/ListTools/RemoveTools, which key on the MCP tool name.
	// Second, it returns (nil, nil) when the service leaves no response (mcpimpl.go:228), so an
	// `if err != nil` guard passes and the caller nil-derefs. Third, and worst, it does NOT
	// apply the channel's `staticvalues` — the params the channel pins server-side and strips
	// from the tool's public input schema. The MCP wire path merges them before invoking
	// (laatoomcptool.go:58); this path does not, so a caller can supply values the wire path
	// forbids and pinned values are simply missing. It also does not check resp.Status, so a
	// service rejection (a 400 with no Go error) is returned as an ordinary result.
	// Nothing in the platform calls this method.
	CallTool(ctx core.RequestContext, name string, args utils.StringMap) (interface{}, error)

	// RemovePrompts unregisters prompts by name from the underlying MCP SDK server. This part
	// does work — unlike the accessors below — because it forwards to the SDK server directly;
	// the accompanying delete from the engine's own map is a no-op on an always-empty map.
	RemovePrompts(ctx core.ServerContext, names ...string) error

	// GetPrompt returns a registered prompt by name.
	//
	// ALWAYS FAILS WITH "prompt not found". The engine's prompt map is allocated and never
	// written: registration from an mcp channel calls the MCP SDK server directly
	// (mcpchannel.go:457), and the only line that would populate the map lives inside the
	// commented-out AddPrompt (mcpimpl.go:347). A not-found here says nothing about whether
	// the prompt is being served.
	GetPrompt(ctx core.ServerContext, name string) (Prompt, error)

	// ListPrompts returns every registered prompt. ALWAYS RETURNS AN EMPTY MAP AND A NIL
	// ERROR, for the reason given on GetPrompt — an empty result is not evidence that the
	// server has no prompts. To enumerate what is really registered, drive a real MCP session
	// against the SDK server.
	ListPrompts(ctx core.ServerContext) (map[string]Prompt, error)

	// RemoveResources unregisters resources by URI. Forwards to the underlying MCP SDK server
	// and then deletes from an always-empty map; since nothing ever registers a resource (see
	// GetResource), this has nothing to remove.
	RemoveResources(ctx core.ServerContext, uris ...string) error

	// GetResource returns a registered resource by URI.
	//
	// ALWAYS FAILS WITH "resource not found". MCP resources are not implemented: the engine's
	// AddResource is commented out in full (mcpimpl.go:384-466), so the resource map is
	// allocated at startup and never written by any code path, and there is no channel kind
	// that registers one.
	GetResource(ctx core.ServerContext, uri string) (Resource, error)

	// ListResources returns every registered resource. ALWAYS RETURNS AN EMPTY MAP AND A NIL
	// ERROR — resources are unimplemented; see GetResource.
	ListResources(ctx core.ServerContext) (map[string]Resource, error)
	/*
		AddResourceTemplate(ctx core.ServerContext, name string, template ResourceTemplate, handler core.Service) error
		RemoveResourceTemplates(ctx core.ServerContext, templates ...string) error
		GetResourceTemplate(ctx core.ServerContext, template string) (ResourceTemplate, error)
		ListResourceTemplates(ctx core.ServerContext) ([]interface{}, error)*/
}

// Resource is a static or generated document an MCP client can read by URI.
//
// NOTHING IMPLEMENTS THIS INTERFACE AND NOTHING CONSTRUCTS ONE. As of 2026-08-27 no type in
// laatoo, laatoosdk or the solutions has a GetUri method, the engine's AddResource is
// commented out (laatooserver/src/engine/mcp/mcpimpl.go:384-466), and no channel kind produces
// a resource. The method contracts below are the interface's stated intent, NOT observed
// behaviour — they are unverified, and no nil, empty or encoding convention should be assumed
// until MCP resources are actually implemented.
type Resource interface {
	// GetUri returns the resource's URI. Intended to be the registry key and the identifier an
	// MCP client reads by.
	GetUri() string

	// GetName returns a short human-readable name for the resource.
	GetName() string

	// GetDescription returns a human-readable description of the resource's contents.
	GetDescription() string

	// GetContent returns the resource body as raw bytes.
	GetContent() []byte

	// GetMimeType returns the body's MIME type.
	GetMimeType() string

	// GetMetadata returns arbitrary additional data carried with the resource.
	GetMetadata() utils.StringMap
}
