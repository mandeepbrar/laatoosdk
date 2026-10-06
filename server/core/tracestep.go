package core

import "laatoo.io/sdk/utils"

// TraceStepKind says what a recorded step is. A viewer groups and styles steps by kind, and an agent
// querying a run filters on it -- "every LLM call in run X" is a kind filter.
//
// Steps are started from a request context (RequestContext.StartTraceStep). The storage side of tracing
// -- the record a sink receives and the attribute names -- is in the tracing package.
type TraceStepKind string

const (
	// TraceStepRun is the top step of a trace: one agent invocation from the request that started it
	// to the response, a failure, or a pause for human input.
	TraceStepRun TraceStepKind = "run"
	// TraceStepPlan is an agent deciding what to do -- a goal agent building its execution graph.
	TraceStepPlan TraceStepKind = "plan"
	// TraceStepTask is one unit of a plan being executed. A retried task records one TraceStepTask per
	// attempt.
	TraceStepTask TraceStepKind = "task"
	// TraceStepLLM is one completion request to an LLM provider, streaming or not.
	TraceStepLLM TraceStepKind = "llm"
	// TraceStepDecision is one decision-model call (AgentManager.Evaluate and its typed forms):
	// its usage is recorded under the same gen_ai attribute keys as an LLM step, so a sink can
	// total both.
	TraceStepDecision TraceStepKind = "decision"
	// TraceStepSkill is one skill invocation.
	TraceStepSkill TraceStepKind = "skill"
	// TraceStepTool is one tool call an agent made -- an MCP tool, or a service exposed to the model
	// as a tool.
	TraceStepTool TraceStepKind = "tool"
	// TraceStepHandoff is one agent passing work to another.
	TraceStepHandoff TraceStepKind = "handoff"
	// TraceStepActivity is one workflow activity an agent drove.
	TraceStepActivity TraceStepKind = "activity"
	// TraceStepResume is a paused run continuing in a later request. It sits in the original run's
	// trace, so a run that waited for a human still reads as one run.
	TraceStepResume TraceStepKind = "resume"
	// TraceStepWorkflow is one workflow run, from its start to its completion, failure or
	// cancellation. It is long-lived (tracing.Tracer.StartLongLivedStep): a durable workflow
	// outlives every request and may end on another pod, so the step is recorded running when the
	// workflow starts and recorded again when it ends.
	TraceStepWorkflow TraceStepKind = "workflow"
)

// TraceStepStatus is how a recorded step ended.
type TraceStepStatus string

const (
	// TraceStepOK is a step that finished without error: the status of an ended step on which
	// neither RecordError nor SetStatus was called.
	TraceStepOK TraceStepStatus = "ok"
	// TraceStepFailed is a step on which RecordError was called, or which was set failed.
	TraceStepFailed TraceStepStatus = "failed"
	// TraceStepPaused is a step that stopped to wait for human input and will be continued by a
	// TraceStepResume in a later request. It is normally set on the run step.
	TraceStepPaused TraceStepStatus = "paused"
	// TraceStepRunning is the status a long-lived step is first recorded with, when it starts. Its
	// second record, written when it ends, carries the final status and replaces the first.
	TraceStepRunning TraceStepStatus = "running"
)

// TraceStep is a handle on one recorded step while it runs, returned by RequestContext.StartTraceStep.
//
// Every method is safe to call on the handle returned when tracing is off, and does nothing there,
// so a caller never checks whether it is recording before using it:
//
//	ctx, step := ctx.StartTraceStep(core.TraceStepTool, "search_orders", utils.StringMap{
//		tracing.AttrToolName: "search_orders",
//	})
//	defer step.End()
//	result, err := callTool(ctx, args) // pass the derived ctx on, so what the tool does nests under the step
//	step.RecordError(err)
type TraceStep interface {
	// SetAttributes adds attributes to the step, overwriting a key already set. Values should be
	// strings, numbers or booleans; anything else is recorded as its string form. Well-known keys
	// are the tracing.Attr* constants.
	SetAttributes(attrs utils.StringMap)
	// SetPayload attaches the full text of something the step consumed or produced under a name,
	// normally one of the tracing.Payload* constants. Setting a name twice keeps the later text.
	SetPayload(name string, content string)
	// RecordError marks the step failed and records err's message. A nil err is ignored.
	RecordError(err error)
	// SetStatus sets how the step ended, overriding what RecordError implied. A run step waiting
	// for human input is set TraceStepPaused.
	SetStatus(status TraceStepStatus)
	// End finishes the step. Calling it more than once has no further effect, and a step that is
	// never ended is never recorded.
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

// NoopTraceStep returns a step handle that records nothing: what StartTraceStep returns when tracing is off,
// and what a RequestContext implementation with no tracing returns from every step method.
func NoopTraceStep() TraceStep {
	return noopTraceStep{}
}

// noopTraceStep is the handle of a step that is not being recorded.
type noopTraceStep struct{}

// SetAttributes does nothing.
func (noopTraceStep) SetAttributes(utils.StringMap) {}

// SetPayload does nothing.
func (noopTraceStep) SetPayload(string, string) {}

// RecordError does nothing.
func (noopTraceStep) RecordError(error) {}

// SetStatus does nothing.
func (noopTraceStep) SetStatus(TraceStepStatus) {}

// End does nothing.
func (noopTraceStep) End() {}

// TraceId returns "": the step belongs to no trace.
func (noopTraceStep) TraceId() string { return "" }

// StepId returns "": the step has no id.
func (noopTraceStep) StepId() string { return "" }

// IsRecording returns false.
func (noopTraceStep) IsRecording() bool { return false }
