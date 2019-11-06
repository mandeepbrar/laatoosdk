package elements

import (
	"laatoo/sdk/server/core"
	"laatoo/sdk/server/ctx"
)

type ObjectLoader interface {
	core.ServerElement
	Register(ctx ctx.Context, objectName string, obj interface{}, metadata core.Info)
	RegisterObjectFactory(ctx ctx.Context, objectName string, factory core.ObjectFactory)
	RegisterObject(ctx ctx.Context, objectName string, objectCreator core.ObjectCreator, objectCollectionCreator core.ObjectCollectionCreator, metadata core.Info)
	CreateCollection(ctx ctx.Context, objectName string, length int) (interface{}, error)
	CreateObject(ctx ctx.Context, objectName string) (interface{}, error)
	GetMetaData(ctx ctx.Context, objectName string) (core.Info, error)
	GetObjectFactory(ctx ctx.Context, name string) (core.ObjectFactory, bool)
}
