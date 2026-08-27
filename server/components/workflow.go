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

// Workflow is a loaded workflow definition. Reach one by name with
// elements.WorkflowManager.GetWorkflow.
//
// THERE IS EXACTLY ONE IMPLEMENTATION ACROSS ALL FOUR ENGINES — dsl.DSLWorkflow in the workflowcore
// plugin (laatoomodules/workflow/dev/plugins/workflowcore/src/sdk/go/dsl/dslworkflow.go).
// goworkflows, cadence, inmemoryworkflows and gcpworkflows all parse the same .dsl format and hand
// that type back; they differ in how they EXECUTE a definition, never in how they represent one.
type Workflow interface {
	// GetDefinition returns the parsed definition — a *dsl.DSLWorkflowConfig carrying the .dsl
	// file's metadata, spec and statement tree.
	//
	// The concrete type lives in a plugin, not in this SDK, so a caller in another module cannot
	// name it without importing that plugin. THAT IS WHY QUESTIONS ABOUT A DEFINITION BELONG ON
	// THIS INTERFACE — IsRetriable is the first — rather than being answered by type-asserting what
	// comes back here. Prefer adding a method over asserting.
	//
	// It is the live definition and not a copy: the accessor backfills a missing id, apiVersion and
	// kind into the config it returns (dslworkflow.go:38-50), and mutating what comes back mutates
	// what every subsequent start uses.
	GetDefinition() interface{}

	// Type returns THE CONSTANT STRING "workflow" (dslworkflow.go:56-58). It is not the workflow's
	// WorkflowType and carries no information.
	//
	// The declared function/process/durable/agent type — the one that decides which engine claims a
	// definition at load — lives in the .dsl metadata, reachable through
	// GetDefinition().Metadata.WorkflowType. Nothing routes on this method. The server logs its
	// value as the workflow's "type" when registering one
	// (laatooserver/src/core/workflowmanager.go:222-224), which is why that line always reads
	// "workflow" whatever the definition declares.
	Type() string

	// GetName returns the registered workflow name, which is the .dsl FILENAME without its
	// extension — not metadata.id, which GetDefinition backfills FROM this name when it is empty.
	// It is the name StartWorkflow, IsWorkflowRegistered, GetWorkflow and sub-workflow references
	// all take.
	GetName() string

	// GetModule returns the module the definition was loaded from — whatever LoadWorkflows was
	// passed for that directory.
	//
	// Its context is what the engines use to create objects, resolve response handlers and build
	// the system requests activities run under, so this is the definition's only link back to
	// plugin-scoped configuration.
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

// WorkflowInstance is a handle on one execution of a workflow.
//
// TREAT ONE AS A SNAPSHOT, NOT A LIVE VIEW. On every engine but Cloud Workflows the instance is a
// dsl.WorkflowInstance struct whose accessors return plain fields (dslworkflowinstance.go:46-75),
// and nothing ever updates the copy a caller is holding: the instance StartWorkflow returned
// reports the status it had at that moment, for ever. To see current state, call
// WorkflowManager.GetWorkflowInstance again. The gcpworkflows instance wraps a Google execution
// object and is equally a snapshot of the read that produced it.
//
// A re-read handle is also not a superset of the started one — GetWorkflowInstance rebuilds it from
// engine state and leaves several fields empty. See InitData, GetVariables and GetError.
type WorkflowInstance interface {
	// GetId returns the instance id — the value to pass back as instanceId to GetWorkflowInstance,
	// SendSignal and CompleteActivity.
	//
	// ITS SHAPE IS THE ENGINE'S AND MUST NEVER BE PARSED: goworkflows builds
	// "<workflowName>_<uuid>" (goworkflows workflowservice.go:415), inmemoryworkflows a bare uuid,
	// and gcpworkflows the full Cloud Workflows execution RESOURCE NAME, which the Executions API
	// only accepts back unshortened.
	//
	// cadence puts the RunId here and the WorkflowId in GetExecutionId
	// (cadence workflowservice.go:344-348) — the reverse of what the two names suggest.
	GetId() string

	// GetExecutionId returns the engine's second identifier for the run, where it has one.
	//
	// It is EMPTY on inmemoryworkflows, equal to GetId on gcpworkflows (Cloud Workflows has only
	// one identifier per execution), the go-workflows execution id on goworkflows, and — see GetId
	// — cadence's WorkflowId rather than an execution id.
	//
	// goworkflows NEEDS it to look an instance up again and remembers it only in a process-local
	// map (goworkflows workflowservice.go:487-497), so GetWorkflowInstance there cannot resolve a
	// run started by another replica or before a restart.
	GetExecutionId() string

	// GetWorkflowName returns the name of the workflow this instance runs — the same name
	// StartWorkflow was given, and the name GetWorkflowInstance's workflowId argument takes.
	GetWorkflowName() string

	// GetStatus returns the instance's status as of the read that produced this handle.
	//
	// goworkflows CANNOT REPORT FAILURE THROUGH THIS. Its status mapping returns Running for
	// anything that is not Active/ContinuedAsNew/Finished, so a failed, cancelled or terminated run
	// comes back as WorkflowStatusRunning (goworkflows workflowservice.go:512-521) and polling for a
	// terminal failure never sees one. cadence and gcpworkflows map the full set.
	//
	// The instance StartWorkflow returns is stamped Running — except on inmemoryworkflows, which
	// runs the whole workflow synchronously before returning and so hands back a terminal status —
	// and never changes afterwards.
	GetStatus() WorkflowStatus

	// GetVariables returns the workflow's variable set.
	//
	// USUALLY EMPTY, and empty is not evidence the workflow has no variables. On the dsl-backed
	// engines the field is populated on the instance the ENGINE executes, not on the handle
	// StartWorkflow or GetWorkflowInstance hands back — goworkflows rebuilds the handle without it
	// (goworkflows workflowservice.go:499-517), so this returns nil there. On gcpworkflows it is the
	// execution's RESULT, available only once the execution has finished and empty while it runs
	// (gcpworkflows instance.go:74-91).
	//
	// Where it is populated it is not a defensive copy: mutating the returned map mutates the
	// instance's own variables.
	GetVariables() utils.StringMap

	// InitData returns the initVal map the workflow was started with.
	//
	// Populated on the instance StartWorkflow returns and NOT repopulated by GetWorkflowInstance,
	// which rebuilds the handle from engine state and leaves it nil on goworkflows, cadence and
	// gcpworkflows. Nil on a re-read instance is the normal answer, not evidence that the workflow
	// was started with nothing.
	InitData() utils.StringMap

	// GetError returns the failure message, empty when the instance has not failed.
	//
	// ALWAYS EMPTY ON goworkflows AND cadence: neither populates the field when rebuilding an
	// instance, and goworkflows' extraction is commented out (goworkflows workflowservice.go:
	// 509-514). An empty string is therefore NOT evidence of success — and on goworkflows GetStatus
	// cannot report the failure either, so a failed run there is indistinguishable from a healthy
	// one through this interface. gcpworkflows returns Google's error payload.
	GetError() string

	// GetPendingActivities returns the ids of activities the instance is waiting on — the human
	// steps a task list would show.
	//
	// EMPTY ON MOST DEPLOYMENTS, SILENTLY. cadence reports them from the execution description.
	// goworkflows tracks them in REDIS and returns an empty slice whenever no redis client is
	// configured (goworkflows workflowservice.go:539-548), so an unconfigured deployment looks
	// exactly like an instance with nothing pending. gcpworkflows returns nil always — Google does
	// not expose the parked step through the Executions API (gcpworkflows instance.go:113-118) —
	// and inmemoryworkflows never populates the field.
	//
	// Never conclude from an empty result that no work is waiting on a human.
	GetPendingActivities() []string

	// Subscribe is A NO-OP THAT RETURNS NIL ON EVERY ENGINE. The handler is never called.
	//
	// dsl.WorkflowInstance.Subscribe returns nil with an "implementation depends on the underlying
	// engine" comment and no implementation (dslworkflowinstance.go:46-48) — and that type is the
	// instance for goworkflows, cadence and inmemoryworkflows — while gcpworkflows' own instance
	// does the same (gcpworkflows instance.go:126-128). Nothing registers the handler anywhere, and
	// nothing returns an error to say so, so a caller sees a successful subscription followed by
	// permanent silence.
	//
	// Use WorkflowManager.Subscribe, which does reach the messaging manager, and filter on the
	// InstanceID carried in the WorkflowEvent.
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

// WorkflowManager is the engine-facing workflow contract. Four engines implement it —
// goworkflows, cadence, inmemoryworkflows and gcpworkflows — and the server implements it once more
// as a ROUTER that dispatches each call to whichever engine claimed the named workflow at load
// (laatooserver/src/core/workflowmanager.go).
//
// Application code should hold elements.WorkflowManager (the router, from
// ctx.GetServerElement(core.ServerElementWorkflowManager)) or use the RequestContext helpers —
// never an engine directly.
//
// An engine claims a definition by the workflowType declared in its .dsl metadata:
// inmemoryworkflows takes `function`, goworkflows takes `process` (or `durable` when its module
// sets mode: durable), and cadence and gcpworkflows BOTH register for `durable` — so only one of
// those two can be enabled at a time, because the second RegisterProvider for a type is refused
// (workflowmanager.go:77-97).
//
// THE ENGINES DIVERGE ON MOST OF THIS INTERFACE, including two methods that succeed while doing
// nothing. Read each method's doc before assuming a call has an effect.
type WorkflowManager interface {
	// LoadWorkflows parses every .dsl under dir, claims those whose metadata.workflowType matches
	// the type this engine registered for, and returns them keyed by name (the filename without
	// .dsl). The server calls this at startup, once per module; application code does not.
	//
	// IT RETURNS (nil, nil) WHEN dir DOES NOT EXIST — on the router and on all four engines. An
	// `if err != nil` guard passes and the caller is holding a nil map.
	//
	// THE ROUTER SWALLOWS PER-ENGINE FAILURES: an engine whose load returns an error is logged and
	// skipped, and the aggregate error stays nil (workflowmanager.go:209-214). A malformed .dsl
	// therefore drops its workflow out of the registry without failing startup, and the first
	// symptom is a NotFound from StartWorkflow much later and somewhere else.
	//
	// gcpworkflows also TRANSPILES AND DEPLOYS each claimed definition to Google here, so on that
	// engine this call performs remote work and can fail for infrastructure reasons.
	LoadWorkflows(ctx core.ServerContext, dir string, module core.Module) (map[string]Workflow, error)

	// StartWorkflow starts workflowName with initVal as its input and insconf overriding declared
	// variables, returning a handle on the new instance.
	//
	// WHETHER IT BLOCKS DEPENDS ON THE ENGINE. inmemoryworkflows runs the ENTIRE workflow
	// synchronously before returning (inmemoryworkflows workflowservice.go:145-205), so the call
	// takes as long as the process does and the returned instance is already Completed or Failed;
	// goworkflows, cadence and gcpworkflows return as soon as the run is scheduled. Do not put this
	// on a request path without knowing which engine claims the definition.
	//
	// inmemoryworkflows also returns a NON-NIL instance ALONGSIDE its error on failure — the only
	// implementation that does.
	//
	// THE CALLER'S IDENTITY IS CAPTURED HERE and every activity runs under it. gcpworkflows refuses
	// a start with no user outright (gcpworkflows workflowservice.go:551-559), because a Cloud
	// Workflows execution re-presents a signed grant on every step; the in-process engines instead
	// fall back to a system request, which BYPASSES AUTHORIZATION for the whole run. Start
	// workflows from a request that carries a user.
	//
	// A name no engine claimed is errors.NotFound from the router (workflowmanager.go:120).
	StartWorkflow(ctx core.RequestContext, workflowName string, initVal utils.StringMap, insconf utils.StringMap) (WorkflowInstance, error)

	// GetWorkflowInstance returns a fresh handle on one running or finished instance.
	//
	// workflowId IS THE WORKFLOW NAME, NOT AN ID. The router looks it up in its name-to-engine map
	// to decide whom to ask (workflowmanager.go:129-131), so passing an instance id or a uuid there
	// routes nowhere and yields NotFound. instanceId is the value GetId returned.
	//
	// The handle it builds is not a superset of the one StartWorkflow gave you: InitData,
	// GetVariables and GetError are left empty on several engines. See WorkflowInstance.
	//
	// goworkflows can only resolve instances started in THIS process — the lookup needs an
	// execution id it keeps in a process-local map (goworkflows workflowservice.go:487-497) — so it
	// fails across replicas and across restarts there.
	GetWorkflowInstance(ctx core.RequestContext, workflowId string, instanceId string) (WorkflowInstance, error)

	// IsWorkflowRegistered reports whether any engine claimed a workflow under this name, walking
	// up to the parent level as the start lookup does (workflowmanager.go:164-173).
	//
	// FALSE MEANS "NO ENGINE CLAIMED IT", which includes a definition that exists on disk and was
	// skipped — because its workflowType matches no registered engine, or because its engine's load
	// failed and the router swallowed the error. It cannot distinguish those from a name that was
	// never written.
	IsWorkflowRegistered(ctx core.ServerContext, name string) bool

	// SendSignal delivers a named signal, with signalVal as its payload, to a running instance.
	//
	// A NO-OP THAT RETURNS SUCCESS ON inmemoryworkflows: it logs "Signal received (ignored in
	// memory)" and returns nil (inmemoryworkflows workflowservice.go:217-220).
	//
	// On gcpworkflows it shares CompleteActivity's resume mechanism but WITHOUT the stored-state
	// merge, so a signal that does not itself carry the execution's callback URL is not actionable
	// (gcpworkflows workflowservice.go:664-667 and the resumeParked comment above it).
	//
	// actId is IGNORED by goworkflows, which derives the signal channel from the instance instead.
	//
	// Nothing in the platform calls this outside the manager proxy; CompleteActivity is the
	// supported way to resume a parked human step.
	SendSignal(ctx core.RequestContext, workflowId string, workflowIns string, actId string, signal string, signalVal utils.StringMap) error

	// CompleteActivity resumes an instance parked on a manual (human) activity, handing the step
	// the data map, and closes out the state snapshot the park wrote.
	//
	// A NO-OP THAT RETURNS SUCCESS ON inmemoryworkflows (inmemoryworkflows workflowservice.go:
	// 222-225) — that engine refuses manual activities outright, so nothing can be parked there.
	//
	// PASSING A NON-NIL err DOES NOT FAIL THE STEP, and both engines that can park still return
	// nil. goworkflows and gcpworkflows log it and leave the execution PARKED, deliberately:
	// neither has a rejection path, and waking the workflow with empty data would make a rejection
	// indistinguishable from an approval that carried no answer (goworkflows workflowengine.go:
	// 456-491, gcpworkflows workflowservice.go:681-706). A rejected step stays parked until it is
	// completed successfully or the run is abandoned.
	//
	// data need not carry everything the step needs: the DSL merges the snapshot taken when the
	// step parked, with the caller's values winning, and deletes that snapshot only after the
	// resume succeeds (dslworkflowinstance.go:200-257).
	//
	// actId is the activity NAME for goworkflows — it is the signal-channel key — while the
	// snapshot lookup accepts either the name or the per-execution activity uuid.
	CompleteActivity(ctx core.RequestContext, workflowId string, workflowIns string, actId string, data utils.StringMap, err error) error

	// Subscribe registers handler for one workflow or activity event type, optionally narrowed to a
	// single engine by wfType.
	//
	// wfType == "" means EVERY registered engine, and on that path the router LOGS AND DISCARDS
	// each engine's error and returns nil regardless (workflowmanager.go:246-251) — a subscription
	// that failed everywhere is reported as success. A non-empty wfType that no engine registered
	// is the one case that genuinely returns NotFound.
	//
	// NOT IMPLEMENTED ON gcpworkflows, which returns errors.NotImplemented (gcpworkflows
	// workflowservice.go:746-748) — so where Cloud Workflows holds the durable slot there are no
	// durable workflow events to subscribe to at all. The other three delegate to the messaging
	// manager, subscribing to a topic named for the event type, and return nil silently when no
	// messaging manager is configured.
	//
	// The handler receives a WorkflowEvent carrying WorkflowID, InstanceID and ActivityID. Filter
	// on those: the subscription is per event type, never per instance.
	//
	// The server does NOT subscribe to EventActivityCompleted or EventActivityFailed itself, on
	// purpose — doing so produced spurious CompleteActivity calls that corrupted the HITL signal
	// queue (workflowmanager.go:17-22).
	Subscribe(ctx core.RequestContext, wfType WorkflowType, eventType WorkflowEventType, handler core.MessageListener) error
}
