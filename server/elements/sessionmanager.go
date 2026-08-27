package elements

import (
	"laatoo.io/sdk/server/core"
)

// SessionManager holds sessions for one server level, across two stores: an in-process map that is
// always consulted first, and a cache component behind it. A session that reports
// IsSerializable()==false lives ONLY in-process — which is what lets it hold live objects such as
// an SSE response stream — and therefore does not survive a restart and is invisible to other
// replicas (laatooserver/src/core/sessionmanager.go:256-268).
type SessionManager interface {
	core.ServerElement

	// GetSession returns the session with this id, checking the in-process store and then the
	// session cache.
	//
	// IT NEVER REPORTS A MISS. On a miss it CREATES a new empty session under the requested id,
	// stores it in both stores, and returns it (sessionmanager.go:336-345). A non-nil result is
	// therefore not evidence that the session existed, and an attacker-supplied or stale id
	// produces a fresh session rather than an error — check the session's own contents, not this
	// method's error, to decide whether a session is real.
	GetSession(ctx core.RequestContext, sessionId string) (core.Session, error)

	// GetUserSession is intended to return the session belonging to a user.
	//
	// IT IS AN UNIMPLEMENTED STUB: the implementation is `return nil, nil` for every input
	// (sessionmanager.go:355-357). Worse, the stub returns a nil *session through a core.Session
	// result, so the caller receives a NON-NIL INTERFACE HOLDING A NIL POINTER
	// (laatooserver/src/core/sessionmanagerproxy.go:15-17) — `if sess != nil` passes and the next
	// method call runs on a nil receiver. Do not call this method.
	GetUserSession(ctx core.RequestContext, userId string) (core.Session, error)
	// SetSession registers an externally constructed session by its ID.
	// Non-serialisable sessions (IsSerializable()==false) are kept only in-process.
	SetSession(ctx core.RequestContext, session core.Session) error
	// DeleteSession removes the session when the SSE connection closes.
	//
	// IT ONLY CLEARS THE IN-PROCESS STORE. The implementation deletes from inMemoryStore and
	// returns nil without touching the session cache
	// (laatooserver/src/core/sessionmanager.go:321-324), so a SERIALISABLE session comes straight
	// back from the cache on the next GetSession — deletion appears to succeed and does not stick.
	// It is correct for the non-serialisable sessions it was written for (SSE streams), which live
	// nowhere else.
	DeleteSession(ctx core.RequestContext, sessionId string) error
}
