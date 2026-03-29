package ai

import (
	"time"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// HITLManager is the server-level coordinator for Human-in-the-Loop steps.
// It is exposed via AgentManager.GetHITLManager() so any agent or skill can
// complete human tasks without importing workflow-specific packages.
//
// The manager is stateless: all context is carried in the HITLTask struct,
// which is round-tripped by the frontend via the AG-UI STATE_SNAPSHOT hitlTask field.
// No server-side task state is stored here.
//
// The pause/issue side is handled by calling the system.manual_task activity,
// which works for both workflow agents (via DSL) and non-workflow agents (via
// direct activity invocation). CompleteTask and FailTask are the resume side,
// routing based on what fields are set in the round-tripped HITLTask:
//   - WorkflowID + InstanceID + ActivityID set → route to workflow TaskManager
//   - Only SessionID + TaskID set → route via session channel (non-workflow agents)
type HITLManager interface {
	// CompleteTask routes task completion.
	// For workflow tasks (WorkflowID set): routes through the workflow TaskManager
	// to resume the paused workflow instance.
	// For non-workflow tasks (SessionID + TaskID only): routes via session channel
	// to unblock the waiting goal agent or pro-code skill.
	CompleteTask(ctx core.RequestContext, task *HITLTask, result utils.StringMap) error

	// FailTask signals a task failure through the appropriate routing path.
	FailTask(ctx core.RequestContext, task *HITLTask, reason string) error
}

// HITLTaskStatus mirrors workflow status for HITL tasks.
type HITLTaskStatus string

const (
	HITLTaskStatusPending   HITLTaskStatus = "pending"
	HITLTaskStatusCompleted HITLTaskStatus = "completed"
	HITLTaskStatusFailed    HITLTaskStatus = "failed"
)

// HITLTask holds context for a paused human-interaction step.
// The frontend round-trips this as the hitlTask STATE_SNAPSHOT field.
//
// Workflow fields (WorkflowID, InstanceID, ActivityID) are optional.
// Set them for workflow-based agents; leave empty for goal agents and pro-code skills.
// AgentID is used for non-workflow routing and should be set when workflow fields are absent.
type HITLTask struct {
	TaskID     string
	SessionID  string
	AgentID    string // set for non-workflow agents (goal agents, skills); empty for workflow tasks
	WorkflowID string // optional — workflow agents only
	InstanceID string // optional — workflow agents only
	ActivityID string // optional — workflow agents only
	Config     *core.HumanTaskConfig
	CreatedAt  time.Time
	Status     HITLTaskStatus
}
