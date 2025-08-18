package elements

import (
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/server/core"
)

type SecurityHandler interface {
	core.ServerElement
	AuthorizeService(ctx core.RequestContext, module string, service string, permission string) (bool, error)
	ServicesAccessibleByRole(ctx core.RequestContext, role string) ([]string, error)
	CanAccessObject(ctx core.RequestContext, module string, service string, object string, objectid string, action string) (bool, error)
	HasPermission(ctx core.RequestContext, permission string) bool
	AuthenticateRequest(ctx core.RequestContext) (authenticated bool, usrId string, tenant string, token string, claims map[string]interface{}, err error)
	RegisterPermission(ctx core.ServerContext, perm string) error
	ListPermissions(ctx core.ServerContext) []string
	SaveRole(ctx core.RequestContext, role interface{}) error
	GetRole(ctx core.RequestContext, name string) (auth.Role, error)
	ListRoles(ctx core.RequestContext) (map[string]auth.Role, error)
	ListServices(ctx core.ServerContext) []string
	AddServiceAccessPolicy(ctx core.ServerContext, tenant string, module string, service string, role string, permission string) error
	AddRolePermissionPolicy(ctx core.ServerContext, tenant string, module string, service string, role string, permission string) error
}
