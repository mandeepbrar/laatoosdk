package server

import (
	"laatoo/sdk/core"
)

type SecurityHandler interface {
	core.ServerElement
	HasPermission(core.RequestContext, string) bool
	AuthenticateRequest(ctx core.RequestContext, loadFresh bool) (string, error)
	AllPermissions(ctx core.RequestContext) []string
}
