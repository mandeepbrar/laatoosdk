package ai

import (
	"time"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// HITLManager is the server-level coordinator for Human-in-the-Loop workflow steps.
// It is exposed via AgentManager.GetHITLManager() so any agent or plugin can create
// and complete human tasks without importing the workflowagents plugin.
type HITLManager interface {
	// CreateTask registers a new human task for a paused workflow activity.
	// Returns a taskID that uniquely identifies this pause point.
	CreateTask(ctx core.RequestContext, workflowID, instanceID, activityID, sessionID string,
		config *core.HumanTaskConfig) (taskID string, err error)

	// CompleteTask marks the task done and resumes the workflow via CompleteActivity.
	CompleteTask(ctx core.RequestContext, taskID string, result utils.StringMap) error

	// FailTask marks the task failed and resumes the workflow with an error signal.
	FailTask(ctx core.RequestContext, taskID string, reason string) error

	// GetTask returns a pending task by ID.
	GetTask(taskID string) (*HITLTask, bool)

	// GetPendingTasksForSession returns all pending tasks for a session.
	GetPendingTasksForSession(sessionID string) []*HITLTask
}

// HITLTaskStatus mirrors workflow status for HITL tasks.
type HITLTaskStatus string

const (
	HITLTaskStatusPending   HITLTaskStatus = "pending"
	HITLTaskStatusCompleted HITLTaskStatus = "completed"
	HITLTaskStatusFailed    HITLTaskStatus = "failed"
)

// HITLTask holds the state for a paused human-interaction workflow step.
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
