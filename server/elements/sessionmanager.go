package elements

import (
	"laatoo.io/sdk/server/core"
)

type SessionManager interface {
	core.ServerElement
	GetSession(ctx core.RequestContext, sessionId string) (core.Session, error)
	GetUserSession(ctx core.RequestContext, userId string) (core.Session, error)
}
