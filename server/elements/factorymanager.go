package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// FactoryManager owns the service factories at one level of the server hierarchy. Every service is
// created BY a factory: a service YAML naming no `factory:` gets the built-in DefaultFactory, so
// the factory manager is initialized and started before the service manager.
//
// Like the service manager, its store is flat per level and a child level receives a filtered COPY
// of its parent's factories at creation time -- there is no walk up the hierarchy at lookup.
// Registration is likewise first-writer-wins: creating a factory under an alias that already exists
// returns the existing one with no error (laatooserver/src/core/factorymanager.go:209-212).
type FactoryManager interface {
	core.ServerElement

	// GetFactory returns a factory by alias. Returns a not-found error when this level has none by
	// that name.
	//
	// Two aliases always exist: "DefaultFactory", used by every service that declares no
	// `factory:`, and "activity", used for workflow activities.
	GetFactory(ctx core.ServerContext, factoryName string) (Factory, error)

	// List returns every factory alias at this level mapped to the name of the module that declared
	// it. Factories registered by the server itself carry an empty module name.
	List(ctx core.ServerContext) utils.StringsMap

	// ChangeLogger changes one factory's log level and format at runtime, with no restart. Despite
	// the parameter name, the argument is a FACTORY alias. Returns a not-found error for an
	// unknown alias.
	ChangeLogger(ctx core.ServerContext, chanName string, logLevel string, logFormat string) error

	// Describe returns an admin-facing snapshot of one factory: Name, Conf, LogLevel, Object (the
	// registered Go type it was built from) and ModuleName. Despite the parameter name, the
	// argument is a FACTORY alias. Returns a not-found error for an unknown alias.
	Describe(ctx core.ServerContext, chanName string) (utils.StringMap, error)
}
