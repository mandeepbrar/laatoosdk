package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

// ScriptManager is the server element that ARBITRATES BETWEEN SCRIPT ENGINES. It compiles nothing
// itself: engines (tengo, yaegi) claim file extensions with RegisterProvider and hand back the
// scripts they compiled with RegisterScript.
//
// Do not confuse it with components.ScriptManager, which is the engine-side interface.
type ScriptManager interface {
	core.ServerElement

	// RegisterScript adds one compiled script to the level's script registry under alias, which is
	// the name a workflow activity definition refers to and the name ScriptActivity.Start resolves
	// (laatooserver/src/core/activities.go:165-170).
	//
	// SILENTLY OVERWRITES an existing alias — a plain map assignment with no duplicate check
	// (laatooserver/src/core/scriptmanager.go:209-212). Two engines that both produce a script
	// named "notify" leave only the one registered later, with no error. Always returns nil.
	//
	// The registered Script must return a non-nil GetScriptManager(), since that is how it is
	// dispatched.
	RegisterScript(ctx core.ServerContext, alias string, act components.Script) error

	// RegisterProvider claims a file extension for an engine. A leading dot is added if absent, so
	// "t" and ".t" are the same claim.
	//
	// A DUPLICATE EXTENSION IS REFUSED, not replaced — the second engine to claim an extension
	// gets an error (scriptmanager.go:144-155).
	//
	// TIMING TRAP: registering a provider does NOT cause the platform to go and load scripts for
	// it. The one directory sweep that calls provider.Load runs inside this element's Initialize,
	// which the server performs BEFORE the service manager initializes any service
	// (laatooserver/src/core/abstractserver_initialize.go:193-218). An engine is a service, so by
	// the time it can call RegisterProvider the sweep has already run against an empty provider
	// map and loaded nothing. An engine must discover and register its own scripts from its own
	// lifecycle; the extension claim is arbitration only.
	RegisterProvider(ctx core.ServerContext, extension string, provider components.ScriptManager) error
}
