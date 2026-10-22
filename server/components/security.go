package components

import (
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/server/components/data"
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
//   - The Casbin-backed handler enforces service access and resource access permissions from the
//     role permissions and roles the namespace authors under config/security/ (and, for customer
//     roles, from the database). This is what a real deployment uses.
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
	// Both paths read the caller's ROLE PERMISSIONS, which authentication stored on the request
	// context under constants.ROLEPERMISSIONS, and enforce the ones this back end defines. Both
	// refuse a caller that holds none, so a service cannot grant itself access by any means
	// available to it at request time.
	AuthorizeService(ctx core.RequestContext, module string, service string, permission string, namespace string) (bool, error)

	// AddServiceAccessPolicy grants a role access to a service at runtime.
	//
	// The Casbin handler refuses the call with an error: access is granted by the role permissions
	// a namespace authors and the roles that list them, never by a row added to a role directly.
	// The fallback handler discards the call entirely and returns nil, which is indistinguishable
	// from success.
	AddServiceAccessPolicy(ctx core.ServerContext, tenant string, module string, service string, role string, permission string, namespace string) error

	// RegisterPermission hands the back end one resource access permission id that a plugin of
	// this namespace declared. The server calls it once per declared id, before ApplyAccessModel,
	// so the back end can check the role permissions it loads against the ids that exist.
	RegisterPermission(ctx core.ServerContext, perm string) error

	// ListPermissions returns the resource access permission ids registered with the back end.
	ListPermissions(ctx core.ServerContext) []string

	// ApplyAccessModel loads the namespace's role permissions and roles and checks them against what
	// has been registered. The server calls it after every module that registers declarations with
	// this back end has started, and before the namespace serves requests.
	//
	// Every resource access permission a role permission lists must have been registered, every
	// service it lists must exist, every role must list role permissions that exist, and no role
	// permission may reuse an id in ancestorRolePermissions, the ids defined by namespaces above this
	// one. An unresolved name is returned as ONE error naming each file, line and name, so a single
	// boot reports every problem. The call is idempotent: a second call rebuilds the model from the
	// same files.
	ApplyAccessModel(ctx core.ServerContext, ancestorRolePermissions []string) error

	// ResolveRolePermissions expands role names into the ids of the role permissions this back end
	// holds for them, read in the tenant of the request. A role this back end does not know
	// contributes nothing and is not an error. The server calls it for every request, at
	// authentication, and stores the union of the answers of every back end up the namespace tree
	// on the request context under constants.ROLEPERMISSIONS.
	ResolveRolePermissions(ctx core.RequestContext, roles []string) ([]string, error)

	// ListRolePermissions returns the role permissions this back end defines. Those defined above
	// this namespace are listed by the server element, which adds them.
	ListRolePermissions(ctx core.RequestContext) ([]auth.RolePermission, error)

	// DeleteRole removes a role built in role management, in the tenant of the request. A role the
	// namespace ships in its files cannot be deleted and is an error. The next request of a user who
	// held the role finds it unknown.
	DeleteRole(ctx core.RequestContext, name string) error

	// HasPermission reports whether the caller holds a resource access permission, through the role
	// permissions on the request context. A caller holding no role permission holds no permissions.
	//
	// It returns a bare bool: "not permitted" and "could not be determined" are the same value.
	HasPermission(ctx core.RequestContext, permission string) bool

	// AddRolePermissionPolicy grants a permission to a role at runtime. The Casbin handler refuses
	// it for the reason AddServiceAccessPolicy is refused.
	AddRolePermissionPolicy(ctx core.ServerContext, tenant string, module string, service string, role string, permission string) error

	// CanAccessObject is the hook for per-object authorization — whether this caller may perform
	// this action on this particular record.
	//
	// NO SHIPPED HANDLER ENFORCES IT. The Casbin handler returns true unconditionally, and the
	// fallback returns its allow-all flag; neither looks at the object, the id or the action.
	// Calling this and branching on the result gives a check that reads like a control and is not
	// one.
	//
	// Record-level authorization is EntityAccessFilter's job, not this method's: an entity that
	// opts in has its entityaccess policy applied by its data component on every read and write,
	// lists included, which a per-object yes/no cannot serve.
	CanAccessObject(ctx core.RequestContext, module string, service string, object string, objectid string, action string) (bool, error)

	// RegisterEntityAccessPolicies hands the handler the entityaccess policy rows a module ships,
	// as CSV file contents keyed by the file's path (used only to name the file in an error).
	//
	// It REPLACES every row previously registered for that module: a hot-reloaded module
	// re-registers and its old rows are gone, and an unloaded module registers an empty map to
	// clear them. A row the handler cannot enforce — a rule outside the pushable grammar — is
	// refused here with an error naming the file, the row and the construct, so boot fails rather
	// than an entity serving unprotected.
	RegisterEntityAccessPolicies(ctx core.ServerContext, module string, policies map[string][]byte) error

	// EntityAccessFilter answers which records of an entity the caller may reach for an action,
	// reduced to a predicate over the record's fields. See data.EntityAccessFilter.
	//
	// A system context is answered as unrestricted. A caller whose roles match no row is answered
	// as denied. An error is a failure to decide and must be treated as a refusal, never as
	// permission.
	EntityAccessFilter(ctx core.RequestContext, entity string, action data.EntityAction) (*data.EntityAccessFilter, error)

	// EnforcesEntityAccess reports whether this handler implements record-level access at all.
	//
	// A data component whose entity opts in checks it at initialization and refuses to start
	// when it is false: the fallback handler answers false, and an opted-in entity under it must
	// fail boot rather than be served without its rules.
	EnforcesEntityAccess() bool

	// SetClaims is a hook for a handler to add claims to a user before a token is minted. Both
	// shipped handlers implement it as an empty body.
	//
	// It is not how extra claims reach a token in practice: the server's token generator merges
	// the claims passed at mint time, independently of this.
	SetClaims(user auth.User, addClaims map[string]interface{}, exp int64)

	// SaveRole persists a role built in role management, in the tenant of the request. The role
	// must implement auth.RolePermissionHolder; every role-permission id it lists must be defined
	// by this back end or by one above it, and its name must not be the name of a role the
	// namespace ships. The change applies to the next request, with no restart. The fallback
	// handler accepts the call and stores nothing.
	SaveRole(ctx core.RequestContext, role interface{}) error

	// GetRole looks up a role by name, in the tenant of the request.
	//
	// The Casbin handler returns the role when it is a default role or a role built in role
	// management, and (nil, nil) when it knows no such role — a nil Role WITH a nil error, so a
	// caller guarding only on err proceeds to dereference nothing. Check the returned Role. The
	// fallback handler returns (nil, nil) for every name.
	GetRole(ctx core.RequestContext, name string) (auth.Role, error)

	// ListRoles returns the roles known for the caller's tenant, falling back to the global
	// tenant when the request carries none: the roles the namespace ships and the roles built in
	// role management. The fallback returns nil.
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
