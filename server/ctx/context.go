package ctx

import (
	"net/http"
	"time"

	glctx "golang.org/x/net/context"
)

type Context interface {
	GetId() string
	GetName() string
	GetPath() string
	GetParent() Context
	Get(key string) (interface{}, bool)
	GetCreationTime() time.Time
	GetElapsedTime() time.Duration
	SetGaeReq(req *http.Request)
	Set(key string, value interface{})
	SetVals(vals map[string]interface{})
	GetString(key string) (string, bool)
	GetBool(key string) (bool, bool)
	GetInt(key string) (int, bool)
	GetStringArray(key string) ([]string, bool)
	SubCtx(name string) Context
	NewCtx(flow string) Context
	GetAppengineContext() glctx.Context
	HttpClient() *http.Client
	GetOAuthContext() glctx.Context
	Dump()
	LogTrace(msg string, args ...interface{})
	LogDebug(msg string, args ...interface{})
	LogInfo(msg string, args ...interface{})
	LogWarn(msg string, args ...interface{})
	LogError(msg string, args ...interface{})
	LogFatal(msg string, args ...interface{})
}
