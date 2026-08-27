package elements

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ObjectLoader is the server's type registry. Every type the platform can instantiate by name --
// entities, services, factories, response handlers, principals -- is registered here, and
// CreateObject is the only supported way to build one.
//
// A registered name is derived from the Go type, never chosen: it is the package path with "/"
// replaced by "." plus the type name (e.g. "laatoo.server.engine.http.cookiesResponseHandler").
// Lookups are exact string matches against that form.
//
// REGISTRATION IS FIRST-WRITER-WINS AND SILENT. Both Register and RegisterObjectFactory skip a
// name that is already present and return nil. Nothing is overwritten, nothing is reported above
// trace level, and the version and metadata supplied by the loser are discarded.
type ObjectLoader interface {
	core.ServerElement

	// Register makes obj's type creatable by name, deriving the name from the type itself.
	//
	// A NO-OP, returning nil, when the derived name is already registered -- the first
	// registration keeps the field and this call's version and metadata are dropped. Two plugins
	// shipping the same Go type path therefore resolve to whichever loaded first, with no warning.
	//
	// obj is used only for its reflect.Type; it is never retained as an instance.
	Register(ctx ctx.Context, obj interface{}, version string, metadata core.Info) error

	// RegisterObjectFactory registers a caller-supplied factory instead of letting the loader
	// build one, for types needing construction logic beyond a zero value.
	//
	// The name is derived by CALLING factory.CreateObject once, so registration has the side
	// effect of constructing one instance. Like Register, it silently skips an already-registered
	// name and returns nil.
	RegisterObjectFactory(ctx ctx.Context, factory core.ObjectFactory, version string) error

	// GetRegName returns the registry name the loader would use for object, whether that name is
	// currently registered, and whether the type was reached through a pointer.
	//
	// The name is computed from reflection alone, so it is returned even when the second result is
	// false. Pointers, slices and arrays are unwrapped to the element type before naming.
	GetRegName(ctx ctx.Context, object interface{}) (string, bool, bool)

	// CreateCollection returns a slice of the named type with the given length, as an interface{}
	// the caller type-asserts. Returns Core_Error_Provider_Not_Found for an unregistered name.
	CreateCollection(ctx ctx.Context, objectName string, length int) (interface{}, error)

	// CreateObject returns a new instance of the named type. Returns
	// Core_Error_Provider_Not_Found for an unregistered name -- the most common cause being a
	// name written with "/" separators instead of "." or a plugin that never loaded.
	CreateObject(ctx ctx.Context, objectName string) (interface{}, error)

	// CreateObjectPointersCollection returns a slice of POINTERS to the named type, which is what
	// the data tier needs when results are unmarshalled in place. Returns
	// Core_Error_Provider_Not_Found for an unregistered name.
	CreateObjectPointersCollection(ctx ctx.Context, objectName string, length int) (interface{}, error)

	// GetMetaData returns the metadata registered alongside a type -- the object spec's
	// configurations and request params for a service or factory.
	//
	// RETURNS (nil, nil) FOR AN UNKNOWN NAME, and also for an empty name. The error is never
	// non-nil for a lookup miss, so an `if err != nil` guard passes and the caller proceeds with a
	// nil Info. Callers must nil-check the first result; the server's own call sites do
	// (serviceelement.go loadMetaData, factory.go loadMetaData) and fall back to defaults.
	GetMetaData(ctx ctx.Context, objectName string) (core.Info, error)

	// GetObjectFactory returns the registered factory for a name, and false when there is none.
	GetObjectFactory(ctx ctx.Context, name string) (core.ObjectFactory, bool)

	// List returns the LOADED PLUGINS, keyed by plugin name with the plugin's path as the value --
	// not the registered objects. It is the loader's inventory of what has been loaded, used by the
	// admin service.
	List(ctx core.ServerContext) utils.StringsMap
}
