package core

import (
	"laatoo.io/sdk/utils"
)

// An ACTION is a single-step primitive -- send a notification, start a workflow, save a record --
// that the server or a plugin registers in code with elements.ActionManager.RegisterAction. Every
// multi-step construct calls actions by name: a workflow step, a goal agent's activity step, a flow
// activity's step, or Go code through the action manager.
//
// An action is NOT a service and carries NO access control. Activities, workflows and agents are
// the service-equivalents a caller is authorized against; an action runs only inside one of them,
// under its request context. What an action touches keeps its own checks -- an action invoking a
// service goes through the service layer, and one reading data goes through the data layer's
// record rules. To expose an action to callers directly, wrap it in an activity.

// ActionParam describes one entry an action reads from its params map.
type ActionParam struct {
	// Name is the key the action reads from params.
	Name string `json:"name" yaml:"name"`
	// Type is the expected value's type, for the designer and for documentation: string, int,
	// float, bool, object, array or any. The server does not coerce values to it.
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Description says what the value means.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Required params must be present and non-nil; ActionManager.ExecuteAction refuses a call
	// missing one before the action runs.
	Required bool `json:"required,omitempty" yaml:"required,omitempty"`
}

// ActionDescriptor describes a registered action. The registered-elements report, the knowledge
// graph, the studio and the AI designer list actions by their descriptors, so the description and
// params are what an author sees when choosing one.
type ActionDescriptor struct {
	// Name is the registered name, unique among the actions and activities of one namespace.
	Name string `json:"name" yaml:"name"`
	// Description says what the action does, in a sentence.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Params lists what the action reads from its params map.
	Params []ActionParam `json:"params,omitempty" yaml:"params,omitempty"`
	// Result describes what the action returns. A workflow step maps from it with result:.
	Result string `json:"result,omitempty" yaml:"result,omitempty"`
	// Streaming marks an action that streams to the conversation that started its workflow, as a
	// streaming activity does: a workflow engine gives a streaming step the session's stream and
	// its streaming response handler. Without it, a streaming action's chunks reach no one.
	Streaming bool `json:"streaming,omitempty" yaml:"streaming,omitempty"`
}

// ActionFunc is an action's implementation. It receives the calling request's context (the
// workflow's, activity's or agent's) and the params the caller passed, and returns the result.
type ActionFunc func(ctx RequestContext, params utils.StringMap) (interface{}, error)
