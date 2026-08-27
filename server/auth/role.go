package auth

import "laatoo.io/sdk/datatypes"

// Role is a named grant a user carries.
//
// The reference implementation is role.DefaultRole
// (laatoomodules/security/dev/plugins/role/src/server/go/defaultrole.go plus its generated
// autogen_DefaultRole.go). Which concrete type is used is configuration: the security handler
// instantiates whatever the `role` setting names, defaulting to role.DefaultRole, and asserts it
// to this interface unchecked (laatooserver/src/core/securityhandler.go:54-65).
//
// A Role object is not what authorization reads at request time. Enforcement compares the ROLE
// NAME carried on the user's StorableRef against Casbin policy; this entity supplies the
// permission list and the tenant grouping around that.
type Role interface {
	datatypes.Serializable

	// GetId returns the role's storage id.
	//
	// Roles reconciled out of Casbin policy are created with id and name set to the same string
	// (LocalSecurityhandler.go:309-310, securityhandler.go:63-64), and loadRolesFromConf compares
	// GetId against the configured anonymous and admin role names
	// (LocalSecurityhandler.go:669, :677). Keeping id and name equal is a load-bearing convention
	// rather than a schema rule — see SetName.
	GetId() string

	// SetId sets the role's storage id. Used when the platform materialises a role it found in
	// policy but not in the database.
	SetId(string)

	// GetName returns the role name — the value AUTHORIZATION ACTUALLY MATCHES ON.
	//
	// Both enforcement paths call enforcer.Enforce(role.Name, ...) with the Name off the user's
	// StorableRef (LocalSecurityhandler.go:212 and :402), and on the token path those refs are
	// built with Id, Name and the claim string all equal (securityhandler.go:519-524).
	GetName() string

	// SetName sets the role name.
	//
	// A role whose Name differs from its Id fails silently in one direction: policy lines written
	// against one string never match users carrying the other, so the grant simply never applies
	// and nothing reports a mismatch. Every place the platform creates a role sets both to the
	// same value; do the same.
	SetName(string)

	// GetPermissions returns the permission strings this role confers.
	//
	// Consumed when the handler builds the permission set for a list of role names
	// (LocalSecurityhandler.go:570-582, the Append at :579) and when roles are loaded from
	// configuration (:672).
	// These back the `accesspermission` predicate check; they are NOT what gates service access,
	// which is Casbin policy keyed on the role name.
	//
	// Returns a plain slice with no error, unlike DefaultUser's own GetPermissions. A nil slice
	// is a legitimate empty set, not a failure.
	GetPermissions() []string

	// SetPermissions replaces the role's permission list.
	SetPermissions([]string)

	// GetTenant returns the tenant this role belongs to. Like User.GetTenant it NEVER RETURNS NIL
	// in practice — DefaultRole returns &ent.TenantInfo, the address of an embedded struct field
	// (defaultrole.go:109-111).
	//
	// Callers dereference it with no guard: LocalSecurityhandler.go:279 chains straight into
	// GetTenantId, and treats an empty tenant id as meaning the role is global. An implementation
	// that returned a genuine nil would panic there.
	GetTenant() TenantInfo
}
