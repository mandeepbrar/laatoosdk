package components

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// Script is one loaded script — a single callable unit produced by a script engine
// (tengo, yaegi) from a file under a module's src/server/registry/scripts/ directory.
//
// A Script is reached as an activity: the platform wraps each registered script in a
// ScriptActivity whose Invoke calls GetScriptManager().InvokeScript
// (laatooserver/src/core/activities.go:174-196). Nothing else in the platform consumes a Script.
type Script interface {
	// GetName returns the script's alias — the file's basename without its extension, which is
	// the name it was registered under and the name an activity definition refers to.
	GetName(ctx core.ServerContext) string

	// GetParams is intended to describe the script's declared parameters.
	//
	// BOTH SHIPPED ENGINES RETURN nil, and nothing in the platform calls this method. Tengo's
	// implementation has its real body commented out
	// (laatoomodules/scripts/dev/plugins/tengo/src/server/go/tengoscript.go:25-28); yaegi's is a
	// bare `return nil` (.../yaegi/src/server/go/yaegiscript.go:27-29). Treat a nil result as
	// "not described", not as "no parameters", and do not build validation on it.
	GetParams(ctx core.ServerContext) map[string]core.Param

	// GetModule returns the module the script was loaded from, or nil when the engine did not
	// record one. No platform code calls this method today.
	GetModule() core.Module

	// GetScriptManager returns the ENGINE that owns this script — the provider that compiled it,
	// not the ScriptManager server element. This is the live dispatch path: ScriptActivity.Invoke
	// calls it and refuses the activity with a Bad Conf error when it returns nil
	// (laatooserver/src/core/activities.go:186-189), so an engine that builds a Script without
	// setting this makes every invocation of that script fail at runtime rather than at load.
	GetScriptManager() ScriptManager
}

// ScriptManager is the contract a SCRIPT ENGINE implements (tengo, yaegi). Do not confuse it with
// elements.ScriptManager, the server element that arbitrates between engines — an engine registers
// itself with that element by extension via RegisterProvider, and registers each script it
// compiles via RegisterScript.
type ScriptManager interface {
	// Load compiles or indexes every script this engine recognises under dir.
	//
	// The central script manager calls this once per module registry/scripts directory, for each
	// distinct registered provider, and SWALLOWS the error: a provider that fails to load is
	// logged and the sweep continues, so a broken script directory does not fail startup
	// (laatooserver/src/core/scriptmanager.go:188-207).
	//
	// TIMING TRAP: that sweep runs inside the script manager element's Initialize, which the
	// server performs BEFORE the service manager initializes any service
	// (laatooserver/src/core/abstractserver_initialize.go:193-218). An engine is a service, so its
	// RegisterProvider call happens after the sweep — by which time the provider map was empty and
	// nothing was loaded. Registering a provider does NOT retroactively re-run Load. An engine
	// must therefore discover and register its own scripts from its own lifecycle (tengo does this
	// from Start, via registerTengoActivities), and must not rely on the central sweep calling
	// Load for it.
	Load(ctx core.ServerContext, dir string) error

	// InvokeScript runs one script.
	//
	// THE RESULT IS NOT RETURNED — it is set on the request context. Both engines call
	// ctx.SetResponse(core.SuccessResponse(result)) and return only an error
	// (tengoactivitymanager.go:146, yaegi.go:77), so the caller reads the value back with
	// ctx.GetResponse(). An implementation that returns nil without setting a response leaves the
	// caller with no result and no error.
	//
	// The engine receives the Script it produced and is expected to type-assert it back to its own
	// concrete type, rejecting a Script belonging to another engine with a BadArg error.
	//
	// ARGUMENT ORDER IS NOT PRESERVED BY THE TENGO ENGINE. args is a map, and tengo flattens it to
	// a positional argument list by ranging over the map (tengoactivitymanager.go:135-138) — Go
	// map iteration order is randomised, so a tengo script taking more than one parameter receives
	// them in a different order on different runs. The yaegi engine passes the map through as a
	// single argument and is unaffected. Both engines prepend the RequestContext as the script's
	// first argument.
	InvokeScript(ctx core.RequestContext, act Script, args utils.StringMap) error
}
