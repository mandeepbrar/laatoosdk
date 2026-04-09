package ai

import (
	"time"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// HITLManager is the server-level coordinator for Human-in-the-Loop steps.
// It is exposed via AgentManager.GetHITLManager() so any agent or skill can
// pause for user input without importing workflow-specific packages.
//
// There are two routing paths, selected by which fields are set in HITLTask:
//
//  Workflow path (WorkflowID + InstanceID + ActivityID set):
//    PauseForUser is NOT used — the DSL system.manual_task activity owns the pause.
//    CompleteTask routes to wfMgr.CompleteActivity to resume the paused workflow.
//
//  Skill / goal-agent path (AgentID + SessionID set, WorkflowID empty):
//    PauseForUser registers a task via TaskManager (distributed, pod-agnostic),
//    streams the question to the user, then blocks the calling goroutine until
//    the user replies or the task times out.
//    CompleteTask routes via TaskManager.CompleteTask, which pub-subs the result
//    back to whichever pod has the blocked goroutine — no sticky sessions needed.
type HITLManager interface {
	// PauseForUser streams question to the user, registers a HITL task via
	// TaskManager, and blocks the calling goroutine until CompleteTask is called
	// with the matching TaskID (from any pod) or the task times out.
	//
	// Use this inside a skill or goal-agent task (non-workflow path only).
	// task must have AgentID + SessionID set; WorkflowID must be empty.
	// The TaskID field in task is populated by PauseForUser from the queue's
	// invocation ID and must be sent to the frontend so it can be round-tripped
	// back when the user replies.
	PauseForUser(ctx core.RequestContext, task *HITLTask, question string) (userReply string, err error)

	// CompleteTask routes task completion.
	// Workflow path (WorkflowID set): calls wfMgr.CompleteActivity.
	// Skill/goal-agent path (WorkflowID empty): calls TaskManager.CompleteTask,
	// which unblocks the goroutine waiting in PauseForUser on the correct pod.
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
