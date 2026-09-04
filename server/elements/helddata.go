package elements

import (
	"laatoo.io/sdk/server/components/data"
	"laatoo.io/sdk/server/core"
)

// The data layer's addressable elements.
//
// These follow the shape Service, Channel, Factory and Module already use, and the distinction
// that shape encodes is the one thing to understand before adding another: AN ELEMENT WRAPS AN
// IMPLEMENTATION, IT IS NOT THE IMPLEMENTATION.
//
// core.Service is not an element. elements.Service is, and its Service() accessor reaches the
// plugin's implementation behind it. The same split applies here: data.DataComponent stays exactly
// what a storage provider implements, and elements.DataComponent is the handle the server
// registers, resolves and reports an address for.
//
// Keeping it that way is not a preference. Making data.DataComponent itself an element would add
// eight methods to an interface implemented by every storage provider, reaching datacommon's
// BaseComponent, data.DataPlugin and any hand-written component -- a plugin-tree change for a
// server-side capability. It would also make a struct deserialized from YAML (data.Dataset) carry
// a context and a parent pointer it has no use for.

// DataComponent is the server's handle on a registered data component -- the element the
// DataManager stores under a (dataconnection, object) key, and which a dataset compiles against.
//
// It wraps, but is not, the provider's own data.DataComponent implementation: Component() reaches
// that.
type DataComponent interface {
	core.ServerElement

	// Component returns the provider's implementation behind this element.
	Component() data.DataComponent

	// GetConnectionName returns the dataconnection this component is bound to -- a factory
	// instance name, which is what a dataconnection is. It is the other half of the key the
	// DataManager registry is built on, and it is never empty: a component is always created by a
	// factory.
	GetConnectionName() string

	// GetObject returns the registered object type this component stores.
	GetObject() string
}

// Dataset is the server's handle on a loaded dataset definition -- the element the DataManager
// stores by dataset id and a /ds/ request resolves.
//
// It wraps, but is not, the *data.Dataset the YAML loader produced: Definition() reaches that.
// The definition stays a plain configuration struct, which is why it is safe to hand to a
// serializer and safe to build in a test with no server running.
type Dataset interface {
	core.ServerElement

	// Definition returns the dataset definition behind this element.
	//
	// Returned by pointer and NOT copied, so a caller must treat it as read-only: the same
	// definition backs every resolution of this dataset.
	Definition() *data.Dataset

	// GetEntity returns the entity this dataset queries.
	GetEntity() string
}
