package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/utils"
)

// Module is the interface a plugin's Go module object satisfies. Like Service, a plugin module
// struct embeds core.Module and the server reflects a field literally named "Module" and assigns
// its own implementation into it (laatooserver core/module.go:110-116); a struct that does not
// embed core.Module under that exact field name fails to load with Core_Error_Type_Mismatch.
//
// Lifecycle, in the order the server drives it:
//
//	Describe -> module config processed -> registry directories loaded from disk
//	  -> configuration values injected into struct fields -> Initialize
//	  -> the element accessors below are called ONCE -> Start
//
// The timing of that "called once" step is the thing to know. Factories, Services, Agents,
// Channels, Rules, Tasks, Topics, Datasets and Permissions are all read immediately after
// Initialize returns (laatooserver core/module.go:380-446) and never again, so a module that
// builds its elements in Start contributes nothing. Each accessor's result is merged OVER
// whatever the module's registry directory already declared -- except Permissions, which
// replaces rather than merges.
type Module interface {
	ConfigurableObject
	// Metadata returns the module's ModuleInfo. The server builds it from
	// src/server/registry/objects/<objectName>.yml and overwrites anything the object loader had
	// supplied, so that file is effectively mandatory for a module with a Go object
	// (laatooserver core/module.go:126-131).
	Metadata() ModuleInfo
	// MetaInfo is a free-form descriptor a module may return. NOTHING CALLS IT: no code in
	// laatooserver or laatoomodules invokes MetaInfo, although two platform modules (dataadapter
	// and entity) implement it. The server's own implementation returns a fresh empty map.
	MetaInfo(ctx ServerContext) utils.StringMap
	// Describe is where a module declares its configurations, before any configuration is parsed.
	// Its error is DISCARDED -- the server calls it as a bare statement
	// (laatooserver core/module.go:142) -- so returning an error here neither aborts the load nor
	// gets logged. Signal failure from Initialize instead, whose error is checked.
	Describe(ServerContext) error
	// Initialize runs after the module's registry directories have been loaded from disk and
	// after configuration-backed struct fields have been injected. conf is the module's raw
	// settings block from config/modules/<name>.yml. Everything the accessors below return must
	// exist by the time this returns.
	Initialize(ctx ServerContext, conf config.Config) error
	// Start runs after every element the module contributed has been created. It is too late to
	// contribute services, channels or any other element here.
	Start(ctx ServerContext) error
	// Factories returns service factory configurations, keyed by factory name, merged over
	// src/server/registry/factories/. Return nil to contribute none.
	Factories(ctx ServerContext) map[string]config.Config
	// Services returns service configurations, keyed by service alias, merged over
	// src/server/registry/services/. This is the programmatic equivalent of a service YAML.
	Services(ctx ServerContext) map[string]config.Config
	// Agents returns agent configurations, merged over src/server/registry/agents/.
	Agents(ctx ServerContext) map[string]config.Config
	// Rules returns rule configurations, merged over src/server/registry/rules/.
	Rules(ctx ServerContext) map[string]config.Config
	// Datasets returns dataset configurations, merged over src/server/registry/datasets/.
	Datasets(ctx ServerContext) map[string]config.Config
	// Permissions returns the module's permission declarations. Unlike every other accessor here
	// the result REPLACES rather than merges: a non-nil return overwrites whatever the registry
	// directory declared (laatooserver core/module.go:443-446). Return nil to leave the
	// file-declared permissions alone.
	Permissions(ctx ServerContext) utils.StringsMap
	// Channels returns channel configurations, keyed by channel name, merged over
	// src/server/registry/channels/.
	Channels(ctx ServerContext) map[string]config.Config
	// Tasks returns task-queue configurations, merged over src/server/registry/tasks/.
	Tasks(ctx ServerContext) map[string]config.Config
	// Topics returns the messaging topics this module declares, merged over
	// src/server/registry/topics/.
	//
	// A topic belongs to the code that publishes and subscribes to it, not to whichever solution
	// happens to host that code. Before modules could declare their own, a topic existed only if a
	// solution author had named it in configuration — and publishing to a topic nobody named
	// silently dropped the message, so a plugin's events depended on someone else remembering them.
	Topics(ctx ServerContext) map[string]config.Config
	// Workflows returns workflow configurations -- but the server NEVER CALLS IT. The only call
	// site is commented out (laatooserver core/module.go:438-442), so a module that returns
	// workflows here contributes nothing and gets no warning; dataadapter implements it and its
	// result is discarded. Ship workflows as .dsl files in the module's registry instead.
	Workflows(ctx ServerContext) map[string]config.Config
	// Activities returns activity configurations -- but the server NEVER CALLS IT. There is no
	// call site anywhere in laatooserver, and the server's own implementation has its single
	// statement commented out and unconditionally returns nil
	// (laatooserver core/moduleimpl.go:83-88). Activities are loaded from the module's registry
	// directory by the activity manager instead.
	Activities(ctx ServerContext) map[string]config.Config
	// GetContext returns the module's own ServerContext, the one its elements are created under.
	GetContext() ServerContext
	// ServerElement returns the ServerElementModule proxy wrapping this module; a plugin reaches
	// its own module object back through ServerElement().(elements.Module).GetUserModule().
	// The server's implementation dereferences the backing module pointer WITHOUT the nil check
	// every other method here uses (laatooserver core/moduleimpl.go:141-143), so it panics on a
	// module impl that was constructed without one.
	ServerElement() ServerElement
	//	GetContext(ctx ServerContext, variable string) (interface{}, bool)
}
