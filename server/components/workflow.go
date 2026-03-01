package components

import (
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
	GetDefinition() interface{}
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
	// Aligned with workflow DSL `activity_type` and core.ActivityType
	MANUAL   WorkflowActivityType = "manual"
	SERVICE  WorkflowActivityType = "service"
	SCRIPT   WorkflowActivityType = "script"
	EXECUTOR WorkflowActivityType = "executor"

	// Legacy values (deprecated; prefer SERVICE/SCRIPT/EXECUTOR + switch/decision statements)
//	AUTOMATIC WorkflowActivityType = "automatic"
//	DECISION  WorkflowActivityType = "decision"
)

// IsAutomatic returns true for any activity type that should be executed
// automatically by the workflow engine (i.e. not a human/manual step).
// This covers the canonical types service, script, executor, the legacy
// "automatic" value, and an empty string (the default when no type is set).
func (t WorkflowActivityType) IsAutomatic() bool {
	switch t {
	case SERVICE, SCRIPT, EXECUTOR, "":
		return true
	}
	return false
}

// IsManual returns true only when the activity type explicitly requires
// human interaction.
func (t WorkflowActivityType) IsManual() bool {
	return t == MANUAL
}

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

type WorkflowType string

const (
	WorkflowTypeFunction WorkflowType = "function"
	WorkflowTypeProcess  WorkflowType = "process"
	WorkflowTypeDurable  WorkflowType = "durable"
	WorkflowTypeAgent    WorkflowType = "agent"
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
	GetWorkflowInstance(ctx core.RequestContext, workflowId string, instanceId string) (WorkflowInstance, error)
	IsWorkflowRegistered(ctx core.ServerContext, name string) bool
	SendSignal(ctx core.RequestContext, workflowId string, workflowIns string, actId string, signal string, signalVal utils.StringMap) error
	CompleteActivity(ctx core.RequestContext, workflowId string, workflowIns string, actId string, data utils.StringMap, err error) error
	//Subscribe to workflow events
	Subscribe(ctx core.RequestContext, wfType WorkflowType, eventType WorkflowEventType, handler core.MessageListener) error
}
