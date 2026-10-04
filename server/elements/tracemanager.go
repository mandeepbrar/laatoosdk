package elements

import (
	"laatoo.io/sdk/server/components/tracing"
	"laatoo.io/sdk/server/core"
)

// TraceManager is the server element that records agent runs step by step, obtained with
// ctx.GetServerElement(core.ServerElementTraceManager).(elements.TraceManager).
//
// It starts and carries steps (tracing.Tracer) and hands finished ones to the trace sink that
// resolves for the namespace a step ran in. A sink registered in a namespace serves that namespace
// and every namespace below it that does not register its own. With no sink resolving, tracing is
// off for that namespace and every step started there is a handle that does nothing.
//
// A plugin rarely needs this interface: tracing.Start and tracing.Current cover adding steps. It
// is for a plugin that supplies a sink.
type TraceManager interface {
	core.ServerElement
	tracing.Tracer

	// RegisterTraceSink registers sink under name for ctx's namespace. Call it from the sink
	// plugin's Initialize. Registering the same sink under the same name again is a no-op; a
	// different sink under a name already taken in that namespace is refused, naming both.
	RegisterTraceSink(ctx core.ServerContext, name string, sink tracing.TraceSink) error

	// DroppedSteps returns how many finished steps were dropped in ctx's namespace since the
	// server started -- because their sink returned an error or panicked, or the export queue was
	// full. Steps are dropped rather than retried so that recording never slows or fails a run.
	DroppedSteps(ctx core.ServerContext) int64
}
