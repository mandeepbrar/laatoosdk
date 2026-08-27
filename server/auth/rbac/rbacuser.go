package rbac

import (
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/server/components/data"
)

// RbacUser is a User carrying role references — the interface every authorization decision in the
// platform goes through.
//
// It is effectively mandatory, not optional: the server asserts the configured user object to
// RbacUser while rebuilding a caller from a token and fails the request with
// CORE_ERROR_TYPE_MISMATCH if the assertion does not hold
// (laatooserver/src/core/securityhandler.go:506-509). A user object that does not implement this
// cannot authenticate at all.
//
// Roles are held as data.StorableRef values, but only the Name is read at enforcement time — see
// GetRoles.
type RbacUser interface {
	auth.User

	// GetRoles returns the role references this user carries.
	//
	// Both Casbin enforcement paths call it and then compare ref.Name against policy:
	// enforcer.Enforce(role.Name, tenantid, namespace, service) in enforceServiceAccess
	// (laatoomodules/security/.../localsecurityhandler/.../LocalSecurityhandler.go:201-215) and
	// enforcer.Enforce(role.Name, tenantid, predicate) in enforceAccess (:391-405). Id and Type
	// are carried for storage and expansion, not for the decision. On the token path all three
	// are set to the same string (securityhandler.go:519-524).
	//
	// The error return is treated as a hard DENY, never as a retryable failure: both call sites do
	// `if err != nil { return false }` (:206-209, :396-399). An implementation that surfaces a
	// transient lookup failure here silently strips the caller of every role for that request.
	// DefaultUser never returns a non-nil error (defaultuser.go:121-123).
	//
	// A nil slice means no roles and therefore no grants — access is denied, not defaulted.
	GetRoles() ([]data.StorableRef, error)

	// SetRoles replaces the user's role references.
	//
	// Called far more often than GetRoles, and from three distinct situations:
	//
	//   - the server, rebuilding a caller from the token's "Roles" claim — which is a
	//     comma-joined string it splits into refs, done BEFORE LoadClaims
	//     (securityhandler.go:517-527), so LoadClaims must not overwrite roles;
	//   - login services minting a session — keyauthlogin/keyauthservice.go:221,
	//     jwttokenlogin/jwttokenloginservice.go:425, tempauth/temptokenservice.go:108,
	//     signup/accountregister.go:49;
	//   - workflow engines restoring the originating caller's authority on a resumed step
	//     (workflow/goworkflows/.../workflowservice.go:279-281,
	//     workflow/cadence/.../workflowservice.go:173-174).
	//
	// It replaces rather than merges: a second call discards the first call's roles.
	//
	// The parameter type is load-bearing and is dropped silently when it is wrong.
	// DefaultUser.Initialize only applies roles supplied through configuration if the value
	// asserts to []data.StorableRef (defaultuser.go:57-63), so passing []string leaves the user
	// with no roles and no error anywhere — the failure appears later as an unexplained
	// authorization denial. keyauthservice.go:220 carries a comment recording exactly that.
	//
	// The error return is discarded by nearly every caller; do not rely on it being read.
	SetRoles([]data.StorableRef) error
}
