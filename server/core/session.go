package core

import (
	"time"

	"laatoo.io/sdk/utils"
)

// Session is a keyed bag of per-conversation state, held by the session manager. Two
// implementations ship: laatooserver/src/core.session, the ordinary cacheable one, and
// engine/http.AGUISession, which holds a live SSE stream and is in-process only.
//
// Every typed getter delegates to utils.StringMap and uses a STRICT type assertion with no
// coercion (utils/stringmap.go:7-77). A key that is present but stored under a different concrete
// type returns the same (zero, false) as a key that is absent — the two are indistinguishable to
// the caller. That matters here because a session round-trips through the cache, and a value that
// went in as an int can come back as something else. When in doubt use Get and assert yourself.
//
// Session data is a plain Go map with no lock (session.go's `data utils.StringMap`), so a session
// shared across concurrent requests must not be written from more than one goroutine.
type Session interface {
	// GetId returns the session id. A session created with an empty id is given a UUID v1
	// (laatooserver/src/core/session.go:25-30); otherwise it keeps the id it was looked up by.
	GetId() string

	// CreationTime returns when this session object was constructed — time.Now() at construction,
	// not when the user authenticated and not a last-touched time. Nothing refreshes it.
	CreationTime() time.Time

	// GetUser returns the session's user id.
	//
	// IT IS ALWAYS EMPTY. The backing field is declared and read but assigned nowhere in the
	// platform (laatooserver/src/core/session.go:19 declares `userid`, :40 returns it, and there is
	// no third occurrence), and AGUISession returns the empty string as a literal
	// (engine/http/aguisession.go:44). Never gate anything on it. The caller's identity comes from
	// ctx.GetUser(), which is populated by the auth layer.
	GetUser() string

	// GetString returns a value stored as a Go string. Anything else stored under the key — a
	// []byte, a fmt.Stringer, a number — yields ("", false) rather than being converted.
	GetString(key string) (string, bool)

	// GetBool returns a value stored as a Go bool. A "true" stored as a string yields
	// (false, false), which reads identically to an absent key.
	GetBool(key string) (bool, bool)

	// GetInt returns a value stored as a Go int.
	//
	// The zero return is -1, NOT 0 (utils/stringmap.go:51), so a caller that ignores the ok flag
	// gets a negative number rather than a plausible-looking zero. The assertion is to int
	// exactly: an int64, or a number that arrived via a JSON decode as float64, yields
	// (-1, false).
	GetInt(key string) (int, bool)

	// GetStringArray returns a value stored as []string, or as []interface{} whose every element
	// is a string. A single non-string element makes the WHOLE call return (nil, false), not a
	// partial slice.
	GetStringArray(key string) ([]string, bool)

	// AllKeys returns the session's keys in map-iteration order, which is randomised — do not
	// depend on it. An empty session returns an empty slice, never nil.
	AllKeys() []string

	// GetStringMap returns a value stored as utils.StringMap or map[string]interface{}.
	//
	// BUG, do not rely on the third form: utils.StringMap.GetStringMap builds a converted map from
	// a map[interface{}]interface{} value and then never returns it, falling through to
	// (nil, false) instead (utils/stringmap.go:104-116). YAML-decoded nested maps take that shape,
	// so a legitimately present value reports as absent.
	GetStringMap(key string) (utils.StringMap, bool)

	// GetStringsMap returns a value stored as map[string]string, or as map[string]interface{} whose
	// every value is a string. As with GetStringArray, one non-string value fails the whole call.
	GetStringsMap(key string) (utils.StringsMap, bool)

	// Get returns the raw stored value with no type handling — the escape hatch when a typed getter
	// returns false and you need to see what is actually there.
	//
	// AGUISession special-cases the key "responseStream", returning its live stream and reporting
	// ok only when that stream is non-nil (engine/http/aguisession.go:47-53).
	Get(key string) (interface{}, bool)

	// Set stores a value, overwriting any existing one. It does not persist: the session is only
	// written out when Save is called.
	//
	// On AGUISession, setting "responseStream" to a value that is not a core.ResponseStream is
	// SILENTLY DISCARDED — the setter returns without storing anything and without error
	// (aguisession.go:55-62).
	Set(key string, val interface{})

	// SetVals stores every entry of vals, overwriting per key; keys already in the session and
	// absent from vals are left alone. A nil map is a no-op.
	SetVals(vals utils.StringMap)

	// Save writes the session back to the session manager's in-memory store and, when the session
	// is serialisable, to the cache.
	//
	// IT RETURNS nil EVEN WHEN THE CACHE WRITE FAILED. sessionManager.save logs the cache error at
	// WARN and returns nil regardless (laatooserver/src/core/sessionmanager.go:109-115), so an
	// `if err != nil` guard here proves only that the in-memory store was updated. On a session
	// constructed without a manager, Save is a no-op returning nil (session.go:76-81), and
	// AGUISession.Save is unconditionally a no-op returning nil (aguisession.go:75).
	Save(RequestContext) error

	// IsSerializable returns false for sessions holding live non-serialisable objects (e.g. SSE
	// streams). The session manager stores these only in its in-process sync.Map and never writes
	// them to the CacheComponent (sessionmanager.go:72-80).
	//
	// A false here means the session does not survive the process: another replica looking it up
	// finds nothing in its own in-memory store, misses the cache, and silently CREATES A NEW EMPTY
	// SESSION under the same id rather than reporting the miss (sessionmanager.go:98-107).
	IsSerializable() bool
}
