package server

import "laatoo/sdk/core"

type ObjectLoader interface {
	core.ServerElement
	Register(ctx core.Context, objectName string, obj interface{})
	RegisterObjectFactory(ctx core.Context, objectName string, factory core.ObjectFactory)
	RegisterObject(ctx core.Context, objectName string, objectCreator core.ObjectCreator, objectCollectionCreator core.ObjectCollectionCreator)
	CreateCollection(ctx core.Context, objectName string, length int) (interface{}, error)
	CreateObject(ctx core.Context, objectName string) (interface{}, error)
	GetObjectCollectionCreator(ctx core.Context, objectName string) (core.ObjectCollectionCreator, error)
	GetObjectCreator(ctx core.Context, objectName string) (core.ObjectCreator, error)
	GetMethod(ctx core.Context, methodName string) (core.ServiceFunc, error)
}
