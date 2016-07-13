package data

import (
	"fmt"
	"laatoo/sdk/core"
	"laatoo/sdk/log"
	"reflect"
)

type StorableConfig struct {
	IdField         string
	Type            string
	SoftDeleteField string
	PreSave         bool
	PostSave        bool
	PostLoad        bool
	Auditable       bool
	Collection      string
	Cacheable       bool
	NotifyNew       bool
	NotifyUpdates   bool
}

//Object stored by data service
type Storable interface {
	Config() *StorableConfig
	GetId() string
	SetId(string)
	PreSave(ctx core.RequestContext) error
	PostSave(ctx core.RequestContext) error
	PostLoad(ctx core.RequestContext) error
	IsDeleted() bool
	Delete()
}

//Factory function for creating storable
//type StorableCreator func() interface{}

func CastToStorableCollection(items interface{}) ([]Storable, error) {
	arr := reflect.ValueOf(items).Elem()
	if arr.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Invalid cast to Storable. Type of Item: %s", arr.Kind())
	}
	length := arr.Len()
	retVal := make([]Storable, length)
	for i := 0; i < length; i++ {
		valPtr := arr.Index(i).Addr().Interface()
		stor, ok := valPtr.(Storable)
		log.Logger.Error(nil, "*****Value of ok in casting", "ok", ok)
		if !ok {
			return nil, fmt.Errorf("Invalid cast to Storable. Item: %s", valPtr)
		}
		retVal[i] = stor
	}
	return retVal, nil
}
