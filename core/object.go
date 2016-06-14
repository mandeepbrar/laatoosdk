package core

type MethodArgs map[string]interface{}

//Creates object
type ObjectCreator func(ctx Context, args MethodArgs) (interface{}, error)

//Creates collection
type ObjectCollectionCreator func(ctx Context, length int, args MethodArgs) (interface{}, error)

//interface that needs to be implemented by any object provider in a system
type ObjectFactory interface {
	//Creates object
	CreateObject(ctx Context, args MethodArgs) (interface{}, error)
	//Creates collection
	CreateObjectCollection(ctx Context, length int, args MethodArgs) (interface{}, error)
}
