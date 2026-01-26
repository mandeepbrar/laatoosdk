package core

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/utils"
)

type StorableConfig struct {
	ObjectType        string
	LabelField        string
	PartialLoadFields []string
	FullLoadFields    []string
	PreSave           bool
	PostSave          bool
	PostUpdate        bool
	PostLoad          bool
	Trackable         bool
	Collection        string
	Cacheable         bool
	RefOps            bool
	Workflow          bool
	Multitenant       bool
}

// Object stored by data service
type Storable interface {
	Constructor(ctx.Context)
	Config() *StorableConfig
	GetId() string
	SetId(string)
	GetLabel() string
	GetVersion() string
	SetValues(ctx.Context, interface{}, utils.StringMap) error
	PreSave(ctx ctx.Context) error
	PostSave(ctx ctx.Context) error
	PostLoad(ctx ctx.Context) error
	IsMultitenant() bool
	Join(item Storable)
	GetObjectRef() interface{}
}
