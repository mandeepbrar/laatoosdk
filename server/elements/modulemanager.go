package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ModuleManager owns the modules loaded at one server level (server, solution, application,
// isolation) and is the introspection surface behind the admin service. Modules are loaded from
// configuration during Initialize; this interface does not create or destroy them.
//
// EVERY METHOD IS SCOPED TO THIS LEVEL'S OWN MODULE SET. A module loaded at the solution level is
// not listed by an application's module manager, so a miss here means "not at this level", not
// "not on this server".
type ModuleManager interface {
	core.ServerElement

	// List returns every module loaded at this level as moduleInstanceName -> pluginName
	// (laatooserver/src/core/modulemanager.go:312-321). The keys are module instance names, which
	// differ from plugin names whenever a plugin is instantiated more than once or under an alias.
	List(ctx core.ServerContext) utils.StringsMap

	// Describe returns a snapshot of one module: Name, Conf, LogLevel, Object, PluginName,
	// PluginDir, Settings, Properties, and the module's loaded Services, Factories, Channels,
	// Rules and Tasks (modulemanager.go:352-373). Returns a NotFound error for an unknown module.
	//
	// The map is a description for display, not a live handle — the config values in it are the
	// module's own config objects and must not be mutated.
	Describe(ctx core.ServerContext, mod string) (utils.StringMap, error)

	// ChangeLogger retunes one module's logging AT RUNTIME, without a restart, by reconfiguring
	// that module's server context logger (modulemanager.go:342-349). Returns NotFound for an
	// unknown module.
	//
	// The change applies to the module's context and everything created under it. Empty strings
	// are passed through to the logger rather than being treated as "leave unchanged", so pass the
	// values you actually want.
	ChangeLogger(ctx core.ServerContext, mod string, logLevel string, logFormat string) error

	// GetModule returns the PLATFORM'S module implementation object for a module instance —
	// the core.Module the platform built (modulemanager.go:376-383).
	//
	// This is NOT the plugin's own Go module object. To reach that, get the elements.Module server
	// element and call GetObject/GetUserModule. Returns a NotFound error for an unknown module.
	GetModule(ctx core.ServerContext, modName string) (core.Module, error)
}
