package components

import (
	"context"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type Workflow interface {
	Spec(ctx context.Context) interface{}
	Type() string
	GetName() string
	GetModule() core.Module
}

type WorkflowInstance interface {
	GetId() string
	InstanceDetails() utils.StringMap
	GetWorkflow() string
	GetStatus() utils.StringMap
	InitData() utils.StringMap
}

type WorkflowActivityType string

const (
	MANUAL    WorkflowActivityType = "manual"
	AUTOMATIC WorkflowActivityType = "automatic"
	DECISION  WorkflowActivityType = "decision"
)

type WorkflowManager interface {
	LoadWorkflows(ctx core.ServerContext, dir string, module core.Module) (map[string]Workflow, error)
	StartWorkflow(ctx core.RequestContext, workflowName string, initVal utils.StringMap, insconf utils.StringMap) (WorkflowInstance, error)
	IsWorkflowRegistered(ctx core.ServerContext, name string) bool
	SendSignal(ctx core.RequestContext, workflowId string, workflowIns string, actId string, signal string, signalVal utils.StringMap) error
	CompleteActivity(ctx core.RequestContext, workflowId string, workflowIns string, actId string, data utils.StringMap, err error) error
}
