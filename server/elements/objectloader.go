package elements

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type ObjectLoader interface {
	core.ServerElement
	Register(ctx ctx.Context, obj interface{}, version string, metadata core.Info) error
	RegisterObjectFactory(ctx ctx.Context, factory core.ObjectFactory, version string) error
	GetRegName(ctx ctx.Context, object interface{}) (string, bool, bool)
	CreateCollection(ctx ctx.Context, objectName string, length int) (interface{}, error)
	CreateObject(ctx ctx.Context, objectName string) (interface{}, error)
	CreateObjectPointersCollection(ctx ctx.Context, objectName string, length int) (interface{}, error)
	GetMetaData(ctx ctx.Context, objectName string) (core.Info, error)
	GetObjectFactory(ctx ctx.Context, name string) (core.ObjectFactory, bool)
	List(ctx core.ServerContext) utils.StringsMap
}
