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
	// Pause records a human wait and returns its opaque handle. It does not block, and it
	// does not touch anyone's stream.
	//
	// The manager builds the record. Callers describe the pause and never construct a
	// HITLTask themselves — an object handed in half-filled, with some fields expected and
	// others expected to be left alone, is a contract every caller has to remember and any
	// caller can get wrong.
	//
	// kind selects the resume handler and is a typed constant rather than a map key so a
	// wrong one fails at compile time instead of in a guard on some other pod. resume is
	// whatever that handler needs to wake this caller, opaque here and never sent onward.
	// data is what the pause is about — the question, the fields outstanding — and is the
	// half a client may see. They are separate parameters because collapsing them into one
	// map would make leaking resume internals the default rather than the mistake.
	//
	// Telling the user is the caller's job, because only the caller knows which stream is
	// theirs: a skill is inside the request it needs to close, while something observing a
	// parked execution from an event has to close the session's stream instead. A manager
	// that closed "the" stream would be right for one of them and wrong for the other.
	//
	// The caller is expected to yield after this returns — a skill by returning a waiting
	// result with the state it needs on re-entry, and a caller that suspends its own
	// execution by parking in whatever way its runtime provides.
	Pause(ctx core.RequestContext, sessionID string, kind HITLPauseKind, resume utils.StringMap, data utils.StringMap) (handle string, err error)

	// Complete supplies the human's answer for the pause named by handle, within sessionID.
	//
	// The session is a parameter because a pause is recorded in its own conversation's
	// memory and there is nowhere else to look it up; a completer always has it, since it
	// is answering inside that conversation. Passing it is safe in a way the identifiers
	// this interface used to carry were not — a client naming its own session grants itself
	// nothing, whereas a client naming a workflow instance named something it did not own.
	//
	// Everything else about the pause is resolved server-side from the record, so no caller
	// hands back the data its own resume will act on.
	//
	// The resume runs inline on this request. Nothing is broadcast: a broadcast is what a
	// pod-pinned waiter needs, and there are none.
	Complete(ctx core.RequestContext, sessionID string, handle string, result utils.StringMap) error

	// Fail abandons the pause named by handle within sessionID, reporting reason to the
	// waiting caller through the same resume handler. Used when a pause cannot be answered
	// rather than when the answer is negative — a rejection is an answer and goes through
	// Complete.
	Fail(ctx core.RequestContext, sessionID string, handle string, reason string) error

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

// HITLResumeSkillKey is the resume-map key naming the skill to re-enter on a skill pause.
//
// A constant rather than a field on HITLTask, because it is not a property of every pause
// — a parked execution has no skill — and because the resume map is already defined as
// what a handler needs to wake its caller. A field would put one kind's vocabulary into
// the shape all kinds share.
const HITLResumeSkillKey = "_skill"

// HITLPauseEnvelope is what a client is handed so it can answer a pause later.
//
// One function because there is one wire shape, and it was being built in five places: by
// each skill that returns a question, and by the server when a parked step is reported.
// They had already drifted -- one carried a "waiting" flag the others did not -- and the
// next field added would have gone into whichever site the author happened to be editing.
//
// It carries the handle and the session that handle belongs to, and nothing else. What the
// pause is for, and how it resumes, stay on the server.
func HITLPauseEnvelope(sessionID, handle string) utils.StringMap {
	return utils.StringMap{
		"taskId":    handle,
		"sessionId": sessionID,
		// Redundant with the envelope's presence, and kept because clients already
		// branch on it: removing it is a client change, not a server one.
		"waiting": true,
	}
}

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

// HITLTask is one recorded agentic pause, as the manager stored it.
//
// Callers do not build these — Pause describes a pause with plain arguments and the
// manager constructs the record. What a resume handler receives is this, read back at
// completion time.
//
// The client round-trips only TaskID. Everything else is resolved server-side, which is
// why this struct carries no workflow, instance, or activity identifier: those were fields
// a client supplied and a resume trusted. Anything a particular kind of resume needs now
// travels in Resume, opaque to everything above it.
type HITLTask struct {
	// TaskID is the opaque handle — the only field that reaches the client.
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
	//
	// It is not for anything a client should see. Keeping it separate from Data is what
	// stops resume internals travelling back out to the browser, which is the failure this
	// interface was reshaped to remove.
	Resume utils.StringMap

	// Data carries what the pause is about rather than how it resumes — the question put
	// to the person, the fields still outstanding, whatever a caller wants back when it is
	// re-entered or a client needs to re-render after a reconnect.
	//
	// Opaque to the manager, like Resume, and deliberately a map: what a pause needs to say
	// about itself differs per caller and will keep differing, and widening this interface
	// each time is not a change one repository can make.
	Data utils.StringMap

	// Config is the optional human-task description a caller attached to the pause.
	Config *core.HumanTaskConfig

	CreatedAt time.Time
	Status    HITLTaskStatus
}
