package components

import "laatoo/sdk/core"

type CacheComponent interface {
	PutObject(ctx core.RequestContext, bucket string, key string, item interface{}) error
	PutObjects(ctx core.RequestContext, bucket string, vals map[string]interface{}) error
	GetObject(ctx core.RequestContext, bucket string, key string, objectType string) (interface{}, bool)
	Get(ctx core.RequestContext, bucket string, key string) (interface{}, bool)
	GetObjects(ctx core.RequestContext, bucket string, keys []string, objectType string) map[string]interface{}
	GetMulti(ctx core.RequestContext, bucket string, keys []string) map[string]interface{}
	Delete(ctx core.RequestContext, bucket string, key string) error
	Increment(ctx core.RequestContext, bucket string, key string) error
	Decrement(ctx core.RequestContext, bucket string, key string) error
}
