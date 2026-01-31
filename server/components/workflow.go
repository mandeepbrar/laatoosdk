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
	Subscribe(ctx core.RequestContext, eventType WorkflowEventType, handler core.MessageListener) error
}

type WorkflowActivityType string

const (
	MANUAL    WorkflowActivityType = "manual"
	AUTOMATIC WorkflowActivityType = "automatic"
	DECISION  WorkflowActivityType = "decision"
)

type WorkflowEventType string

// Workflow Lifecycle Events
const (
	EventWorkflowStarted   WorkflowEventType = "workflow.execution.started"
	EventWorkflowCompleted WorkflowEventType = "workflow.execution.completed"
	EventWorkflowFailed    WorkflowEventType = "workflow.execution.failed"
	EventWorkflowCanceled  WorkflowEventType = "workflow.execution.canceled"
)

// Activity lifecycle Events
const (
	EventActivityScheduled WorkflowEventType = "activity.task.scheduled"
	EventActivityStarted   WorkflowEventType = "activity.task.started"
	EventActivityCompleted WorkflowEventType = "activity.task.completed"
	EventActivityFailed    WorkflowEventType = "activity.task.failed"
	EventActivityRetrying  WorkflowEventType = "activity.task.retrying"
)

// WorkflowEvent represents an event specifically related to workflow or activity execution
type WorkflowEvent struct {
	core.Event
	WorkflowID string `json:"workflow_id"`
	InstanceID string `json:"instance_id"`
	ActivityID string `json:"activity_id,omitempty"`
}

type WorkflowManager interface {
	LoadWorkflows(ctx core.ServerContext, dir string, module core.Module) (map[string]Workflow, error)
	StartWorkflow(ctx core.RequestContext, workflowName string, initVal utils.StringMap, insconf utils.StringMap) (WorkflowInstance, error)
	GetWorkflowInstance(ctx core.RequestContext, instanceId string) (WorkflowInstance, error)
	IsWorkflowRegistered(ctx core.ServerContext, name string) bool
	SendSignal(ctx core.RequestContext, workflowId string, workflowIns string, actId string, signal string, signalVal utils.StringMap) error
	CompleteActivity(ctx core.RequestContext, workflowId string, workflowIns string, actId string, data utils.StringMap, err error) error
	//Subscribe to workflow events
	Subscribe(ctx core.RequestContext, eventType WorkflowEventType, handler core.MessageListener) error
}
