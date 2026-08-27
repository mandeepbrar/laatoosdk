package elements

import (
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

// SecurityHandler is the server element a plugin reaches for security decisions, obtained from a
// context with ServerElementSecurityHandler.
//
// It is the platform-side manager, distinct from components.SecurityHandler, which is the
// pluggable authorization back end a solution configures. The split is by responsibility:
// authentication — verifying a credential and building the caller's identity — lives here and
// nowhere else, while every authorization method below DELEGATES to the configured back end.
//
// That delegation matters when reading the authorization methods: what they actually enforce is
// decided by which back end a solution installed, and one of the two shipped back ends answers
// every question from a single allow-all flag. See components.SecurityHandler.
//
// Most services never call this. Authorization already ran before Invoke — a channel authorizes
// its service on the way in — so reaching for it is for a service making a second, finer-grained
// decision of its own, not for repeating the check that already happened.
type SecurityHandler interface {
	core.ServerElement

	// AuthorizeService decides whether the caller may invoke a service. An empty permission is
	// decided by the service-access policy; a non-empty one routes to HasPermission instead.
	// Both paths refuse a caller holding no roles.
	AuthorizeService(ctx core.RequestContext, module string, service string, permission string, namespace string) (bool, error)

	// ServicesAccessibleByRole lists the services one role can reach — a query over the policy
	// for building a UI or auditing, not an authorization decision.
	ServicesAccessibleByRole(ctx core.RequestContext, role string) ([]string, error)

	// CanAccessObject is the per-object authorization hook, and no shipped back end enforces it:
	// it answers without consulting the object, the id or the action. Enforce row-level access
	// in the service. See components.SecurityHandler.CanAccessObject.
	CanAccessObject(ctx core.RequestContext, module string, service string, object string, objectid string, action string) (bool, error)

	// HasPermission reports whether the caller holds a permission, resolved through its roles. A
	// bare bool: refusal and failure-to-determine are the same value.
	HasPermission(ctx core.RequestContext, permission string) bool

	// AuthenticateRequest authenticates the caller of a request.
	//
	// authComp is the authentication component of the service the request is bound for, or nil.
	// The channel already holds its service reference when authentication runs, so this is passed
	// rather than discovered. When it declares a verification descriptor, the credential is
	// verified against that declaration instead of against the platform's own signing key — which
	// is what allows a machine caller to present a credential this platform did not mint.
	//
	// nil is the ordinary case, not an error: a service that declares nothing authenticates exactly
	// as before, and the task manager and queue consumers authenticate with no channel and
	// therefore no service reference at all.
	AuthenticateRequest(ctx core.RequestContext, authComp components.AuthenticationComponent) (authenticated bool, usrId string, tenant string, token string, claims map[string]interface{}, err error)

	// RegisterPermission declares a permission name so it can be granted and enforced. A plugin
	// registers the permissions it means to check, usually at startup.
	RegisterPermission(ctx core.ServerContext, perm string) error

	// ListPermissions returns the registered permission names.
	ListPermissions(ctx core.ServerContext) []string

	// SaveRole persists a role. Both shipped back ends accept the call and store nothing.
	SaveRole(ctx core.RequestContext, role interface{}) error

	// GetRole looks up a role by name. Both shipped back ends return a nil Role AND a nil error
	// for every name, so guard on the returned Role rather than on err.
	//
	// A role needs no entry here to work: authorization matches role names as strings against
	// the policy.
	GetRole(ctx core.RequestContext, name string) (auth.Role, error)

	// ListRoles returns the roles known for the caller's tenant, falling back to the global
	// tenant when the request carries none.
	ListRoles(ctx core.RequestContext) (map[string]auth.Role, error)

	// ListServices returns the services known to the policy.
	ListServices(ctx core.ServerContext) []string

	// AddServiceAccessPolicy grants a role access to a service at runtime.
	//
	// The tenant argument is carried by the signature but not honoured by the Casbin back end,
	// which enforces with an empty tenant — so a grant naming a tenant matches nobody and only a
	// wildcard grant takes effect. The fallback back end discards the call and returns nil.
	AddServiceAccessPolicy(ctx core.ServerContext, tenant string, module string, service string, role string, permission string, namespace string) error

	// AddRolePermissionPolicy grants a permission to a role at runtime, with the same tenant
	// caveat as AddServiceAccessPolicy.
	AddRolePermissionPolicy(ctx core.ServerContext, tenant string, module string, service string, role string, permission string) error
}
