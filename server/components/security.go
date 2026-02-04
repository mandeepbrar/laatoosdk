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

type AuthenticationComponent interface {
	SetTokenGenerator(core.ServerContext, func(auth.User, map[string]interface{}, int64) (string, auth.User, error))
}
