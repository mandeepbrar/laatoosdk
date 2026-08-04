package components

import (
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type SecurityHandler interface {
	InitializeProps(ctx core.ServerContext, props utils.StringMap) error
	ServicesAccessibleByRole(ctx core.RequestContext, role string) ([]string, error)
	ListServices(ctx core.ServerContext) []string
	AuthorizeService(ctx core.RequestContext, module string, service string, permission string, namespace string) (bool, error)
	AddServiceAccessPolicy(ctx core.ServerContext, tenant string, module string, service string, role string, permission string, namespace string) error

	RegisterPermission(ctx core.ServerContext, perm string) error
	ListPermissions(ctx core.ServerContext) []string
	HasPermission(ctx core.RequestContext, permission string) bool
	AddRolePermissionPolicy(ctx core.ServerContext, tenant string, module string, service string, role string, permission string) error

	CanAccessObject(ctx core.RequestContext, module string, service string, object string, objectid string, action string) (bool, error)
	SetClaims(user auth.User, addClaims map[string]interface{}, exp int64)

	SaveRole(ctx core.RequestContext, role interface{}) error
	GetRole(ctx core.RequestContext, name string) (auth.Role, error)
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
