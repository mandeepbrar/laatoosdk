// Package tracing is the contract for recording what an agent run did, step by step: the run
// itself, each plan and task, each LLM call, each skill and tool call, each handoff and each
// workflow activity, nested under the step that caused it.
//
// Three parties meet here:
//
//   - The SERVER records steps. It implements Tracer as part of the trace manager element
//     (elements.TraceManager) and starts the steps it owns -- an agent's run, every LLM
//     completion, every skill invocation -- itself.
//
//   - A PLUGIN adds steps the server cannot see, such as a goal agent's plan or a chat backend's
//     tool call, through Start. Start resolves the trace manager from the context and hands back a
//     step that does nothing when tracing is off, so the plugin code is the same in every
//     configuration (example below).
//
//   - A SINK stores finished steps. A plugin implementing TraceSink registers it with the trace
//     manager; the server hands it batches of StepRecord, already stamped with the tenant and user
//     each step ran as, so a sink writing tenant-scoped records never has to infer identity.
//
// Tracing is OFF until a sink is registered. With no sink resolving for a context's namespace,
// Start returns the context unchanged and a step whose methods do nothing.
//
// Every type in a method signature here is a core, utils, builtin or SDK type, which is what lets
// a plugin's own package main satisfy or consume these interfaces across the shared-object
// boundary.
//
// A plugin recording one tool call:
//
//	ctx, step := tracing.Start(ctx, tracing.StepTool, "search_orders", utils.StringMap{
//		tracing.AttrToolName: "search_orders",
//	})
//	defer step.End()
//	step.SetPayload(tracing.PayloadInput, string(argsJSON))
//	result, err := callTool(ctx, args)
//	step.RecordError(err)
package tracing

