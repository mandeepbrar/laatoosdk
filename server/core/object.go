package core

import (
	"laatoo/sdk/common/config"
	"laatoo/sdk/server/ctx"
)

//Creates object
type ObjectCreator func(ctx.Context) interface{}

//Creates collection
type ObjectCollectionCreator func(cx ctx.Context, length int) interface{}

//interface that needs to be implemented by any object provider in a system
type ObjectFactory interface {
	//Creates object
	CreateObject(ctx.Context) interface{}
	//Creates collection
	CreateObjectCollection(cx ctx.Context, length int) interface{}
	//Get Metadata for the object
	Info() Info
}

type Initializable interface {
	Initialize(ctx ctx.Context, conf config.Config) error
}

type Serializable interface {
	ReadAll(ctx.Context, Codec, SerializableReader) error
	WriteAll(ctx.Context, Codec, SerializableWriter) error
}
