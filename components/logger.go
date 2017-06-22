package components

import "laatoo/sdk/core"

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