import (
	"time"

	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// StepKind says what a step is. A viewer groups and styles steps by kind, and an agent querying a
// run filters on it -- "every LLM call in run X" is a kind filter.
type StepKind string

const (
	// StepRun is the top step of a trace: one agent invocation from the request that started it
	// to the response, a failure, or a pause for human input.
	StepRun StepKind = "run"
	// StepPlan is an agent deciding what to do -- a goal agent building its execution graph.
	StepPlan StepKind = "plan"
	// StepTask is one unit of a plan being executed. A retried task records one StepTask per
	// attempt, each carrying AttrRetryAttempt.
	StepTask StepKind = "task"
	// StepLLM is one completion request to an LLM provider, streaming or not.
	StepLLM StepKind = "llm"
	// StepSkill is one skill invocation.
	StepSkill StepKind = "skill"
	// StepTool is one tool call an agent made -- an MCP tool, or a service exposed to the model
	// as a tool.
	StepTool StepKind = "tool"
	// StepHandoff is one agent passing work to another.
	StepHandoff StepKind = "handoff"
	// StepActivity is one workflow activity an agent drove.
	StepActivity StepKind = "activity"
	// StepResume is a paused run continuing in a later request. It sits in the original run's
	// trace, so a run that waited for a human still reads as one run.
	StepResume StepKind = "resume"
)

// StepStatus is how a step ended.
type StepStatus string

const (
	// StatusOK is a step that finished without error. It is the status of an ended step on which
	// neither RecordError nor SetStatus was called.
	StatusOK StepStatus = "ok"
	// StatusFailed is a step on which RecordError was called, or which was set failed explicitly.
	StatusFailed StepStatus = "failed"
	// StatusPaused is a step that stopped to wait for human input and will be continued by a
	// StepResume in a later request. It is normally set on the run step.
	StatusPaused StepStatus = "paused"
)

// Well-known attribute keys. The LLM keys use the OpenTelemetry generative-AI semantic convention
// names, so the same records can later be exported to an OpenTelemetry backend unchanged.
// A plugin may set keys of its own alongside these.
const (
	// AttrLLMProvider is the provider that served a completion, e.g. "anthropic".
	AttrLLMProvider = "gen_ai.system"
	// AttrLLMModel is the model a completion was requested from.
	AttrLLMModel = "gen_ai.request.model"
	// AttrLLMInputTokens is the prompt token count the provider reported.
	AttrLLMInputTokens = "gen_ai.usage.input_tokens"
	// AttrLLMOutputTokens is the completion token count the provider reported.
	AttrLLMOutputTokens = "gen_ai.usage.output_tokens"
	// AttrLLMCostUSD is the cost of a completion in US dollars, when the provider reported one.
	AttrLLMCostUSD = "laatoo.llm.cost_usd"
	// AttrLLMTimeToFirstTokenMs is, for a streaming completion, the milliseconds until the first
	// event arrived.
	AttrLLMTimeToFirstTokenMs = "laatoo.llm.time_to_first_token_ms"
	// AttrLLMRequestId is the provider's own id for the completion request.
	AttrLLMRequestId = "laatoo.llm.request_id"
	// AttrAgent is the alias of the agent a run step belongs to.
	AttrAgent = "laatoo.agent.name"
	// AttrAgentType is the type of that agent, e.g. "goal" or "workflow".
	AttrAgentType = "laatoo.agent.type"
	// AttrSessionId is the conversation session a run belongs to, when it has one.
	AttrSessionId = "laatoo.session.id"
	// AttrSkillName is the skill a StepSkill invoked.
	AttrSkillName = "laatoo.skill.name"
	// AttrToolName is the tool a StepTool called.
	AttrToolName = "laatoo.tool.name"
	// AttrHandoffTarget is the agent a StepHandoff passed work to.
	AttrHandoffTarget = "laatoo.handoff.target"
	// AttrHandoffMode is how a handoff ran, e.g. "direct" or "async".
	AttrHandoffMode = "laatoo.handoff.mode"
	// AttrActivityName is the activity a StepActivity executed.
	AttrActivityName = "laatoo.activity.name"
	// AttrRetryAttempt is the 1-based attempt number of a retried step.
	AttrRetryAttempt = "laatoo.retry.attempt"
	// AttrPauseHandle is the handle of the pause a StatusPaused run is waiting on, and the pause a
	// StepResume continued.
	AttrPauseHandle = "laatoo.pause.handle"
)

// Well-known payload names. A payload is the full text of something a step consumed or produced,
// kept apart from the step's attributes because it can be large: a sink stores it separately and a
// viewer loads it only when the step is opened.
const (
	// PayloadInput is what a step was given -- skill or tool parameters, a task's input.
	PayloadInput = "input"
	// PayloadOutput is what a step returned.
	PayloadOutput = "output"
	// PayloadPrompt is the messages sent to an LLM, in order.
	PayloadPrompt = "prompt"
	// PayloadCompletion is the text an LLM returned. For a streaming call it is the assembled
	// stream.
	PayloadCompletion = "completion"
	// PayloadPlan is a plan an agent built, e.g. a goal agent's execution graph.
	PayloadPlan = "plan"
)

// Step is a handle on one step while it runs. Every method is safe to call on the handle Start
// returns when tracing is off, and does nothing there; a caller never checks whether it is
// recording before using it.
type Step interface {
	// SetAttributes adds attributes to the step, overwriting a key already set. Values should be
	// strings, numbers or booleans; anything else is recorded as its string form.
	SetAttributes(attrs utils.StringMap)
	// SetPayload attaches the full text of something the step consumed or produced under a name,
	// normally one of the Payload* constants. Setting a name twice keeps the later text.
	SetPayload(name string, content string)
	// RecordError marks the step failed and records err's message. A nil err is ignored.
	RecordError(err error)
	// SetStatus sets how the step ended, overriding what RecordError implied. A run step waiting
	// for human input is set StatusPaused.
	SetStatus(status StepStatus)
	// End finishes the step. Calling it more than once has no further effect, and a step that is
	// never ended is never exported.
	End()
	// TraceId returns the id of the trace the step belongs to -- the run's id -- or "" when the
	// step is not recording.
	TraceId() string
	// StepId returns the step's own id, or "" when the step is not recording.
	StepId() string
	// IsRecording reports whether the step is being recorded, so a caller can skip building an
	// expensive payload that would be thrown away.
	IsRecording() bool
}

// Tracer starts and finds steps. The server implements it as part of the trace manager element;
// a plugin normally reaches it through Start and Current rather than resolving it itself.
type Tracer interface {
	// Enabled reports whether steps started from ctx are recorded, which is whether a trace sink
	// resolves for ctx's namespace.
	Enabled(ctx core.RequestContext) bool
	// StartStep starts a step as a child of ctx's current step -- or as the top step of a new
	// trace when ctx has none -- and returns a context carrying it, to pass to whatever the step
	// does, together with the step's handle. When tracing is off it returns ctx unchanged and a
	// handle that does nothing.
	StartStep(ctx core.RequestContext, kind StepKind, name string, attrs utils.StringMap) (core.RequestContext, Step)
	// CurrentStep returns the innermost step ctx carries, or a handle that does nothing when it
	// carries none. It takes the base context so that code holding only a ctx.Context, such as a
	// log handler, can read the current step's ids.
	CurrentStep(ctx ctx.Context) Step
	// RunStep returns the top step of the trace ctx belongs to when that step was started in this
	// process, or a handle that does nothing otherwise. It is how code at any depth marks the
	// run paused or adds run-level attributes.
	RunStep(ctx ctx.Context) Step
	// Inject writes ctx's current trace position into a string map that can travel where a
	// context cannot -- a message, a task, a workflow's start parameters, a pause record. It
	// returns an empty map when ctx carries no step.
	Inject(ctx ctx.Context) utils.StringMap
	// Extract returns a context continuing the trace a carrier written by Inject describes, so
	// steps started from it nest under the step that was current when the carrier was written. An
	// empty or unreadable carrier returns ctx unchanged.
	Extract(ctx core.RequestContext, carrier utils.StringMap) core.RequestContext
}

// StepRecord is one finished step as a sink receives it.
type StepRecord struct {
	// TraceId is the id of the trace -- the run -- the step belongs to.
	TraceId string
	// StepId is the step's own id.
	StepId string
	// ParentStepId is the id of the step this one is nested under, or "" for the run step.
	ParentStepId string
	// Kind says what the step is.
	Kind StepKind
	// Name is the step's display name, e.g. a model, skill, tool or agent name.
	Name string
	// Start is when the step started.
	Start time.Time
	// End is when the step ended.
	End time.Time
	// Status is how the step ended.
	Status StepStatus
	// Error is the message RecordError recorded, or "".
	Error string
	// Attributes are the step's attributes, after redaction.
	Attributes utils.StringMap
	// Payloads are the full texts attached with SetPayload, by name, after redaction.
	Payloads map[string]string
	// Redacted reports whether redaction masked any attribute or payload of this step.
	Redacted bool
	// TenantId is the id of the tenant the step ran as, captured when the step started. A sink
	// writing tenant-scoped records writes this step under this tenant.
	TenantId string
	// TenantName is that tenant's name.
	TenantName string
	// UserId is the id of the user the step ran as, captured when the step started.
	UserId string
	// UserName is that user's name.
	UserName string
	// Namespace is the namespace the step ran in, e.g. "::portal".
	Namespace string
}

// TraceSink stores finished steps. A plugin implements it and registers it with the trace manager
// (elements.TraceManager.RegisterTraceSink) at Initialize; from then on the server hands it every
// finished step of every trace started in that namespace and the namespaces below it.
type TraceSink interface {
	// ExportSteps stores a batch of finished steps. The server calls it off the request path, and
	// a step's parent may arrive in a later batch than the step itself. A returned error drops the
	// batch: the server logs it and counts the dropped steps, and the requests that produced them
	// are unaffected.
	ExportSteps(ctx core.ServerContext, steps []*StepRecord) error
}

// Start starts a step as a child of ctx's current step, through the trace manager ctx resolves.
// It returns a context carrying the step and the step's handle. When no trace manager resolves, or
// tracing is off for ctx's namespace, it returns ctx unchanged and a handle that does nothing --
// so code calling Start needs no check of its own.
func Start(ctx core.RequestContext, kind StepKind, name string, attrs utils.StringMap) (core.RequestContext, Step) {
	tracer := tracerOf(ctx)
	if tracer == nil {
		return ctx, NoopStep()
	}
	return tracer.StartStep(ctx, kind, name, attrs)
}

// Current returns the innermost step ctx carries, through the trace manager ctx resolves, or a
// handle that does nothing when there is none.
func Current(ctx core.RequestContext) Step {
	tracer := tracerOf(ctx)
	if tracer == nil {
		return NoopStep()
	}
	return tracer.CurrentStep(ctx)
}

// tracerOf resolves the trace manager from ctx, returning nil when ctx is nil or the server
// provides none -- the case on a server built before the trace manager existed.
func tracerOf(ctx core.RequestContext) Tracer {
	if ctx == nil {
		return nil
	}
	// the trace manager element; nil on a server that has none
	tracer, ok := ctx.GetServerElement(core.ServerElementTraceManager).(Tracer)
	if !ok {
		return nil
	}
	return tracer
}

// NoopStep returns a step handle that records nothing. Start returns one when tracing is off, and
// an implementation of Tracer returns one for the same reason.
func NoopStep() Step {
	return noopStep{}
}

// noopStep is the handle of a step that is not being recorded.
type noopStep struct{}

// SetAttributes does nothing.
func (noopStep) SetAttributes(utils.StringMap) {}

// SetPayload does nothing.
func (noopStep) SetPayload(string, string) {}

// RecordError does nothing.
func (noopStep) RecordError(error) {}

// SetStatus does nothing.
func (noopStep) SetStatus(StepStatus) {}

// End does nothing.
func (noopStep) End() {}

// TraceId returns "": the step belongs to no trace.
func (noopStep) TraceId() string { return "" }

// StepId returns "": the step has no id.
func (noopStep) StepId() string { return "" }

// IsRecording returns false.
func (noopStep) IsRecording() bool { return false }
