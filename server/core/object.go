package core

import (
	"laatoo/sdk/common/config"
	"laatoo/sdk/server/ctx"
)

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
	//Get Metadata for the object
	Info() Info
}

type Initializable interface {
	Initialize(ctx ctx.Context, conf config.Config) error
}
