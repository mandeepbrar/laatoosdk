package ai

import (
	"time"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// HITLManager is the server-level coordinator for Human-in-the-Loop workflow steps.
// It is exposed via AgentManager.GetHITLManager() so any agent or plugin can complete
// human tasks without importing workflow-specific packages.
//
// The manager is stateless: all workflow context is carried in the HITLTask struct,
// which is round-tripped by the frontend via the AG-UI STATE_SNAPSHOT hitlTask field.
// No server-side task state is stored here.
type HITLManager interface {
	// CompleteTask routes task completion through TaskManager.
	// The task struct carries the workflow context (workflowID, instanceID, activityID)
	// that was round-tripped from the frontend — no server-side task state is stored.
	CompleteTask(ctx core.RequestContext, task *HITLTask, result utils.StringMap) error

	// FailTask signals a task failure through TaskManager.
	FailTask(ctx core.RequestContext, task *HITLTask, reason string) error
}

// HITLTaskStatus mirrors workflow status for HITL tasks.
type HITLTaskStatus string

const (
	HITLTaskStatusPending   HITLTaskStatus = "pending"
	HITLTaskStatusCompleted HITLTaskStatus = "completed"
	HITLTaskStatusFailed    HITLTaskStatus = "failed"
)

// HITLTask holds context for a paused human-interaction workflow step.
// The frontend round-trips this as the hitlTask STATE_SNAPSHOT field.
// Agents construct it from the incoming request and pass it to HITLManager.
type HITLTask struct {
	TaskID     string
	WorkflowID string
	InstanceID string
	ActivityID string
	SessionID  string
	Config     *core.HumanTaskConfig
	CreatedAt  time.Time
	Status     HITLTaskStatus
}
