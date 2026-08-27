package elements

import (
	"laatoo.io/sdk/server/core"
)

// Isolation is the innermost level of the server hierarchy -- server, solution, application,
// isolation -- and represents one live data context of an application, typically a tenant.
//
// An isolation is created from a directory under the SOLUTION's isolations/, not under the
// application, and it has its own module, service, channel and factory managers seeded from copies
// of the application's.
type Isolation interface {
	core.ServerElement

	// GetIsolationId returns the id of the thing this isolation exists for -- the tenant id, for a
	// tenant isolation. It comes from the isolation config's required "IsoId" key, so it is
	// distinct from the isolation's element NAME (the directory name), which GetName returns; the
	// two are commonly different and code that needs the tenant must use this one.
	//
	// The value is fixed at creation and never empty: a missing "IsoId" fails the isolation's
	// creation with Core_Missing_Conf.
	GetIsolationId() string
	//GetApplet(name string) (Applet, bool)
}
