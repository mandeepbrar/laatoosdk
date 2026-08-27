package ctx

import (
	"context"
	"log/slog"
	"time"
)

// Context is the base context every Laatoo context extends: both core.ServerContext and
// core.RequestContext embed it, so everything declared here is available on the ctx a service
// receives without reaching for anything else.
//
// That availability is the point. The recurring cost this interface exists to remove is a plugin
// author writing a helper — or adding a dependency — for something already here: CreateUUID is
// the canonical case, re-implemented from crypto/rand more than once. Read this list before
// writing a utility.
//
// # The params store, and the one distinction worth knowing
//
// Set/Get and the typed accessors read and write a key/value store carried by the context. How
// that store is shared across derived contexts differs by which derivation you use, and the two
// look interchangeable at the call site:
//
//   - SubCtx and the With* methods SHARE the parent's store by reference. A Set on the child is
//     visible to the parent and to every other context derived from it.
//   - NewCtx COPIES the store. A Set on the child is invisible to the parent.
//
// Neither is wrong; picking the wrong one is. A value written into a shared store outlives the
// scope that wrote it, and a value written into a copied store silently fails to propagate.
type Context interface {
	context.Context

	// CreateUUID returns a new random (v4) UUID. This is the platform's id generator: never add
	// a uuid dependency and never hand-roll one from crypto/rand.
	//
	// Note it is unrelated to the context's own id from GetId, which is time-ordered (v1).
	CreateUUID() string

	// GetId returns the context's own trace id, shared by every context derived through SubCtx
	// or With* and freshly minted by NewCtx. Use it to correlate log lines belonging to one
	// request.
	GetId() string

	// GetName returns the name this context was created with. SubCtx and NewCtx set it; the
	// With* methods carry the parent's name through unchanged.
	GetName() string

	// GetPath returns the accumulated chain of names from the root to this context, which is
	// what identifies where in the server hierarchy the context sits. Each derivation appends
	// to it, so it grows as a request descends.
	GetPath() string

	// GetParent returns the context this one was derived from, or nil at the root.
	GetParent() Context

	// Get returns a raw value from the params store. Prefer the typed accessors below unless
	// the value is neither a string, bool, int nor []string.
	Get(key string) (interface{}, bool)

	// GetCreationTime returns when this context was created. SubCtx and the With* methods
	// INHERIT the parent's creation time; NewCtx starts a fresh one — so elapsed time measures
	// the whole request in the first case and only the new scope in the second.
	GetCreationTime() time.Time

	// GetElapsedTime returns the time since GetCreationTime, with the same
	// inherited-versus-fresh distinction.
	GetElapsedTime() time.Duration

	// Set stores a value in the params store. Whether it is visible to the parent depends on
	// how this context was derived — see the note on the interface.
	Set(key string, value interface{})

	// SetVals stores every entry of vals. A nil map is a no-op rather than an error.
	SetVals(vals map[string]interface{})

	// GetString returns a stored value when it is a string. A value of any other type reports
	// not-found rather than being formatted, so the boolean means "present AND a string".
	GetString(key string) (string, bool)

	// GetBool returns a stored bool, and also parses a stored string ("true", "1", …) into one.
	// Note the asymmetry with GetInt, which does NOT accept both.
	GetBool(key string) (bool, bool)

	// GetInt returns a stored value when it is a STRING parseable as an integer.
	//
	// A value stored as an actual int reports not-found — the conversion runs only in the
	// string direction, unlike GetBool. On failure the int returned is -1, not 0, so a caller
	// ignoring the boolean gets a value that looks deliberate.
	GetInt(key string) (int, bool)

	// GetStringArray returns a stored value when it is exactly a []string. No element-wise
	// conversion is attempted, so a []interface{} of strings reports not-found.
	GetStringArray(key string) ([]string, bool)

	// SubCtx derives a named child for a nested scope, SHARING the parent's params store and
	// creation time. Use it to scope logging within one request; use NewCtx when the child's
	// writes must not reach the parent.
	SubCtx(name string) Context

	// NewCtx derives a named child with a COPY of the params store, its own trace id and a
	// fresh creation time — an independent scope rather than a nested one.
	//
	// newpath controls only how the path is composed: true starts the path afresh at name,
	// false appends name to the parent's path.
	NewCtx(name string, newpath bool) Context

	// GetAppengineContext returns the App Engine context associated with the originating
	// request, or nil when the context did not come from one. App Engine specific.
	GetAppengineContext() context.Context

	// GetOAuthContext returns a context carrying the credentials for outbound OAuth-authenticated
	// calls.
	GetOAuthContext() context.Context

	// WithCancel derives a context cancellable through the returned func, sharing the params
	// store. As with the standard library, the caller must call the returned func to release
	// resources.
	WithCancel() (Context, context.CancelFunc)

	// WithDeadline derives a context cancelled at the given time, sharing the params store.
	WithDeadline(timeout time.Time) (Context, context.CancelFunc)

	// WithTimeout derives a context cancelled after the given duration, sharing the params
	// store.
	WithTimeout(timeout time.Duration) (Context, context.CancelFunc)

	// WithValue derives a context carrying a value on the underlying context.Context — the
	// standard library mechanism, reached with Value, and distinct from the params store that
	// Set and Get use.
	WithValue(key, val interface{}) Context

	// WithContext derives a context backed by the given parent context.Context, keeping this
	// context's name, path and params store. Use it to adopt a cancellation or deadline
	// established elsewhere.
	WithContext(parent context.Context) Context

	// CompleteContext records the end of the scope, logging its path and elapsed time.
	//
	// It releases nothing and cancels nothing: a context obtained from WithCancel, WithDeadline
	// or WithTimeout still needs its CancelFunc called. This is a tracing marker, not a
	// teardown.
	CompleteContext()

	// Dump writes every key and value in the params store to the log. A debugging aid — it does
	// not redact, so avoid it on a context carrying credentials.
	Dump()

	// LogTrace logs at the most detailed level. Note it is emitted at the same level as
	// LogDebug, so it cannot be filtered separately.
	LogTrace(msg string, args ...slog.Attr)

	// LogDebug logs diagnostic detail useful while developing.
	LogDebug(msg string, args ...slog.Attr)

	// LogInfo logs a normal, expected event.
	LogInfo(msg string, args ...slog.Attr)

	// LogWarn logs something recoverable that a reader should still notice.
	LogWarn(msg string, args ...slog.Attr)

	// LogError logs a failure.
	LogError(msg string, args ...slog.Attr)

	// LogFatal logs a failure at error level and RETURNS — it does not exit the process,
	// despite the name, and execution continues on the next line. Anything that must stop after
	// a fatal condition has to return or panic explicitly.
	LogFatal(msg string, args ...slog.Attr)
}
