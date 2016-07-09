package data

import (
	"fmt"
	"laatoo/sdk/core"
	"reflect"
)

//Object stored by data service
type Storable interface {
	GetId() string
	GetObjectType() string
	PreSave(ctx core.RequestContext) error
	PostSave(ctx core.RequestContext) error
	PostLoad(ctx core.RequestContext) error
	GetIdField() string
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
		if !ok {
			return nil, fmt.Errorf("Invalid cast to Storable. Item: %s", valPtr)
		}
		retVal[i] = stor
	}
	return retVal, nil
}
