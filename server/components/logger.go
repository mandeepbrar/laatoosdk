package components

import (
	"log/slog"

	"laatoo.io/sdk/ctx"
)

// Logger is the server's logging component. Most code should call the package helpers in
// laatoo.io/sdk/server/log, which route to the logger already on the context; this interface is
// what those resolve to, and what ctx.GetServerElement(core.ServerElementLogger) returns.
//
// Two implementations back it, chosen by configuration: SlogLogger over log/slog
// (laatooserver/src/log/sloglogger.go) and the older SimpleLogger
// (laatooserver/src/log/simplelogger.go). They differ on Trace and on SetFormat; everything else
// below applies to both.
//
// LEVELS ARE THE INTEGERS IN laatoo.io/sdk/server/log — FATAL=1, ERROR=2, WARN=3, INFO=4, DEBUG=5,
// TRACE=6 — so HIGHER IS MORE VERBOSE, the opposite of slog's own ordering. A message is emitted
// when the configured level is at or above the message's own level, and Fatal is never suppressed.
type Logger interface {
	// Trace records the most verbose level of message. Emitted only when the configured level is
	// TRACE.
	//
	// ON THE SLOG BACKEND TRACE AND DEBUG ARE THE SAME LEVEL. Both map to slog.LevelDebug
	// (sloglogger.go:229-234) and a configured TRACE also maps to LevelDebug
	// (sloglogger.go:372-385), so with SlogLogger a level of DEBUG emits every Trace call too and
	// the two cannot be separated in the output. SimpleLogger does distinguish them
	// (simplelogger.go:35-44).
	//
	// args are evaluated by the caller before the call, so an expensive attribute is paid for even
	// when the message is discarded. Guard those on GetLevel — NOT on IsTrace, whose doc explains
	// why.
	Trace(reqContext ctx.Context, msg string, args ...slog.Attr)

	// Debug records a diagnostic message, emitted at level DEBUG or above (that is, DEBUG or
	// TRACE). On the slog backend it is indistinguishable from Trace — see Trace.
	Debug(reqContext ctx.Context, msg string, args ...slog.Attr)

	// Info records a normal operational message, emitted at level INFO or above.
	Info(reqContext ctx.Context, msg string, args ...slog.Attr)

	// Warn records a message about a condition that did not stop the operation, emitted at level
	// WARN or above.
	Warn(reqContext ctx.Context, msg string, args ...slog.Attr)

	// Error records a failure, emitted at level ERROR or above. It only writes a log line: it does
	// not create, wrap or return an error, and callers must still return one of their own.
	Error(reqContext ctx.Context, msg string, args ...slog.Attr)

	// Fatal records a message at the highest severity. IT DOES NOT EXIT, PANIC, OR STOP ANYTHING.
	//
	// Neither implementation terminates the process. SlogLogger logs it at slog.LevelError,
	// identically to Error (sloglogger.go:244-246); SimpleLogger prints it and returns
	// (simplelogger.go:60-62). Code written as `logger.Fatal(...)` followed by an assumed-
	// unreachable path keeps running, carrying whatever invalid state prompted the call. Return an
	// error, or exit explicitly.
	//
	// It is the one level NEVER suppressed by the configured level — SimpleLogger applies no level
	// test to it at all — which makes it the level for something that must appear in a production
	// log regardless of configuration.
	Fatal(reqContext ctx.Context, msg string, args ...slog.Attr)

	// SetLevel changes verbosity at runtime, taking one of the FATAL..TRACE constants from
	// laatoo.io/sdk/server/log (higher is more verbose).
	//
	// It takes effect immediately on both implementations; on SlogLogger it updates the shared
	// slog.LevelVar the handler filters on (sloglogger.go:387-393). AN OUT-OF-RANGE VALUE IS NOT
	// REJECTED: the slog mapping sends anything it does not recognise to slog.LevelError
	// (sloglogger.go:372-385), so passing 0 or a negative number quietly silences everything below
	// Error rather than failing.
	//
	// The change applies to this logger instance, which many contexts may share.
	SetLevel(int)

	// SetFormat selects the output format — "json", or one of the happy/happycolor console formats.
	//
	// A SILENT NO-OP ON THE SLOG BACKEND. SlogLogger.SetFormat records the string and does nothing
	// else (sloglogger.go:365-371): the handler is chosen from the format at construction and is
	// never rebuilt, so the output keeps the shape it started with while GetFormat reports the new
	// name. Only SimpleLogger actually swaps the printer, and only for a format present in its
	// table — an unrecognised name there leaves the previous printer in place, also without
	// complaint (simplelogger.go:64-70).
	//
	// Set the format in configuration; this cannot be relied on to change it.
	SetFormat(string)

	// GetLevel returns the configured level as one of the FATAL..TRACE constants in
	// laatoo.io/sdk/server/log. This — not the Is* predicates — is what to compare against when
	// guarding an expensive log call: GetLevel() >= log.DEBUG.
	GetLevel() int

	// GetFormat returns the format name last set. On the slog backend that is what SetFormat was
	// told rather than what the handler is emitting; see SetFormat.
	GetFormat() string

	// IsTrace reports whether the configured level is EXACTLY TRACE. IsDebug, IsInfo and IsWarn are
	// the same exact-equality test against their own level, on both implementations
	// (sloglogger.go:394-405, simplelogger.go:75-86).
	//
	// THEY ARE NOT "IS THIS LEVEL ENABLED" TESTS, which is what they read as and how they get used.
	// At level TRACE, Debug messages ARE emitted but IsDebug() returns false — so the usual guard
	//
	//	if logger.IsDebug() { logger.Debug(ctx, "…", expensiveAttr()) }
	//
	// suppresses exactly the output an operator raised the verbosity to obtain, and does so
	// precisely at the most verbose setting. Compare GetLevel() instead.
	IsTrace() bool

	// IsDebug reports whether the configured level is EXACTLY DEBUG — not whether Debug output is
	// enabled. See IsTrace.
	IsDebug() bool

	// IsInfo reports whether the configured level is EXACTLY INFO — not whether Info output is
	// enabled. See IsTrace.
	IsInfo() bool

	// IsWarn reports whether the configured level is EXACTLY WARN — not whether Warn output is
	// enabled. See IsTrace.
	IsWarn() bool
}
