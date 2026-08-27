package components

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ModInfo is the description of one loaded module handed to a ModuleManagerPlugin's Load.
//
// ModSettings is the module instance's settings block MERGED OVER the plugin's own defaults — but
// only when the module instance actually declared a settings: block. A plugin whose config.yml has
// settings: but whose module YAML has none gets nil here, not the plugin defaults
// (laatooserver/src/core/modulemanager_createmodules.go:152-159).
//
// Configurations is the module's declared configuration schema; ModProps is the merged contents of
// the plugin's properties/ directory.
type ModInfo struct {
	ModName            string
	PluginName         string
	PluginDir          string
	ParentModName      string
	Mod                core.Module
	UserObj            core.Module
	PluginConf         config.Config
	ModConf            config.Config
	ModSettings        config.Config
	Configurations     map[string]core.Configuration
	ModProps           utils.StringMap
	IsExtended         bool
	ExtendedPluginName string
	ExtendedPluginConf config.Config
	ExtendedPluginDir  string
	Hot                bool
}

// GetContext reads a context variable from the described module's own server context.
//
// Only variables declared under the plugin's config.yml params: block are there to read. The
// platform copies a param into the module context when the module instance supplies a value, or
// falls back to the param's declared default; an UNDECLARED settings key is never copied and is
// therefore invisible here and to {{var "..."}} templates
// (laatooserver/src/core/modulemanager_createmodules.go:164-181). The ctx argument is accepted and
// IGNORED — the lookup always uses info.Mod's context.
func (info *ModInfo) GetContext(ctx core.ServerContext, variable string) (interface{}, bool) {
	return info.Mod.GetContext().Get(variable)
}

// ModuleManagerPlugin is implemented by a plugin's MODULE OBJECT (the type declared in the
// plugin's module registry entry) when that plugin needs to see every other module the server
// loads. It is how the UI aggregators, masterdata collectors and similar cross-cutting plugins
// discover contributions from unrelated plugins.
//
// Registration is by type assertion, not configuration: the module manager asserts each module's
// user object to this interface and, if it satisfies it, records it
// (laatooserver/src/core/modulemanager_managerplugins.go:55-71). Implement the whole interface or
// none of it — a partial implementation simply is not detected, silently.
type ModuleManagerPlugin interface {
	// GetName returns a label for this plugin, used only in the module manager's logging.
	GetName() string

	// Load is called ONCE PER MODULE, for every module in the server, with that module's ModInfo.
	// This is where a plugin harvests whatever the other module contributes.
	//
	// It runs on the OTHER module's server context (a subcontext of it), not on this plugin's, so
	// context variables read here are the loaded module's — which is what makes ModInfo.GetContext
	// meaningful.
	//
	// An error aborts the sweep: the whole module-loading pass fails, and modules the iteration
	// had not yet reached are never offered (modulemanager_managerplugins.go:76-90).
	//
	// modInfo.Hot distinguishes a hot reload from initial startup.
	Load(ctx core.ServerContext, modInfo *ModInfo) error

	// Loaded is called once after Load has been called for every module in the pass, signalling
	// that the plugin has seen the complete set and may now build whatever it aggregates.
	//
	// It is invoked on THIS plugin's own module context, unlike Load
	// (modulemanager_managerplugins.go:93-103). Do work that needs the full picture here, not at
	// the end of Load — Load has no way to know which module is last.
	Loaded(ctx core.ServerContext) error

	// Unloaded is intended as the post-removal counterpart of Unloading.
	//
	// THE SERVER NEVER CALLS IT. As of 2026-08-27 the only call site in the entire repository is a
	// plugin's own unit test; the hot-reload path calls Unloading and nothing else
	// (laatooserver/src/core/modulemanager_hotload.go:179-192). Cleanup placed here does not run.
	// Put it in Unloading.
	Unloaded(ctx core.ServerContext, insName, modName string) error

	// Unloading is called on every registered ModuleManagerPlugin, for each module being unloaded,
	// BEFORE that module's elements are torn down — during a hot reload, and for a module being
	// removed (modulemanager_hotload.go:179-192).
	//
	// insName is the module instance name and modName is its plugin name. Note the argument order
	// relative to the names: insName first.
	//
	// An error aborts the unload sweep, so remaining plugins are not notified. This is the only
	// unload hook the server actually invokes — see Unloaded.
	Unloading(ctx core.ServerContext, insName, modName string) error
}
