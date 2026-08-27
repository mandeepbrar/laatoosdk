package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

// WorkflowManager is the server element application code uses to start and inspect workflows —
// from ctx.GetServerElement(core.ServerElementWorkflowManager).
//
// It is the ROUTER in front of the engines, not an engine: every components.WorkflowManager method
// it carries dispatches to whichever engine claimed the named workflow, walking up to the parent
// server level when the name is not registered here. Those methods' docs record where the engines
// diverge, and they diverge a lot — read them before relying on a call having an effect.
type WorkflowManager interface {
	core.ServerElement
	components.WorkflowManager

	// RegisterProvider claims a WorkflowType for an engine, so every definition declaring that type
	// is routed to it. Engines call this from their service factory as they are constructed;
	// application code never does.
	//
	// ONE ENGINE PER TYPE, AND THE SECOND ONE FAILS. Registering a type that is already taken
	// returns an error rather than replacing or chaining
	// (laatooserver/src/core/workflowmanager.go:77-84). This is load-bearing: cadence and
	// gcpworkflows both register `durable`, so a solution can enable only one of them, and
	// goworkflows joins that contest when its module sets mode: durable instead of the `process`
	// type it otherwise takes.
	//
	// Registering also adds the engine to the set LoadWorkflows walks at startup and that a
	// wfType-less Subscribe fans out to.
	RegisterProvider(ctx core.ServerContext, wfType components.WorkflowType, mgr components.WorkflowManager) error

	// GetWorkflow returns a registered workflow definition by name, and whether one was found.
	//
	// Until this existed there was no way to reach a definition at all: components.WorkflowManager
	// offers LoadWorkflows, StartWorkflow, GetWorkflowInstance, IsWorkflowRegistered, SendSignal,
	// CompleteActivity and Subscribe — none of which hand back a Workflow — and WorkflowInstance
	// carries no route back to the definition it was started from. Anything needing to ask a
	// property of a definition, retriability being the first, had nowhere to ask.
	//
	// IT IS ON THE ELEMENT AND DELIBERATELY NOT ON components.WorkflowManager. The server already
	// receives every Workflow from each engine when it loads them, so it can answer this from what
	// it holds. Putting it on the engine-facing interface instead would oblige all four engines to
	// implement a method that adds nothing, and would break each of them silently — Go checks
	// interface satisfaction at the assertion site, so they would compile and fail at startup.
	//
	// The (value, bool) shape and the ServerContext mirror IsWorkflowRegistered, which asks the
	// adjacent question about the same registry.
	GetWorkflow(ctx core.ServerContext, name string) (components.Workflow, bool)
}
