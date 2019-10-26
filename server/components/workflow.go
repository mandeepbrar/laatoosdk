package components

import "laatoo/sdk/server/core"

type WorkflowInitiator interface {
	StartWorkflow(ctx core.RequestContext, workflowName string, initVal interface{}) error
}
