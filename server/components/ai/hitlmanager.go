package ai

import (
	"time"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// HITLManager is the server-level owner of every agentic Human-in-the-Loop pause —
// a skill, a goal agent, or a workflow agent waiting on a person. It is exposed via
// AgentManager.GetHITLManager() so any caller can pause without importing
// workflow-specific packages, and without knowing how any other kind of pause resumes.
//
// The manager records a pause and returns. It never blocks the calling goroutine:
// a goroutine parked on a channel holds the caller's continuation — its locals and its
// position in the loop — which pins the pause to one pod and to one process lifetime.
// A caller that has somewhere to keep its own continuation does not need one.
//
// What the manager knows is deliberately small: a pause has a kind, a session, and an
// opaque Resume map it never interprets. Waking the caller belongs to a resume handler
// registered for that kind. This is what keeps engine and workflow vocabulary out of
// this interface — adding a second kind of pause is a registration, not a new branch here.
type HITLManager interface {
	// Pause records a human wait and returns its opaque handle. It does not block.
	//
	// The caller is expected to yield after this returns: a skill returns a "waiting"
	// result carrying whatever state it needs on re-entry, and a caller that suspends its
	// own execution parks in whatever way its runtime provides. The question is streamed
	// to the user as the turn closes, with the handle, so the client can hand the handle
	// back when the person answers.
	//
	// task.TaskID is populated with the handle and must not be set by the caller.
	Pause(ctx core.RequestContext, task *HITLTask, question string) (handle string, err error)

	// Complete supplies the human's answer for the pause named by handle.
	//
	// The manager resolves the handle to the recorded pause, dispatches to the resume
	// handler registered for that pause's kind, and clears the record. Resolution happens
	// server-side precisely so no caller has to hand back the resume data — a client that
	// could supply it could also forge it.
	//
	// The resume runs inline on this request. Nothing is broadcast: a broadcast is what a
	// pod-pinned waiter needs, and there are none.
	Complete(ctx core.RequestContext, handle string, result utils.StringMap) error

	// Fail abandons the pause named by handle, reporting reason to the waiting caller
	// through the same resume handler. Used when a pause cannot be answered rather than
	// when the answer is negative — a rejection is an answer and goes through Complete.
	Fail(ctx core.RequestContext, handle string, reason string) error

	// RegisterResumeHandler binds a resume strategy to a pause kind, at startup.
	//
	// This is the seam that keeps the manager free of any caller's vocabulary: the
	// component that knows how to wake a given kind of caller registers that knowledge
	// here, rather than the manager acquiring a branch for each one. Registering a kind
	// twice replaces the previous handler.
	RegisterResumeHandler(ctx core.ServerContext, kind HITLPauseKind, handler HITLResumeHandler) error
}

// HITLResumeHandler wakes the caller of one kind of pause and delivers the human's answer.
//
// It receives the recorded pause — including its opaque Resume map, which is meaningful to
// the handler and to nothing else — and the result the completer supplied. An error means
// the caller was not resumed; the manager reports it and leaves the record in place, since a
// pause whose resume failed is still pending rather than finished.
type HITLResumeHandler func(ctx core.RequestContext, task *HITLTask, result utils.StringMap) error

// HITLPauseKind names how a paused caller is woken, and therefore which registered
// resume handler owns it. It is not a description of who paused — two different agents
// that both park their own execution share one kind.
type HITLPauseKind string

const (
	// HITLPauseSkill is a caller that returned rather than blocking, and is resumed by
	// being invoked again with its stored state plus the human's answer.
	HITLPauseSkill HITLPauseKind = "skill"

	// HITLPauseParked is a caller that suspended its own execution — a workflow agent
	// whose engine holds the parked step — and is resumed by signalling that execution.
	// The details of the signal live in the pause's Resume map and in the handler.
	HITLPauseParked HITLPauseKind = "parked"
)

// HITLTaskStatus is the lifecycle of a recorded pause.
type HITLTaskStatus string

const (
	HITLTaskStatusPending   HITLTaskStatus = "pending"
	HITLTaskStatusCompleted HITLTaskStatus = "completed"
	HITLTaskStatusFailed    HITLTaskStatus = "failed"
)

// HITLTask is one recorded agentic pause.
//
// The client round-trips only TaskID. Everything else is resolved server-side from the
// recorded pause, which is why this struct carries no workflow, instance, or activity
// identifier: those were fields a client supplied and a resume trusted. Anything a
// particular kind of resume needs now travels in Resume, opaque to everything above it.
type HITLTask struct {
	// TaskID is the opaque handle. Set by Pause; the only field that reaches the client.
	TaskID string

	// Kind selects the registered resume handler.
	Kind HITLPauseKind

	// SessionID is the conversation the pause belongs to, and scopes where it is recorded.
	// A pause outliving its session has nothing left to resume into.
	SessionID string

	// AgentID identifies the paused caller for a skill-kind resume.
	AgentID string

	// Resume carries whatever the registered handler needs to wake this caller. The
	// manager stores and returns it without interpretation; keys prefixed with "_" are
	// the convention for values private to one runtime.
	Resume utils.StringMap

	// Config is the optional human-task description a caller attached to the pause.
	Config *core.HumanTaskConfig

	CreatedAt time.Time
	Status    HITLTaskStatus
}
