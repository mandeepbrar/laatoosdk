package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ActionManager is the registry of ACTIONS at one server level: the single-step primitives the
// server and plugins register in code (core.ActionFunc plus a core.ActionDescriptor). It is the
// element at core.ServerElementActionsManager.
//
// AN ACTION IS NOT A SERVICE AND CARRIES NO ACCESS CONTROL. Workflows, activities and agents are
// the service-equivalents a caller is authorized against; an action runs inside one of them, under
// its request context, and nothing exposes an action on a channel, to a UI action or to a job. To
// expose one, wrap it in an activity.
//
// Actions resolve NEAREST NAMESPACE FIRST: an action registered at one level is callable at every
// level below it, and a nearer level may register its own of the same name, which wins there.
//
// ACTIONS AND ACTIVITIES RESOLVE SEPARATELY. A workflow step says which it runs -- `action:` reaches
// ExecuteAction, `activity:` reaches ActivityManager.ExecuteActivity -- and neither falls back to
// the other, so an action and an activity may share a name. The built-in action "ExecuteActivity"
// is how a caller holding only actions runs an activity.
type ActionManager interface {
	core.ServerElement

	// RegisterAction registers fn under desc.Name.
	//
	// CALL IT FROM A MODULE'S OR FACTORY'S Initialize, so the action exists before any workflow,
	// activity or agent starts.
	//
	// A NAME ALREADY TAKEN BY ANOTHER ACTION AT THIS LEVEL IS REFUSED, NOT REPLACED, and the error
	// names both registrants. A re-registration by the same registrant (a hot reload) replaces its
	// own action. An activity of the same name is not a conflict.
	RegisterAction(ctx core.ServerContext, desc core.ActionDescriptor, fn core.ActionFunc) error

	// ExecuteAction runs the named action with params, under ctx, and returns its result.
	//
	// A required param (core.ActionParam.Required) that is absent or nil is refused with a missing
	// argument error before the action runs. Nothing else is checked: there is no authorization.
	// An unknown name answers a not-found error; an ACTIVITY of the name is never run in its place.
	ExecuteAction(ctx core.RequestContext, name string, params utils.StringMap) (interface{}, error)

	// GetAction returns the descriptor of the named action, resolved nearest-first, or nil when no
	// action has that name.
	GetAction(ctx core.ServerContext, name string) *core.ActionDescriptor

	// ListActions returns the descriptors of every action reachable from this level -- its own and
	// every enclosing level's -- each name once, the nearest registration winning.
	ListActions(ctx core.ServerContext) []core.ActionDescriptor
}
