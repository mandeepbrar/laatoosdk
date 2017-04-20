package server

import (
	"laatoo/sdk/core"
)

type SecurityHandler interface {
	core.ServerElement
	HasPermission(core.RequestContext, string) bool
	AuthenticateRequest(ctx core.RequestContext) error
	AllPermissions(ctx core.RequestContext) []string
}
