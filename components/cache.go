package components

import (
	"fmt"
	"laatoo/sdk/core"
)

type CacheComponent interface {
	PutObject(ctx core.RequestContext, key string, item interface{}) error
	GetObject(ctx core.RequestContext, key string, val interface{}) bool
	GetMulti(ctx core.RequestContext, keys []string, val map[string]interface{}) bool
	Delete(ctx core.RequestContext, key string) error
}

func GetCacheKey(objectType string, variants ...interface{}) string {
	return fmt.Sprintf("%s_%#v", objectType, variants)
}
