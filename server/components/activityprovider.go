package components

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ActivityProvider is the contract a plugin implements to serve an activity TYPE the server does not
// build in -- the flow type, served by inmemoryworkflows, is the shipped one. It registers itself
// with elements.ActivityManager.RegisterActivityProvider from a factory's Initialize, and the
// activity manager hands it every activity definition of that type.
//
// TIMING: activity definitions are loaded during the activity manager's Initialize, before any
// plugin service or factory initializes, so a provider is NOT offered definitions as they load.
// Each provider-backed activity binds to its provider when the activity's own service STARTS, and
// calls LoadActivity then; a provider registered later than that finds the activity already
// failed to start.
type ActivityProvider interface {
	// LoadActivity prepares one definition of the provider's type -- for a flow, parsing and
	// validating its Body. An error fails the activity's start, naming the activity, so a
	// malformed definition is refused at boot rather than at its first call.
	LoadActivity(ctx core.ServerContext, def *core.ActivityDefinition) error

	// ExecuteActivity runs a loaded definition with the caller's params, within the caller's
	// request, and returns its result. The activity manager checks def.Input before the call and
	// projects the result to def.Output after it, so a provider need not.
	ExecuteActivity(ctx core.RequestContext, def *core.ActivityDefinition, params utils.StringMap) (utils.StringMap, error)

	// UnloadActivity drops whatever LoadActivity prepared, when the activity's module unloads.
	UnloadActivity(ctx core.ServerContext, def *core.ActivityDefinition) error
}
