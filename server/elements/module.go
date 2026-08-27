package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// Module is the server element representing ONE CONFIGURED INSTANCE of a plugin — what
// config/modules/<name>.yml declares. One plugin can back many modules with different settings.
//
// Reach it with ctx.GetServerElement(core.ServerElementModule); GetName() returns the module
// instance name, not the plugin name.
type Module interface {
	core.ServerElement

	// GetObject returns the plugin's own Go module object — the type named in the plugin's module
	// registry entry, the same value a ModuleManagerPlugin sees as ModInfo.UserObj.
	//
	// nil for a module whose plugin ships no Go module object (a pure YAML/UI plugin), which is
	// the common case. Callers must nil-check.
	//
	// IDENTICAL TO GetUserModule — both return the same field
	// (laatooserver/src/core/moduleproxy.go:17-19 and :32-34). Neither is the platform-side
	// core.Module that ModuleManager.GetModule hands back; that one is the module's own
	// implementation object, not the plugin's.
	GetObject() core.Module

	// GetModuleProperties returns the merged contents of the plugin's properties/ directory —
	// localisation strings and similar static per-module data, deep-merged across the plugin and
	// any plugin it extends (laatooserver/src/core/module.go:218-228). It is NOT the module's
	// settings and NOT its context variables.
	GetModuleProperties() utils.StringMap

	// GetUserModule returns the plugin's own Go module object. See GetObject — this returns
	// exactly the same value.
	GetUserModule() core.Module

	// Metadata returns the module's declared configuration schema, built from the plugin's
	// registry/objects/<ModuleObject>.yml. This is what makes a module's settings readable through
	// the typed GetConfiguration accessors; a setting absent from it is not declared and cannot be
	// read that way.
	Metadata() core.ModuleInfo
}
