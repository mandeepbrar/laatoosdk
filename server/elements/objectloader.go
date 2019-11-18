package elements

import (
	"laatoo/sdk/server/core"
	"laatoo/sdk/server/ctx"
)

type ObjectLoader interface {
	core.ServerElement
	Register(ctx ctx.Context, obj interface{}, metadata core.Info) error
	RegisterObjectFactory(ctx ctx.Context, factory core.ObjectFactory) error
	RegisterObject(ctx ctx.Context, objectCreator core.ObjectCreator, objectCollectionCreator core.ObjectCollectionCreator, metadata core.Info) error
	GetRegName(ctx ctx.Context, object interface{}) string
	CreateCollection(ctx ctx.Context, objectName string, length int) (interface{}, error)
	CreateObject(ctx ctx.Context, objectName string) (interface{}, error)
	GetMetaData(ctx ctx.Context, objectName string) (core.Info, error)
	GetObjectFactory(ctx ctx.Context, name string) (core.ObjectFactory, bool)
}
