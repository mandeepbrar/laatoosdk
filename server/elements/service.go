package elements

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// Service is the server's handle on a configured service instance -- the element the service
// manager stores, channels bind to, and other services are injected as.
//
// It wraps, but is not, the plugin's own core.Service implementation: Service() reaches that.
// A service resolved from the cross-pod NATS registry also satisfies this interface, and that
// remote form answers several methods degenerately (see each method).
type Service interface {
	core.ServerElement

	// Metadata returns the service's descriptor -- request params, configurations, description,
	// streaming flag -- built from the object spec YAML and then amended by the service's own
	// Describe hook.
	Metadata() core.ServiceInfo

	// Service returns the plugin's core.Service implementation behind this element.
	//
	// RETURNS nil for a service resolved from the remote (NATS) registry, which has no local
	// implementation to return.
	Service() core.Service

	// GetDescription returns the service description, taken from the object spec and overridable
	// at runtime through the service's Describe hook. For a remote service it is a generated
	// "Remote NATS Proxy for <alias>" string.
	GetDescription() string

	// GetModule returns the module that supplied this service, or nil for a service declared in
	// server, solution or application configuration rather than by a plugin -- and always nil for
	// a remote service.
	GetModule() core.Module

	// ServiceContext returns the server context the service was created in. This is the context
	// its configuration, cache selection, logger and middleware were resolved against, and the one
	// the channel derives per-request contexts from.
	ServiceContext() core.ServerContext

	// GetConfiguration returns the raw service YAML config, which includes keys the object spec
	// never declared. Reading a key from here is not the same as reading it via
	// GetStringConfiguration on the service, which sees only declared `configurations:`.
	GetConfiguration() config.Config

	// HandleRequest runs one invocation: it AUTHORIZES the call against the service's
	// accesspermission, builds a request, populates its params from vals according to the object
	// spec's declared request params, and invokes the middleware chain.
	//
	// Calling this does not bypass security -- an unauthorized call returns
	// Core_Error_Unauthorized. Two behaviours to know:
	//   - The middleware chain STOPS at the first service that sets a response, so a middleware
	//     that responds prevents the target service from running at all.
	//   - A middleware entry that does not implement core.UserInvokableService is skipped
	//     silently rather than failing.
	//
	// vals keys not declared in the object spec's request params are ignored; encoding names a
	// codec per parameter for values arriving as bytes.
	HandleRequest(ctx core.RequestContext, vals utils.StringMap, encoding utils.StringsMap) error
}
