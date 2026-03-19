package elements

import (
	"laatoo.io/sdk/server/core"
)

type SessionManager interface {
	core.ServerElement
	GetSession(ctx core.RequestContext, sessionId string) (core.Session, error)
	GetUserSession(ctx core.RequestContext, userId string) (core.Session, error)
	// SetSession registers an externally constructed session by its ID.
	// Non-serialisable sessions (IsSerializable()==false) are kept only in-process.
	SetSession(ctx core.RequestContext, session core.Session) error
	// DeleteSession removes the session when the SSE connection closes.
	DeleteSession(ctx core.RequestContext, sessionId string) error
}
