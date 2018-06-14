package server

import (
	"laatoo/sdk/core"
)

type SessionManager interface {
	core.ServerElement
	GetSession(ctx core.ServerContext, sessionId string) (core.Session, error)
	GetUserSession(ctx core.ServerContext, userId string) (core.Session, error)
	Broadcast(core.ServerContext, func(core.ServerContext, core.Session) error) error
}
