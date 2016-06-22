package registry

import (
	"laatoo/framework/core/objects"
	"laatoo/sdk/core"
)

//register the object factory in the global register
func RegisterObjectFactory(objectName string, factory core.ObjectFactory) {
	objects.RegisterObjectFactory(objectName, factory)
}

//register the object factory in the global register
func RegisterObject(objectName string, objectCreator core.ObjectCreator, objectCollectionCreator core.ObjectCollectionCreator) {
	objects.RegisterObject(objectName, objectCreator, objectCollectionCreator)
}

//register the object factory in the global register
func RegisterInvokableMethod(methodName string, method core.ServiceFunc) {
	objects.RegisterInvokableMethod(methodName, method)
}
