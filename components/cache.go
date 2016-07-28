package components

import (
	"fmt"
	"laatoo/sdk/core"
)

type CacheComponent interface {
	PutObject(ctx core.RequestContext, bucket string, key string, item interface{}) error
	GetObject(ctx core.RequestContext, bucket string, key string, objectType string) (interface{}, bool)
	GetMulti(ctx core.RequestContext, bucket string, keys []string, objectType string) map[string]interface{}
	Delete(ctx core.RequestContext, bucket string, key string) error
}

func GetCacheKey(objectType string, variants ...interface{}) string {
	return fmt.Sprintf("%s_%#v", objectType, variants)
}
