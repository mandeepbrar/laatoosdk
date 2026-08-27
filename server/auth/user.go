package auth

import "laatoo.io/sdk/datatypes"

// TenantInfo identifies the tenant a user or role belongs to.
//
// The only implementation in the platform is the data.TenantInfo struct
// (server/components/data/tenantinfo.go), which entity codegen embeds BY VALUE into every
// multitenant entity — user.DefaultUser and role.DefaultRole among them. Because it is embedded
// by value, the accessors returning it hand back the address of a struct field and are therefore
// never nil. See User.GetTenant for what that costs.
//
// Serializable carries TenantId and TenantName as two plain strings, so a tenant survives a round
// trip through storage but not through a token: PopulateClaims writes only the id.
type TenantInfo interface {
	datatypes.Serializable

	// GetTenantId returns the tenant identifier the data layer partitions on.
	//
	// The empty string is NOT a detectable absence — it is a value that gets filtered on. A
	// multitenant read stamps args["TenantId"] = GetTenantId() unconditionally
	// (laatoomodules/datastore/dev/plugins/common/basecomponent.go:212), so a caller with no
	// tenant reads exactly the rows whose TenantId is also empty, and ValidateTenant
	// (basecomponent.go:566) grants ownership of them. Authorization is worse: the tenant column
	// passed to Casbin is a local hardcoded to ""
	// (laatoomodules/security/.../localsecurityhandler/.../LocalSecurityhandler.go:196-200 and
	// :386-390), so role grants are compared against the empty tenant no matter who is calling.
	//
	// Test a user for "no tenant" with GetTenant().GetTenantId() == "", never against nil.
	GetTenantId() string

	// GetTenantName returns the tenant's display name.
	//
	// A label only — nothing in the platform authorizes or partitions on it. It is also not
	// carried on a token: LoadClaims rebuilds the tenant with SetTenant(id, ""), so the name is
	// permanently empty for the rest of a token-authenticated request.
	GetTenantName() string
}

