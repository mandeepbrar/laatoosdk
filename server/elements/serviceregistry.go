package elements

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
)

// ServiceRegistry is the CROSS-POD service catalogue: it publishes a pod's services to the cluster
// and resolves services hosted by other pods.
//
// IT IS OPT-IN. A registry exists only when a level's config declares `serviceregistry: {type:
// nats}`; without that, the service manager is local-only and never consults one, however many
// pods are running. The only implemented type is "nats", backed by the embedded NATS server's
// KeyValue store plus request-reply; an unknown type is a Core_Bad_Conf failure at boot rather
// than a silent fall back to local-only.
//
// This is not a core.ServerElement and is not published into contexts -- the service manager holds
// it and is the only caller.
type ServiceRegistry interface {
	// RegisterService publishes a service to the cluster catalogue AND binds a request-reply
	// handler that will execute it for any caller who can reach its subject. This is a security
	// boundary, which is why the registry is opt-in: a registered service answers remote
	// invocations.
	//
	// The bound handler rebuilds the caller's identity from the wire -- principal, roles and
	// tenant -- and authorization is ENFORCED against it. An identity that was claimed but cannot
	// be rebuilt is REFUSED rather than downgraded to anonymous, since a nil caller would reach
	// the system-request path where authorization is off.
	//
	// RE-REGISTERING AN ALIAS SILENTLY REPLACES: any existing local subscription for that alias is
	// unsubscribed and the catalogue entry is overwritten, with no error and no warning.
	RegisterService(ctx core.ServerContext, serviceAlias string, svc Service, conf config.Config) error

	// GetService resolves a service hosted anywhere in the cluster, returning a proxy that
	// dispatches over NATS request-reply.
	//
	// RETURNS (nil, nil, nil) WHEN THE SERVICE IS NOT IN THE CATALOGUE -- a miss is not an error,
	// so an `if err != nil` guard passes and the caller must nil-check the Service. The service
	// manager relies on exactly this to fall through to its own "service missing" error.
	//
	// The config.Config result is ALWAYS nil from the NATS implementation: a remote service's
	// configuration lives on the pod hosting it and is not carried across.
	//
	// The returned proxy is degenerate in places -- Service() and GetModule() return nil,
	// IsStreaming() is false, and each call carries a fixed 10-second timeout.
	GetService(ctx core.ServerContext, serviceName string) (Service, config.Config, error)
}
