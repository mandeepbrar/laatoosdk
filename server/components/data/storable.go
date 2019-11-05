package data

import (
	"fmt"
	"laatoo/sdk/server/core"
	"reflect"

	"github.com/golang/protobuf/proto"
)

type StorableConfig struct {
	IdField           string
	LabelField        string
	PartialLoadFields []string
	FullLoadFields    []string
	Type              string
	SoftDeleteField   string
	PreSave           bool
	PostSave          bool
	PostUpdate        bool
	PostLoad          bool
	Auditable         bool
	Collection        string
	Cacheable         bool
	RefOps            bool
	Multitenant       bool
}

//Object stored by data service
type Storable interface {
	proto.Message
	core.Initializable
	Config() *StorableConfig
	GetId() string
	SetId(string)
	GetLabel() string
	SetValues(interface{}, map[string]interface{})
	PreSave(ctx core.RequestContext) error
	PostSave(ctx core.RequestContext) error
	PostLoad(ctx core.RequestContext) error
	IsDeleted() bool
	IsMultitenant() bool
	Delete()
	Join(item Storable)
}

type StorableMT interface {
	Storable
	GetTenant() string
	SetTenant(tenant string)
}

type StorableRef struct {
	Id     string   `json:"Id" bson:"Id" gorm:"column:Id" protobuf:"bytes,51,opt,name=id,proto3" sql:"type:varchar(100);`
	Type   string   `json:"Type" bson:"Type" gorm:"column:Type" protobuf:"bytes,59,opt,name=type,proto3" sql:"type:varchar(100);`
	Name   string   `json:"Name" bson:"Name" gorm:"column:Name" protobuf:"bytes,60,opt,name=name,proto3" sql:"type:varchar(300);`
	Entity Storable `json:"-" bson:"-" sql:"-" protobuf:"group,64,opt,name=Entity,proto3"`
}

func StorableArrayToMap(items []Storable) map[string]Storable {
	res := make(map[string]Storable, len(items))
	for _, item := range items {
		res[item.GetId()] = item
	}
	return res
}

//Factory function for creating storable
//type StorableCreator func() interface{}

func CastToStorableCollection(items interface{}) ([]Storable, []string, error) {
	arr := reflect.ValueOf(items)
	if arr.Kind() == reflect.Ptr {
		arr = arr.Elem()
	}
	if arr.Kind() != reflect.Slice {
		return nil, nil, fmt.Errorf("Invalid cast to Storable. Type of Item: %s", arr.Kind())
	}
	length := arr.Len()
	retVal := make([]Storable, length)
	ids := make([]string, length)
	j := 0
	for i := 0; i < length; i++ {
		itemKind := arr.Index(i).Kind()
		var valPtr interface{}
		if itemKind == reflect.Ptr {
			valPtr = arr.Index(i).Interface()
		} else {
			valPtr = arr.Index(i).Addr().Interface()
		}
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
	arr := reflect.ValueOf(items)
	if arr.Kind() == reflect.Ptr {
		arr = arr.Elem()
	}
	if arr.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Invalid cast to Storable. Type of Item: %s", arr.Kind())
	}
	length := arr.Len()
	retVal := make(map[string]Storable, length)
	for i := 0; i < length; i++ {
		itemKind := arr.Index(i).Kind()
		var valPtr interface{}
		if itemKind == reflect.Ptr {
			valPtr = arr.Index(i).Interface()
		} else {
			valPtr = arr.Index(i).Addr().Interface()
		}
		stor, ok := valPtr.(Storable)
		if !ok {
			return nil, fmt.Errorf("Invalid cast to Storable. Item: %s %s %t", valPtr, arr.Index(i).Kind(), arr.Index(i).IsNil())
		}
		if stor.IsDeleted() {
			continue
		}
		retVal[stor.GetId()] = stor
	}
	return retVal, nil
}
