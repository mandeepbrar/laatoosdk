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

	// IsRetriable reports whether this process may be started again from its original input after
	// a failure.
	//
	// It is on this interface rather than left to callers reading GetDefinition() because the
	// question is engine-agnostic and the answer is not: a definition is whatever format its engine
	// loaded, so asking through GetDefinition() forces every caller to type-assert into one
	// engine's types and to stop working the moment a second format exists.
	//
	// FALSE IS THE MEANINGFUL DEFAULT. A restart re-runs every step, so it is safe only where each
	// one is idempotent — a property only the definition's author knows, and the processes that
	// lack it are exactly the ones where a second run charges, provisions or publishes twice. An
	// implementation that cannot determine the answer must return false, never true: an omitted
	// declaration means the author did not say, and that is not permission.
	//
	// This does NOT describe per-activity retry. Statement.Retry and Statement.Timeout are declared
	// on the DSL and executed by no engine; this reports whether the WHOLE process may be started
	// again, which is a different capability and the one that is actually available.
	IsRetriable() bool
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
