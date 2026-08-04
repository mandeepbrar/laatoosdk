package auth

// MachinePrincipal is a User that is a machine caller rather than a person — a task worker, a
// workflow callback, or an external system holding a credential.
//
// It extends User rather than replacing it, the same way rbac.RbacUser adds roles: everything that
// already accepts a User accepts a machine caller unchanged, so no existing call site, context
// accessor or authorization path needs to know this type exists. Code that does care asserts to
// MachinePrincipal, exactly as 32 sites already assert to rbac.RbacUser.
//
// The consequence to hold onto: because a machine satisfies User, it will answer the person-shaped
// methods with empty values. What keeps machine callers out of user listings and user-scoped role
// resolution is the realm they authenticate into, not the type — a surface that must exclude them
// filters by realm, or asserts to MachinePrincipal, rather than assuming every User is a person.
type MachinePrincipal interface {
	User
	// GetClientId returns the identifier the credential presented for this caller, as issued by
	// whoever registered it. Distinct from GetId, which is the platform's identifier.
	GetClientId() string
	// GetIssuer returns the issuer that minted this caller's credential. For an externally-minted
	// credential this is the external issuer; for a platform-minted one it is the platform.
	GetIssuer() string
}
