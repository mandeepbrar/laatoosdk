package components

import (
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// SecurityHandler is the authorization back end: the pluggable half of security that decides
// whether a caller may reach a service or hold a permission. Authentication — verifying who the
// caller is — is not here; that lives in the server's own security manager.
//
// A solution selects one by naming it in the security configuration's authorizationhandler.
//
// # Which implementation is in use decides whether anything is enforced
//
// Two ship, and they are not interchangeable:
//
//   - The Casbin-backed handler enforces service access and permissions from the policy files
//     under config/security/. This is what a real deployment uses.
//   - A fallback handler is installed when no authorizationhandler is named. It holds a single
//     allow-all flag and answers EVERY authorization question with it, consulting no policy at
//     all. With the flag off it denies everything; with it on it permits everything.
//
// The fallback is not a degraded version of the first — it is a different thing wearing the same
// interface. A solution that meant to configure a handler and did not gets uniform answers rather
// than an error, so verify from the boot log which one was installed before concluding that a
// policy is being applied.
//
// # Not every method is implemented by both
//
// Several are hooks that no shipped handler fills in. Each is called out below. Treating one of
// those as an enforced control is the dangerous reading, so they say plainly what they return.
type SecurityHandler interface {
	// InitializeProps receives the handler's configuration properties at startup.
	InitializeProps(ctx core.ServerContext, props utils.StringMap) error

	// ServicesAccessibleByRole lists the services one role can reach, for building a UI or
	// auditing a policy. It is a query over the policy, not an authorization decision.
	ServicesAccessibleByRole(ctx core.RequestContext, role string) ([]string, error)

	// ListServices returns the services known to the policy.
	ListServices(ctx core.ServerContext) []string

	// AuthorizeService decides whether the caller may invoke a service.
	//
	// The permission argument selects between two different checks rather than adding to one:
	// empty means the service declared no permission and access is decided by the service-access
	// policy; non-empty routes to HasPermission and the service-access policy is not consulted.
	//
	// Both paths read the caller's ROLES and both refuse a caller that has none, so a service
	// cannot grant itself access by any means available to it at request time.
	AuthorizeService(ctx core.RequestContext, module string, service string, permission string, namespace string) (bool, error)

	// AddServiceAccessPolicy grants a role access to a service at runtime.
	//
	// The tenant argument is accepted by this signature but the Casbin handler enforces with an
	// empty tenant, so a grant scoped to a named tenant matches nobody — only a wildcard grant
	// takes effect. The fallback handler discards the call entirely and returns nil, which is
	// indistinguishable from success.
	AddServiceAccessPolicy(ctx core.ServerContext, tenant string, module string, service string, role string, permission string, namespace string) error

	// RegisterPermission declares a permission name so it can be granted and enforced.
	RegisterPermission(ctx core.ServerContext, perm string) error

	// ListPermissions returns the registered permission names.
	ListPermissions(ctx core.ServerContext) []string

	// HasPermission reports whether the caller holds a permission, resolved through the caller's
	// roles. A caller with no roles holds no permissions.
	//
	// It returns a bare bool: "not permitted" and "could not be determined" are the same value.
	HasPermission(ctx core.RequestContext, permission string) bool

	// AddRolePermissionPolicy grants a permission to a role at runtime, with the same tenant
	// caveat as AddServiceAccessPolicy.
	AddRolePermissionPolicy(ctx core.ServerContext, tenant string, module string, service string, role string, permission string) error

	// CanAccessObject is the hook for per-object authorization — whether this caller may perform
	// this action on this particular record.
	//
	// NO SHIPPED HANDLER ENFORCES IT. The Casbin handler returns true unconditionally, and the
	// fallback returns its allow-all flag; neither looks at the object, the id or the action.
	// Row-level authorization must therefore be enforced by the service itself. Calling this and
	// branching on the result gives a check that reads like a control and is not one.
	CanAccessObject(ctx core.RequestContext, module string, service string, object string, objectid string, action string) (bool, error)

	// SetClaims is a hook for a handler to add claims to a user before a token is minted. Both
	// shipped handlers implement it as an empty body.
	//
	// It is not how extra claims reach a token in practice: the server's token generator merges
	// the claims passed at mint time, independently of this.
	SetClaims(user auth.User, addClaims map[string]interface{}, exp int64)

	// SaveRole persists a role. Both shipped handlers accept the call and store nothing.
	SaveRole(ctx core.RequestContext, role interface{}) error

	// GetRole looks up a role by name.
	//
	// Both shipped handlers return (nil, nil) for every name — a nil Role WITH a nil error, so a
	// caller guarding only on err proceeds to dereference nothing. Check the returned Role.
	//
	// A role does not need to exist here to be usable: authorization matches role names as
	// strings against the policy, so a role named in a grant works whether or not this can
	// resolve it.
	GetRole(ctx core.RequestContext, name string) (auth.Role, error)

	// ListRoles returns the roles known for the caller's tenant, falling back to the global
	// tenant when the request carries none. The Casbin handler answers from roles held in
	// memory; the fallback returns nil.
	ListRoles(ctx core.RequestContext) (map[string]auth.Role, error)
}

// AuthenticationComponent is implemented by a service that participates in authentication —
// either by minting a token for a caller it has authenticated, or by declaring how credentials
// presented on its own channel are verified.
type AuthenticationComponent interface {
	// SetTokenGenerator injects the platform's token-issuing function at startup, for services
	// registered under the security handler's authservices configuration.
	SetTokenGenerator(core.ServerContext, func(auth.User, map[string]interface{}, int64) (string, auth.User, error))

	// GetVerificationDescriptor returns how an inbound machine credential presented on this
	// service's channel is verified, and whether one is declared at all.
	//
	// Returning false is the ordinary case and means "I declare nothing" — the request then
	// authenticates exactly as it does today, against the platform's own signing key. A service
	// that returns true is verified against its declaration instead, which is what lets a caller
	// present a credential this platform never minted.
	//
	// The service declares; it does not verify. Key resolution and the signature check happen in
	// the security manager, so credential cryptography lives in one audited place and no service
	// holds key material.
	GetVerificationDescriptor(core.ServerContext) (*VerificationDescriptor, bool)
}
