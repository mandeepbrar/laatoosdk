package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// TaskManager dispatches background work to queues and reports its completion.
//
// A queue must be DECLARED before use, in a plugin's registry/tasks/<task>.yml, which binds the
// queue name to a processor service and to the backend ("manager") that carries it. Pushing to an
// undeclared queue is a Core_Bad_Conf error, not a silent no-op, because the backend lookup goes
// through that declaration.
//
// BACKENDS DIFFER IN WHAT THEY CAN ANSWER, and the difference is not abstracted away:
//   - embedded NATS/JetStream: durable, and answers GetTask from a KV record it keeps alongside
//     the stream (the stream itself drops a message once acked).
//   - gcptasks: answers GetTask only for tasks Cloud Tasks still HOLDS -- pending, delayed or
//     retrying -- and only when the queue was configured to name its tasks. A dispatched task is
//     gone.
//   - beanstalktasks and tunnymemtasks: GetTask always returns NotImplemented; nothing is retained.
//
// Queues are synchronous by default. A queue declared `async: true` means the processor returning
// signals ACCEPTANCE, not completion, and its single completion event must be published by whoever
// finishes the work, through CompleteTask.
type TaskManager interface {
	core.ServerElement

	// PushTask enqueues task on queue and returns the invocation id assigned to it. The task's
	// payload is JSON-marshalled, and the caller's user and tenant are captured onto it so the
	// processor runs with the pusher's identity rather than as system.
	//
	// Returns Core_Bad_Conf when no declared backend owns queue. On any other failure the returned
	// id is still populated while err is non-nil -- check the error, not the id.
	//
	// WHETHER THE TASK IS RETAINED IS THE BACKEND'S CHOICE, not this call's; see the interface
	// comment. Only the NATS backend records a pushed task for later lookup, and even there the
	// record is best effort: a failure to record is logged and does not fail the push.
	PushTask(ctx core.RequestContext, queue string, task interface{}, metadata utils.StringMap) (string, error)

	// SubscribeTaskCompletion subscribes handler to the completion events of a queue.
	//
	// topic is a QUEUE NAME, not a full topic: a trailing ".>" or ".*" is stripped and the name is
	// mapped onto the platform's static completion topic for that queue. An empty topic (or a bare
	// ">"/"*") subscribes to every queue's completions. The invocation id is carried inside each
	// TaskCompletionMessage, so listeners dispatch on it rather than subscribing per invocation.
	//
	// Returns an internal error when no messaging manager is configured.
	SubscribeTaskCompletion(ctx core.ServerContext, topic string, handler core.MessageListener, subscriberId string) error

	// CompleteTask publishes the completion event for one invocation: its result, its metadata and
	// its error, if any. This is the call an ASYNC queue's worker makes when the work actually
	// finishes; for a synchronous queue the framework already published one when the processor
	// returned, and calling this again publishes a second event for the same id that no subscriber
	// can distinguish from the first.
	//
	// The completing user's id and roles are added to metadata (on a copy -- the caller's map is
	// not mutated) so downstream workflow steps run as that user.
	//
	// RETURNS nil WITHOUT PUBLISHING ANYTHING when no messaging manager is configured; the only
	// trace is a warning in the log. A caller that treats a nil error as "the completion was
	// delivered" is wrong on such a deployment.
	CompleteTask(ctx core.RequestContext, queue string, invocationId string, result interface{}, metadata utils.StringMap, err error) error

	// ProcessTask runs a task through its queue's registered processor service. Backends call this
	// on delivery; application code normally pushes instead.
	//
	// Returns Core_Bad_Conf when the queue has no processor. REFUSES a task whose identity could
	// not be rebuilt (TASK_IDENTITY_UNAVAILABLE) rather than running it as system, since a system
	// request has authorization disabled and would grant a broken task more access than a valid
	// one. Publishes the completion itself only for a synchronous queue.
	ProcessTask(ctx core.ServerContext, task *components.Task) (interface{}, error)

	// GetTask returns a previously pushed task by the id PushTask returned. The queue is
	// required because it selects which configured backend owns the task.
	//
	// Returns Core_Bad_Conf when no backend owns queue, and a NotImplemented error when the
	// backend serving it cannot address tasks by id -- always on beanstalktasks and tunnymemtasks,
	// and on gcptasks when task naming is disabled.
	//
	// RETURNS (nil, nil) ON THE NATS BACKEND when the task's record has aged out of the record
	// bucket: absent is not treated as an error there, so an `if err != nil` guard passes and the
	// caller must nil-check the task.
	GetTask(ctx core.RequestContext, queue string, id string) (*components.Task, error)

	// List returns each queue that has a running processor mapped to the module name of that
	// processor service. Queues declared with no processor do not appear.
	List(ctx core.ServerContext) utils.StringsMap

	// CreateEmptyTaskObj builds a Task pre-populated with an empty instance of the level's
	// configured principal type, so a codec can unmarshal the identity in place. Backends use it
	// when decoding a delivered task.
	//
	// The returned Task's User is nil when the configured principal type could not be created --
	// a misconfigured or unregistered user object. That is logged, not returned, and ProcessTask
	// later refuses such a task.
	CreateEmptyTaskObj(ctx core.ServerContext) *components.Task
}
