package core

type MethodArgs map[string]interface{}

//Creates object
type ObjectCreator func() interface{}

//Creates collection
type ObjectCollectionCreator func(length int) interface{}

//interface that needs to be implemented by any object provider in a system
type ObjectFactory interface {
	//Creates object
	CreateObject() interface{}
	//Creates collection
	CreateObjectCollection(length int) interface{}
}

type Initializable interface {
	Init(ctx Context, args MethodArgs) error
}
