package utils

import (
	"fmt"
	"reflect"
	"strings"

	"laatoo.io/sdk/config"
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/utils"
)

/*
type ObjectCreator interface {
	CreateObjectPointersCollection(objectName string, length int) (interface{}, error)
	CreateObject(objectName string) (interface{}, error)
}*/

type LookupFunc func(interface{}, string, interface{}) (interface{}, error)

// object: object for which fields are to be set
// newvals: values to be set on the object
// mappings: if fields from the map need to be set to specific fields of the object
// field processor: if values need to be transformed from the map before being set on the object
func SetObjectFields(ctx ctx.Context, object interface{}, newVals map[string]interface{},
	mappings map[string]string, fieldProcessor map[string]LookupFunc) error {
	var err error
	entVal := reflect.ValueOf(object).Elem()
	for k, v := range newVals {
		objField := k
		objVal := v
		if mappings != nil {
			newfld, ok := mappings[k]
			if ok {
				objField = newfld
			}
		}
		if fieldProcessor != nil {
			tfunc, ok := fieldProcessor[objField]
			if ok {
				objVal, err = tfunc(ctx, objField, objVal)
				if err != nil {
					return err
				}
			}
		}
		if objVal == nil {
			continue
		}
		f := entVal.FieldByName(objField)
		if f.IsValid() {
			if f.CanSet() {
				kind := f.Kind()
				switch kind {
				case reflect.Slice:
					{
						if f.Type().String() == "[]string" {
							arr := utils.CastToStringArray(objVal)
							f.Set(reflect.ValueOf(arr))
						} else if f.Type().String() == "[]config.Config" {
							arr := config.CastToConfigArray(objVal)
							f.Set(reflect.ValueOf(arr))
						} else if f.Type().String() == "[]map[string]interface{}" {
							arr := utils.CastToMapArray(objVal)
							f.Set(reflect.ValueOf(arr))
						}
						continue
					}
				case reflect.Struct:
					{
						objCreator := ctx.(ObjectCreator)
						objType, isPtr := GetRegisteredName(f.Type())
						vtype, visPtr := GetRegisteredName(reflect.ValueOf(v).Type())
						if objType == vtype {
							if (isPtr && visPtr) || (!isPtr && !visPtr) {
								f.Set(reflect.ValueOf(v))
							} else {
								if visPtr {
									f.Set(reflect.ValueOf(v).Elem())
								}
							}
						} else {
							structVals, ok := objVal.(map[string]interface{})
							if ok {
								structobj, err := objCreator.CreateObject(objType)
								if err != nil {
									return err
								}
								err = SetObjectFields(ctx, structobj, structVals, mappings, fieldProcessor)
								if err != nil {
									return err
								}
								if isPtr {
									f.Set(reflect.ValueOf(structobj).Convert(f.Type()))
								} else {
									f.Set(reflect.ValueOf(structobj).Elem().Convert(f.Type()))
								}
							}
						}
					}
				default:
					{
						f.Set(reflect.ValueOf(objVal).Convert(f.Type()))
					}
				}
			}
		}
	}
	return nil
}

func GetRegisteredName(typ reflect.Type) (regName string, isptr bool) {
	for {
		kind := typ.Kind()
		if kind == reflect.Ptr {
			typ = typ.Elem()
			isptr = true
			continue
		}
		if kind == reflect.Array || kind == reflect.Slice {
			typ = typ.Elem()
			isptr = false
			continue
		}
		break
	}

	pkg := typ.PkgPath()
	if pkg != "" {
		regName = fmt.Sprintf("%s.%s", strings.ReplaceAll(typ.PkgPath(), "/", "."), typ.Name())
	} else {
		regName = typ.Name()
	}
	return
}

func GetObjectFields(object interface{}, fields []string) map[string]interface{} {
	entVal := reflect.ValueOf(object).Elem()
	vals := make(map[string]interface{}, len(fields))
	for _, v := range fields {
		f := entVal.FieldByName(v)
		vals[v] = f.Interface()
	}
	return vals
}

// ObjectCreator is the narrow slice of the context that this package needs in order to build
// nested objects while setting fields. Both core.ServerContext and core.RequestContext satisfy it,
// which is why SetObjectFields type-asserts its ctx to this interface — an assertion that PANICS,
// unchecked, and only on the code path that hits a struct field (see SetObjectFields above).
//
// Objects are created from the platform's object register, so a name must have been registered —
// by a plugin's manifest, by entity codegen, or as an internal object — before either method can
// build it. An unregistered name is CORE_ERROR_PROVIDER_NOT_FOUND, never a zero value
// (laatooserver/src/core/objectloader.go:126-139).
//
// Registered names use DOTS for the Go package path (laatoo.server.engine.http.xyz), because
// registration converts the path's slashes to dots and lookup is an exact string match.
type ObjectCreator interface {
	// CreateObjectPointersCollection returns an empty slice of POINTERS to the named type
	// ([]*T), sized length — the shape the data tier fills when reading many rows.
	//
	// Contrast core.ServerContext.CreateCollection, which returns []T. Neither runs the per-object
	// constructor or field initialisers that CreateObject applies; the elements are the slice's
	// zero values.
	CreateObjectPointersCollection(objectName string, length int) (interface{}, error)

	// CreateObject returns a new instance of the named registered type as a POINTER
	// (interface{} holding *T), ready to be type-asserted.
	//
	// This is the only correct way to make a platform object — never a struct literal. Beyond
	// allocating, it initialises the declared pointer fields the platform injects (StorageInfo,
	// TenantInfo and similar) and calls the type's Constructor and SetSelfReference when it has
	// them (laatooserver/src/core/objectadapter.go:74-104). A struct literal skips all of that and
	// the omission does not surface until something dereferences one of those fields.
	CreateObject(objectName string) (interface{}, error)
}
