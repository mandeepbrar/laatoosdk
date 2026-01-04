package log

import (
	"log/slog"

	"laatoo.io/sdk/ctx"
)

func Trace(reqContext ctx.Context, msg string, args ...slog.Attr) {
	reqContext.LogTrace(msg, args...)
}
func Debug(reqContext ctx.Context, msg string, args ...slog.Attr) {
	reqContext.LogDebug(msg, args...)
}
func Info(reqContext ctx.Context, msg string, args ...slog.Attr) {
	reqContext.LogInfo(msg, args...)
}
func Warn(reqContext ctx.Context, msg string, args ...slog.Attr) {
	reqContext.LogWarn(msg, args...)
}
func Error(reqContext ctx.Context, msg string, args ...slog.Attr) {
	reqContext.LogError(msg, args...)
}
func Fatal(reqContext ctx.Context, msg string, args ...slog.Attr) {
	reqContext.LogFatal(msg, args...)
}
func Dump(context ctx.Context) {
	context.Dump()
}

const (
	FATAL = 1
	ERROR = 2
	WARN  = 3
	INFO  = 4
	DEBUG = 5
	TRACE = 6
)
