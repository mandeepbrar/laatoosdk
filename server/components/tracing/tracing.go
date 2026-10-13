// Package tracing is the storage side of recording what an agent run did, step by step: the run
// itself, each plan and task, each LLM call, each skill and tool call, each handoff and each
// workflow activity, nested under the step that caused it.
//
// Steps are STARTED from a request context, with core.RequestContext.StartTraceStep, and their handle
// and kinds are core.TraceStep and core.TraceStepKind. This package holds what the rest of the system needs
// around them:
//
//   - The well-known attribute keys and payload names a step is annotated with (Attr*, Payload*),
//     so the server, the plugins and a sink agree on what an LLM step's token count is called.
//   - StepRecord, one finished step as storage receives it, already stamped with the tenant and
//     user it ran as, so a sink writing tenant-scoped records never has to infer identity.
//   - TraceSink, the contract a storage plugin implements and registers with the trace manager
//     (elements.TraceManager.RegisterTraceSink).
//   - Tracer, the surface the trace manager element offers beyond the request context: reading the
//     current step from a base context, and carrying a trace position where a context cannot go --
//     a message, a task, a workflow's start parameters, a pause record. A workflow engine plugin
//     uses it to deliver the trace from a workflow's start to each activity it invokes.
//
// Tracing is OFF until a sink is registered. With no sink resolving for a context's namespace,
// StartTraceStep returns the context unchanged and a step whose methods do nothing.
//
// Every type in a method signature here is a core, utils, builtin or SDK type, which is what lets
// a plugin's own package main satisfy or consume these interfaces across the shared-object
// boundary.
package tracing

import (
	"time"

	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
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
	// AttrLLMCachedInputTokens is the part of the input the provider read from its prompt cache and
	// billed at the cached rate, when it reported any (OpenTelemetry's GenAI name). Absent means none
	// was reported, not zero.
	AttrLLMCachedInputTokens = "gen_ai.usage.cache_read.input_tokens"
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
	// AttrSkillName is the skill a skill step invoked.
	AttrSkillName = "laatoo.skill.name"
	// AttrToolName is the tool a tool step called.
	AttrToolName = "laatoo.tool.name"
	// AttrHandoffTarget is the agent a handoff step passed work to.
	AttrHandoffTarget = "laatoo.handoff.target"
	// AttrHandoffMode is how a handoff ran, e.g. "direct" or "async".
	AttrHandoffMode = "laatoo.handoff.mode"
	// AttrActivityName is the activity an activity step executed.
	AttrActivityName = "laatoo.activity.name"
	// AttrRetryAttempt is the 1-based attempt number of a retried step.
	AttrRetryAttempt = "laatoo.retry.attempt"
	// AttrPauseHandle is the handle of the pause a paused run is waiting on, and the pause a
	// resume step continued.
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
	// PayloadLog is log records written while the step was current, one per line.
	PayloadLog = "log"
)

// Tracer is the trace manager's surface beyond the request context. RequestContext.StartTraceStep,
// CurrentTraceStep and RootTraceStep are implemented through it; a plugin calls it directly only to carry a
// trace where a request context cannot go, or to read the current step from a base context.
type Tracer interface {
	// Enabled reports whether steps started from ctx are recorded, which is whether a trace sink
	// resolves for ctx's namespace.
	Enabled(ctx core.RequestContext) bool
	// StartTraceStep is what RequestContext.StartTraceStep does: starts a step as a child of ctx's current
	// step, or as the top step of a new trace when ctx has none, and returns a context carrying it.
	StartTraceStep(ctx core.RequestContext, kind core.TraceStepKind, name string, attrs utils.StringMap) (core.RequestContext, core.TraceStep)
	// CurrentTraceStep returns the innermost step ctx carries, or a handle that does nothing when it
	// carries none. It takes the base context so that code holding only a ctx.Context, such as a
	// log handler, can read the current step.
	CurrentTraceStep(ctx ctx.Context) core.TraceStep
	// RootTraceStep returns the top step of the trace ctx belongs to when that step was started in this
	// process, or a handle that does nothing otherwise.
	RootTraceStep(ctx ctx.Context) core.TraceStep
	// Inject writes ctx's current trace position into a string map that can travel where a
	// context cannot -- a message, a task, a workflow's start parameters, a pause record. It
	// returns an empty map when ctx carries no step.
	Inject(ctx ctx.Context) utils.StringMap
	// Extract returns a context continuing the trace a carrier written by Inject describes, so
	// steps started from it nest under the step that was current when the carrier was written. An
	// empty or unreadable carrier returns ctx unchanged. It also accepts a carrier returned by
	// StartLongLivedStep, so work done for a long-lived step -- a workflow's activities -- nests
	// under it from any request or pod.
	Extract(ctx core.RequestContext, carrier utils.StringMap) core.RequestContext
	// StartLongLivedStep starts a step whose end will happen in a different request or process --
	// a workflow run, a manual activity waiting for a person. Unlike StartTraceStep it is recorded
	// at once, with status core.TraceStepRunning, and it returns a carrier identifying it as well
	// as a context carrying it. Keep the carrier with the work (a workflow instance, a pause
	// record) and pass it to EndLongLivedStep when the work ends. When tracing is off it returns
	// ctx unchanged and an empty carrier.
	StartLongLivedStep(ctx core.RequestContext, kind core.TraceStepKind, name string, attrs utils.StringMap) (core.RequestContext, utils.StringMap)
	// EndLongLivedStep ends the long-lived step a carrier from StartLongLivedStep identifies,
	// recording it again with the final status, errMsg ("" for none) and any further attributes.
	// It needs only the carrier, so it can run on any pod after any restart. The new record
	// replaces the running one. An empty or unreadable carrier is ignored.
	EndLongLivedStep(ctx core.RequestContext, carrier utils.StringMap, status core.TraceStepStatus, errMsg string, attrs utils.StringMap)
}

// StepRecord is one finished step as a sink receives it.
type StepRecord struct {
	// TraceId is the id of the trace -- the run -- the step belongs to.
	TraceId string
	// StepId is the step's own id.
	StepId string
	// ParentStepId is the id of the step this one is nested under, or "" for the top step.
	ParentStepId string
	// Kind says what the step is.
	Kind core.TraceStepKind
	// Name is the step's display name, e.g. a model, skill, tool or agent name.
	Name string
	// Start is when the step started.
	Start time.Time
	// End is when the step ended.
	End time.Time
	// Status is how the step ended.
	Status core.TraceStepStatus
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
	//
	// A step id arrives TWICE only for a long-lived step: first with status
	// core.TraceStepRunning when it starts, then with its final status when it ends, possibly from
	// another pod days later. The later record replaces the earlier one; a sink stores by step id
	// and overwrites. No other step is ever sent twice.
	ExportSteps(ctx core.ServerContext, steps []*StepRecord) error
}
