package elements

import (
	//	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// Engine is a transport that hosts a routing tree: http, websocket, grpc, rpc or mcp. Each engine
// owns a root channel that every channel declared for it descends from, and supplies the response
// handler used for a channel that configures none.
type Engine interface {
	core.ServerElement

	// GetRootChannel returns the engine's root channel -- the value a channel's `parent:` names
	// when it sits directly under the engine.
	//
	// NOT A PURE GETTER: the http and mcp engines construct a FRESH proxy on every call and
	// reassign the root channel's own proxy reference to it, so a Channel obtained earlier and the
	// one returned now are different objects wrapping the same channel.
	GetRootChannel(ctx core.ServerContext) Channel

	// GetRequestParams returns the transport-level parameters of an in-flight request.
	//
	// ONLY THE HTTP ENGINE RETURNS ANYTHING. It returns the route parameters bound by the matched
	// path plus "__requesturi" holding the full request URI. The websocket, grpc, rpc and mcp
	// engines all return nil, so a caller must handle a nil map rather than assume route params
	// are available on every transport.
	GetRequestParams(ctx core.RequestContext) utils.StringMap

	// GetDefaultResponseHandler returns the handler a channel uses when it declares no
	// `responsehandler:`. The channel manager treats a nil result as a Core_Bad_Conf failure at
	// channel start rather than serving without one.
	GetDefaultResponseHandler(ctx core.ServerContext) core.ResponseHandler
}
