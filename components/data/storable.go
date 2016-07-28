package data

import (
	"fmt"
	"laatoo/sdk/core"
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
	RefOps          bool
}

//Object stored by data service
type Storable interface {
	core.Initializable
	Config() *StorableConfig
	GetId() string
	SetId(string)
	PreSave(ctx core.RequestContext) error
	PostSave(ctx core.RequestContext) error
	PostLoad(ctx core.RequestContext) error
	IsDeleted() bool
	Delete()
	Join(item Storable)
}

//Factory function for creating storable
//type StorableCreator func() interface{}

func CastToStorableCollection(items interface{}) ([]Storable, []string, error) {
	fmt.Errorf("Type of items ... ", items)
	arr := reflect.ValueOf(items).Elem()
	if arr.Kind() != reflect.Slice {
		return nil, nil, fmt.Errorf("Invalid cast to Storable. Type of Item: %s", arr.Kind())
	}
	length := arr.Len()
	retVal := make([]Storable, length)
	ids := make([]string, length)
	j := 0
	for i := 0; i < length; i++ {
		valPtr := arr.Index(i).Addr().Interface()
		stor, ok := valPtr.(Storable)
		if !ok {
			return nil, nil, fmt.Errorf("Invalid cast to Storable. Item: %s", valPtr)
		}
		if stor.IsDeleted() {
			continue
		}
		ids[j] = stor.GetId()
		retVal[j] = stor
		j++
	}
	return retVal[0:j], ids, nil
}

func CastToStorableHash(items interface{}) (map[string]Storable, error) {
	arr := reflect.ValueOf(items).Elem()
	if arr.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Invalid cast to Storable. Type of Item: %s", arr.Kind())
	}
	length := arr.Len()
	retVal := make(map[string]Storable, length)
	for i := 0; i < length; i++ {
		valPtr := arr.Index(i).Addr().Interface()
		stor, ok := valPtr.(Storable)
		if !ok {
			return nil, fmt.Errorf("Invalid cast to Storable. Item: %s", valPtr)
		}
		if stor.IsDeleted() {
			continue
		}
		retVal[stor.GetId()] = stor
	}
	return retVal, nil
}
