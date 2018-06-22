package components

import "laatoo/sdk/ctx"

type CacheComponent interface {
	PutObject(ctx ctx.Context, bucket string, key string, item interface{}) error
	PutObjects(ctx ctx.Context, bucket string, vals map[string]interface{}) error
	GetObject(ctx ctx.Context, bucket string, key string, objectType string) (interface{}, bool)
	Get(ctx ctx.Context, bucket string, key string) (interface{}, bool)
	GetObjects(ctx ctx.Context, bucket string, keys []string, objectType string) map[string]interface{}
	GetMulti(ctx ctx.Context, bucket string, keys []string) map[string]interface{}
	Delete(ctx ctx.Context, bucket string, key string) error
	Increment(ctx ctx.Context, bucket string, key string) error
	Decrement(ctx ctx.Context, bucket string, key string) error
}
