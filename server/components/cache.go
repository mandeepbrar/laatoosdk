package components

import (
	"time"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type CacheComponent interface {
	PutTempObject(ctx core.RequestContext, bucket string, key string, item interface{}, ttl time.Duration) error
	PutObject(ctx core.RequestContext, bucket string, key string, item interface{}) error
	PutObjects(ctx core.RequestContext, bucket string, vals utils.StringMap) error
	GetObject(ctx core.RequestContext, bucket string, key string, objectType string) (interface{}, bool)
	GetIntoObject(ctx core.RequestContext, bucket string, key string, obj interface{}) error
	Get(ctx core.RequestContext, bucket string, key string) (interface{}, bool)
	GetObjects(ctx core.RequestContext, bucket string, keys []string, objectType string) utils.StringMap
	GetMulti(ctx core.RequestContext, bucket string, keys []string) utils.StringMap
	Delete(ctx core.RequestContext, bucket string, key string) error
	Increment(ctx core.RequestContext, bucket string, key string) error
	Decrement(ctx core.RequestContext, bucket string, key string) error
}
