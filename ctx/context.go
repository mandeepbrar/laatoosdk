package ctx

import (
	"context"
	"log/slog"
	"time"
)

type Context interface {
	context.Context
	CreateUUID() string
	GetId() string
	GetName() string
	GetPath() string
	GetParent() Context
	Get(key string) (interface{}, bool)
	GetCreationTime() time.Time
	GetElapsedTime() time.Duration
	Set(key string, value interface{})
	SetVals(vals map[string]interface{})
	GetString(key string) (string, bool)
	GetBool(key string) (bool, bool)
	GetInt(key string) (int, bool)
	GetStringArray(key string) ([]string, bool)
	SubCtx(name string) Context
	NewCtx(name string, newpath bool) Context
	GetAppengineContext() context.Context
	GetOAuthContext() context.Context
	WithCancel() (Context, context.CancelFunc)
	WithDeadline(timeout time.Time) (Context, context.CancelFunc)
	WithTimeout(timeout time.Duration) (Context, context.CancelFunc)
	WithValue(key, val interface{}) Context
	WithContext(parent context.Context) Context
	CompleteContext()
	Dump()
	LogTrace(msg string, args ...slog.Attr)
	LogDebug(msg string, args ...slog.Attr)
	LogInfo(msg string, args ...slog.Attr)
	LogWarn(msg string, args ...slog.Attr)
	LogError(msg string, args ...slog.Attr)
	LogFatal(msg string, args ...slog.Attr)
}
