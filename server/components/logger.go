package components

import (
	"log/slog"

	"laatoo.io/sdk/ctx"
)

type Logger interface {
	Trace(reqContext ctx.Context, msg string, args ...slog.Attr)
	Debug(reqContext ctx.Context, msg string, args ...slog.Attr)
	Info(reqContext ctx.Context, msg string, args ...slog.Attr)
	Warn(reqContext ctx.Context, msg string, args ...slog.Attr)
	Error(reqContext ctx.Context, msg string, args ...slog.Attr)
	Fatal(reqContext ctx.Context, msg string, args ...slog.Attr)

	SetLevel(int)
	SetFormat(string)
	GetLevel() int
	GetFormat() string
	IsTrace() bool
	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
}