// User is the identity of the caller of a request.
//
// The reference implementation is user.DefaultUser
// (laatoomodules/security/dev/plugins/user/src/server/go/defaultuser.go). A User reaches a service
// by one of two very different routes, and they populate different subsets of this interface:
//
//   - loaded from the data store during login — every field populated (except Password, which
//     PostLoad blanks);
//   - reconstructed from a verified token by the server on each subsequent request
//     (laatooserver/src/core/securityhandler.go:504-543) — a FRESH object carrying only what
//     PopulateClaims wrote plus the roles the server sets, with every other field at its zero
//     value.
//
// On the token path GetEmail and GetRealm are always "", and GetUserAccount is always nil,
// because nothing restores them. A service needing those must re-read the user by GetId.
type User interface {
	datatypes.Serializable

	// GetId returns the user's storage id — the value services key on and the only field the
	// platform treats as the identity.
	//
	// On the token path it is set by LoadClaims from claims["UserId"], which the token generator
	// writes separately from PopulateClaims (securityhandler.go:337). A credential that leaves it
	// empty is refused rather than admitted as an authenticated caller with no id
	// (securityhandler.go:537).
	GetId() string

	// SetId sets the user's storage id. Used while constructing a user rather than during a
	// request: from configuration in Initialize, and by LoadClaims from the token's subject.
	SetId(string)

	// GetUsernameField returns the NAME of the entity field holding the login name — not the
	// login name. DefaultUser returns "Username" (defaultuser.go:93-95).
	//
	// Callers build a lookup condition from it, pairing it with GetUserName:
	// utils.StringMap{usr.GetUsernameField(): username}
	// (laatoomodules/security/dev/plugins/common/common.go:66, signup/accountregister.go:73,
	// oauthlogin/oauthloginservice.go:201).
	//
	// It must name a real, queryable field on the user entity. A wrong name is not an error at
	// any layer — the condition simply matches no rows, and every login is rejected as bad
	// credentials with nothing in the logs pointing at the field name.
	GetUsernameField() string

	// GetUserName returns the login name. Carried on a token: PopulateClaims writes
	// claims["UserName"] and LoadClaims reads it back, so this is one of the few fields populated
	// on the token path.
	GetUserName() string

	// LoadClaims populates this user from the claims of an already-verified credential.
	//
	// DefaultUser reads exactly six keys and nothing else: UserId, UserName, Name, Roles, Account
	// and Tenant (defaultuser.go:196-233). A token whose claims use other names — an externally
	// minted one naming its subject "sub" — yields a user with an empty id, which the server then
	// refuses (securityhandler.go:537). The names are a platform convention, not a standard.
	//
	// Nothing outside those six is restored. Email, Realm, Permissions and Status are left at
	// their zero values for the whole request.
	//
	// Three of DefaultUser's assertions are UNGUARDED and panic inside the authentication path if
	// the claim is present but not a string: claims["UserName"].(string) at defaultuser.go:204,
	// claims["Name"].(string) at :207 and claims["Tenant"].(string) at :231. A token carrying a
	// number or an object under any of those keys takes down the request handler rather than
	// failing the credential.
	//
	// The Roles branch is dead on the JWT path: the token generator writes Roles as a
	// comma-joined string (securityhandler.go:333), so the []data.StorableRef assertion at
	// defaultuser.go:212 never succeeds. Roles are set by the server calling SetRoles just before
	// this (securityhandler.go:517-527), and an implementation that "fixes" LoadClaims to parse
	// them would be overwriting that.
	//
	// Implementations must guard every read: this runs on attacker-supplied input.
	LoadClaims(map[string]interface{})

	// PopulateClaims writes the claims a generated token will carry. DefaultUser writes UserName,
	// Name, Tenant and Account (defaultuser.go:164-169).
	//
	// It is NOT the whole token, and must not be treated as the place the identity is carried:
	// the generator adds claims["UserId"] itself immediately afterwards, along with Roles and exp
	// (securityhandler.go:333-343). It is also not symmetric with LoadClaims — Email and Realm
	// are written by neither, which is why they never survive a token.
	//
	// Whatever is written here lands in the JWT payload, which is signed but NOT encrypted and is
	// readable by anyone holding the token. Never write a password, hash, key or private field.
	PopulateClaims(map[string]interface{})

	// GetEmail returns the user's stored email address (defaultuser.go:171-173).
	//
	// Not a claim: on a token-authenticated request this is always "". Do not identify or notify
	// a caller from it inside a service — re-read the user by GetId first.
	GetEmail() string

	// GetRealm returns the authentication realm recorded on the user (defaultuser.go:130-132).
	//
	// Nothing currently gates on it: every realm comparison in the platform is commented out
	// (laatoomodules/security/dev/plugins/common/common.go:57,
	// signup/accountregister.go:66). Like GetEmail it is not a claim, so it is "" on the token
	// path. Treat it as informational until something reinstates the check.
	GetRealm() string

	// GetTenant returns the caller's tenant. IN PRACTICE IT NEVER RETURNS NIL.
	//
	// DefaultUser returns &usr.TenantInfo — the address of an embedded struct field
	// (defaultuser.go:89-91) — and DefaultRole does the same. So every `GetTenant() == nil` guard
	// in the platform is dead code, and a user built with no tenant surfaces as an EMPTY tenant
	// id rather than a detectable absence.
	//
	// Three such guards, all in laatoomodules/datastore/dev/plugins/common/basecomponent.go:
	// RequireTenantScope (:196) exists to refuse a multitenant read by a caller with no tenant
	// and never fires; PreProcessConditionMap (:212) then scopes the query to TenantId "";
	// ValidateTenant (:562) grants ownership of any row whose tenant is also "".
	//
	// Users in exactly that state are routine, not exotic: the Anonymous and System users the
	// server builds at startup (securityhandler.go:207-238) and machine users minted by
	// keyauthlogin (keyauthservice.go:212-218) are all created without tenant info.
	//
	// Check GetTenant().GetTenantId() == "" instead. An implementation returning a genuine nil
	// would panic several unguarded callers, e.g. LocalSecurityhandler.go:279.
	GetTenant() TenantInfo

	// GetUserAccount returns the user's profile entity, or nil.
	//
	// DefaultUser returns usr.Account.Entity, and nil when that is unset (defaultuser.go:139-144).
	// Entity is populated only when the reference has been expanded by a read — it is
	// `json:"-"` (server/components/data/storageinfo.go:138) so deserialization never fills it,
	// and LoadClaims rebuilds only the ref's Id/Type/Name (defaultuser.go:217-228). Nil is
	// therefore the normal answer on a token-authenticated request.
	//
	// Callers must nil-check. laatoomodules/security/dev/plugins/mfa/.../mfaemailtask.go:129-130
	// does not, and calls GetFullName straight through.
	//
	// DefaultUser's assertion Account.Entity.(auth.UserAccount) at defaultuser.go:143 is also
	// unguarded, and no account entity that ships implements UserAccount — see the doc on that
	// interface. Expanding the reference is what would make this panic.
	GetUserAccount() UserAccount
}
