package errors

import (
	"log/slog"
	"sync"

	"laatoo.io/sdk/ctx"
)

// Kind is the kind of failure a declared error reports. It decides the status a client is
// answered with and the plain sentence a client reads when the error itself is not for it.
type Kind string

const (
	// KindBadInput is a request the caller can correct: a missing, malformed or refused value.
	KindBadInput Kind = "badinput"
	// KindNotSignedIn is a request with no valid credential.
	KindNotSignedIn Kind = "notsignedin"
	// KindNotPermitted is a signed-in caller asking for what it may not do.
	KindNotPermitted Kind = "notpermitted"
	// KindNotFound is a request for a record or resource that does not exist.
	KindNotFound Kind = "notfound"
	// KindConflict is a request that collides with what already exists, such as a taken name.
	KindConflict Kind = "conflict"
	// KindNotSetUp is a failure caused by missing or wrong configuration rather than by the caller.
	KindNotSetUp Kind = "notsetup"
	// KindUnavailable is a dependency that cannot be reached: a database, a queue, a remote service.
	KindUnavailable Kind = "unavailable"
	// KindInternal is a fault in the platform's own code.
	KindInternal Kind = "internal"
)

// kinds is every Kind a declaration may name.
var kinds = map[Kind]bool{
	KindBadInput: true, KindNotSignedIn: true, KindNotPermitted: true, KindNotFound: true,
	KindConflict: true, KindNotSetUp: true, KindUnavailable: true, KindInternal: true,
}

// Valid reports whether k is one of the declared kinds. Pure.
func (k Kind) Valid() bool {
	return kinds[k]
}

// Audience says who may read a declared error.
type Audience string

const (
	// AudiencePerson is an error sent to the client: its code, field, values and sentence.
	AudiencePerson Audience = "person"
	// AudienceLog is an error that stays in the log; a client reads only its kind's plain sentence.
	AudienceLog Audience = "log"
)

// Valid reports whether a is one of the declared audiences. Pure.
func (a Audience) Valid() bool {
	return a == AudiencePerson || a == AudienceLog
}

// Declaration is one declared error as a plugin (or the server) states it: in a registry/errors
// file, or under errors: in a service's configuration. The sentence comes from the declaring
// module's default property file, under errors.<Code>.
type Declaration struct {
	// Module is the module instance whose plugin declares the error.
	Module string
	// Code names the error within its plugin, e.g. DesignNameTaken.
	Code string
	// Description says when the error is raised, for a reader of the declaration.
	Description string
	// Kind is the kind of failure.
	Kind Kind
	// Audience says who may read it.
	Audience Audience
	// Field is the request field the error concerns, or "".
	Field string
	// Values names the values a raise supplies and the sentence may use.
	Values []string
	// Sentence is the default-language sentence template, or "" when the property is missing.
	Sentence string
}

// DeclaredErrorKey is the info key under which a raised declared error carries its *Raised.
const DeclaredErrorKey = "DeclaredError"

// Raised is what one raise of a declared error records. Declaration is nil when the code was not
// declared for the raising module, or when no resolver is set.
type Raised struct {
	// Code is the code the raise named.
	Code string
	// Values are the named values the raise supplied.
	Values []slog.Attr
	// Cause is the error that led to this one, or nil; it is logged, never sent to a client.
	Cause error
	// Declaration is the resolved declaration, or nil.
	Declaration *Declaration
}

// DeclarationResolver finds the declaration of code for the module ctx belongs to. The server
// sets one at start-up; without it every raise is unresolved.
type DeclarationResolver func(c ctx.Context, code string) (*Declaration, bool)

var (
	resolverLock sync.RWMutex
	resolver     DeclarationResolver
)

// SetDeclarationResolver installs the function a raise resolves its code with. The server calls
// it once while starting; a later call replaces the earlier one.
func SetDeclarationResolver(r DeclarationResolver) {
	resolverLock.Lock()
	defer resolverLock.Unlock()
	resolver = r
}

// resolve answers code's declaration for the module c belongs to, or nil.
func resolve(c ctx.Context, code string) *Declaration {
	resolverLock.RLock()
	r := resolver
	resolverLock.RUnlock()
	if r == nil {
		return nil
	}
	if d, ok := r(c, code); ok {
		return d
	}
	return nil
}

// Raise returns the declared error code of the module c belongs to, carrying the named values,
// e.g. errors.Raise(ctx, "DesignNameTaken", slog.String("kind", "agent"), slog.String("name", n)).
func Raise(c ctx.Context, code string, values ...slog.Attr) error {
	return raise(c, code, nil, values)
}

// RaiseFrom is Raise with the error that caused it, which is kept for the log and never sent to a
// client, e.g. a database driver's error under the platform's own "could not reach the database".
func RaiseFrom(c ctx.Context, cause error, code string, values ...slog.Attr) error {
	return raise(c, code, cause, values)
}

// raise builds the declared error: an *Error whose code is the declared code and whose info carries
// the *Raised, the values and, through throwError, the cause and the stack.
func raise(c ctx.Context, code string, cause error, values []slog.Attr) error {
	raised := &Raised{Code: code, Values: values, Cause: cause, Declaration: resolve(c, code)}
	message := "Declared error " + code
	if raised.Declaration != nil && raised.Declaration.Description != "" {
		message = raised.Declaration.Description
	}
	info := append([]slog.Attr{slog.Any(DeclaredErrorKey, raised)}, values...)
	return throwError(c, message, code, cause, info...)
}

// DeclaredError answers the declared error err carries as its own outermost error, or nil. Like
// HasErrorCode it reads the error it is handed; the server's reply path walks every wrapping.
func DeclaredError(err error) *Raised {
	laatooErr, ok := err.(*Error)
	if !ok || laatooErr == nil {
		return nil
	}
	for _, attr := range laatooErr.info {
		if attr.Key != DeclaredErrorKey {
			continue
		}
		if raised, ok := attr.Value.Any().(*Raised); ok {
			return raised
		}
	}
	return nil
}

// IsKind reports whether err is a resolved declared error of kind k.
func IsKind(err error, k Kind) bool {
	raised := DeclaredError(err)
	return raised != nil && raised.Declaration != nil && raised.Declaration.Kind == k
}
