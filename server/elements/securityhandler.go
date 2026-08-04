package elements

import (
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

type SecurityHandler interface {
	core.ServerElement
	AuthorizeService(ctx core.RequestContext, module string, service string, permission string, namespace string) (bool, error)
	ServicesAccessibleByRole(ctx core.RequestContext, role string) ([]string, error)
	CanAccessObject(ctx core.RequestContext, module string, service string, object string, objectid string, action string) (bool, error)
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
	RegisterPermission(ctx core.ServerContext, perm string) error
	ListPermissions(ctx core.ServerContext) []string
	SaveRole(ctx core.RequestContext, role interface{}) error
	GetRole(ctx core.RequestContext, name string) (auth.Role, error)
	ListRoles(ctx core.RequestContext) (map[string]auth.Role, error)
	ListServices(ctx core.ServerContext) []string
	AddServiceAccessPolicy(ctx core.ServerContext, tenant string, module string, service string, role string, permission string, namespace string) error
	AddRolePermissionPolicy(ctx core.ServerContext, tenant string, module string, service string, role string, permission string) error
}
