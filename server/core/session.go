package core

import (
	"time"

	"laatoo.io/sdk/utils"
)

type Session interface {
	GetId() string
	CreationTime() time.Time
	GetUser() string
	GetString(key string) (string, bool)
	GetBool(key string) (bool, bool)
	GetInt(key string) (int, bool)
	GetStringArray(key string) ([]string, bool)
	AllKeys() []string
	GetStringMap(key string) (utils.StringMap, bool)
	GetStringsMap(key string) (utils.StringsMap, bool)
	Get(key string) (interface{}, bool)
	Set(key string, val interface{})
	SetVals(vals utils.StringMap)
	Save(RequestContext) error
	// IsSerializable returns false for sessions holding live non-serialisable
	// objects (e.g. SSE streams). The session manager stores these only in its
	// in-process sync.Map and never writes them to the CacheComponent.
	IsSerializable() bool
}
