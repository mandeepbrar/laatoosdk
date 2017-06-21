package log

import "laatoo/sdk/core"

func Trace(reqContext core.Context, msg string, args ...interface{}) {
	reqContext.LogTrace(msg, args...)
}
func Debug(reqContext core.Context, msg string, args ...interface{}) {
	reqContext.LogDebug(msg, args...)
}
func Info(reqContext core.Context, msg string, args ...interface{}) {
	reqContext.LogInfo(msg, args...)
}
func Warn(reqContext core.Context, msg string, args ...interface{}) {
	reqContext.LogWarn(msg, args...)
}
func Error(reqContext core.Context, msg string, args ...interface{}) {
	reqContext.LogError(msg, args...)
}
func Fatal(reqContext core.Context, msg string, args ...interface{}) {
	reqContext.LogFatal(msg, args...)
}

const (
	FATAL = 1
	ERROR = 2
	WARN  = 3
	INFO  = 4
	DEBUG = 5
	TRACE = 6
)

type Logger interface {
	Trace(reqContext core.Context, msg string, args ...interface{})
	Debug(reqContext core.Context, msg string, args ...interface{})
	Info(reqContext core.Context, msg string, args ...interface{})
	Warn(reqContext core.Context, msg string, args ...interface{})
	Error(reqContext core.Context, msg string, args ...interface{})
	Fatal(reqContext core.Context, msg string, args ...interface{})

	SetLevel(int)
	SetFormat(string)
	IsTrace() bool
	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
}
