package elements

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ServiceManager owns the services visible at one level of the server hierarchy (server, solution,
// application, isolation) and their lifecycle -- creation from config and modules, initialization,
// start, and unload.
//
// LOOKUP DOES NOT WALK UP TO A PARENT. Each level's manager holds its own flat store, checks that,
// then the cross-pod registry if one is configured, and then fails. Levels inherit instead by COPY:
// a child manager is built with a snapshot of its parent's store taken at creation time
// (laatooserver/src/core/abstractserver_create.go:567-587). A service registered on a parent AFTER
// its children were built is therefore invisible to those children.
//
// SERVICE ALIAS COLLISIONS ARE SILENT. Creating a service under an alias that already exists
// returns the PRE-EXISTING service with no error and no warning
// (laatooserver/src/core/servicemanager.go:264-268) -- the second declaration's config, factory and
// object are all discarded. Two modules declaring the same alias resolve to whichever was processed
// first, and nothing in the boot log says so.
type ServiceManager interface {
	core.ServerElement

	// GetService resolves a service by its alias -- the service YAML filename minus ".yml", unless
	// the file declares `name:`.
	//
	// Checks this level's own store, then the remote (NATS) service registry when
	// `serviceregistry:` is configured, then returns Core_Error_Missing_Service. There is no walk
	// up the element hierarchy. A service resolved remotely is NOT cached: it is a
	// request-reply proxy over NATS and is re-resolved on every lookup, deliberately, so a service
	// is not pinned to wherever it happened to live at first call.
	GetService(ctx core.ServerContext, alias string) (Service, error)

	// GetServiceContext returns the server context of the named service -- the context its config,
	// cache and logger were resolved in. Resolves the service exactly as GetService does and
	// propagates its error unchanged.
	GetServiceContext(ctx core.ServerContext, alias string) (core.ServerContext, error)

	// List returns every service alias visible at this level mapped to the name of the module that
	// declared it. Services declared outside a module carry an empty module name.
	List(ctx core.ServerContext) utils.StringsMap

	// Describe returns an admin-facing snapshot of one service: Name, Object (the servicemethod),
	// Conf, LogLevel, Params and Permission. Returns a not-found error for an unknown alias --
	// including for a service that exists only in the remote registry, which this does not consult.
	Describe(ctx core.ServerContext, svc string) (utils.StringMap, error)

	// ChangeLogger changes one service's log level and format at runtime, with no restart.
	// Returns a not-found error for an unknown alias; like Describe, it only sees this level's
	// local store.
	ChangeLogger(ctx core.ServerContext, svc string, logLevel string, logFormat string) error

	// CreateParams builds request Param descriptors from a `params:` config block, reading type,
	// description, collection, stream and required for each entry.
	//
	// An UNRECOGNISED type name is not an error: it becomes datatypes.Object with the literal
	// string kept as the custom object type to instantiate, so the failure surfaces later as a
	// missing object rather than here as a bad config.
	CreateParams(ctx core.ServerContext, paramsConf config.Config) (map[string]core.Param, error)
}
