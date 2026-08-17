package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

type WorkflowManager interface {
	core.ServerElement
	components.WorkflowManager
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
