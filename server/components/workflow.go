package components

import (
	"context"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type WorkflowStatus string

const (
	WorkflowStatusRunning    WorkflowStatus = "running"
	WorkflowStatusCompleted  WorkflowStatus = "completed"
	WorkflowStatusFailed     WorkflowStatus = "failed"
	WorkflowStatusCanceled   WorkflowStatus = "canceled"
	WorkflowStatusTerminated WorkflowStatus = "terminated"
	WorkflowStatusPending    WorkflowStatus = "pending"
)

type Workflow interface {
	Spec(ctx context.Context) interface{}
	Type() string
	GetName() string
	GetModule() core.Module
}

type WorkflowInstance interface {
	GetId() string
	GetExecutionId() string
	GetWorkflowName() string
	GetStatus() WorkflowStatus
	GetVariables() utils.StringMap
	InitData() utils.StringMap
	GetError() string
	GetPendingActivities() []string
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
	GetWorkflowInstance(ctx core.RequestContext, instanceId string) (WorkflowInstance, error)
	IsWorkflowRegistered(ctx core.ServerContext, name string) bool
	SendSignal(ctx core.RequestContext, workflowId string, workflowIns string, actId string, signal string, signalVal utils.StringMap) error
	CompleteActivity(ctx core.RequestContext, workflowId string, workflowIns string, actId string, data utils.StringMap, err error) error
}
