package utils

import (
	"math/rand"
	"reflect"

	"golang.org/x/crypto/bcrypt"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Remove(arr []string, elem string) []string {
	for i, v := range arr {
		if v == elem {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

func StrContains(s []string, e string) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

func CastToInterfaceArray(val interface{}) []interface{} {
	if val == nil {
		return nil
	}
	itemVal := reflect.ValueOf(val)
	if itemVal.Kind() != reflect.Slice && itemVal.Kind() != reflect.Array {
		return nil
	}
	leng := itemVal.Len()
	res := make([]interface{}, leng)
	for i := 0; i < leng; i++ {
		res[i] = itemVal.Index(i).Interface()
	}
	return res
}

/*
func CastToArrayType(val interface{}, typ reflect.Type) reflect.Value {
	log.Println("***************************** casting")
	if val == nil {
		return reflect.MakeSlice(typ, 0, 0)
	}
	itemVal := reflect.ValueOf(val)
	log.Println("***************************** setting field item val", val, itemVal.Kind())
	if itemVal.Kind() != reflect.Slice && itemVal.Kind() != reflect.Array {
		return reflect.MakeSlice(typ, 0, 0)
	}
	leng := itemVal.Len()
	log.Println("***************************** setting field item val", leng)
	res := reflect.MakeSlice(typ, leng, leng)
	log.Println("***************************** setting field")
	for i := 0; i < leng; i++ {
		log.Println("***************************** setting field val", itemVal.Index(i), itemVal.Index(i).Interface())
		obj := reflect.New(res.Index(i).Type())
		obj.SetBytes(itemVal.Index(i).Bytes())
		//converted := itemVal.Index(i).Convert()
		res.Index(i).Set(converted)
	}
	return res
}*/

func CastToStringArray(val interface{}) []string {
	if val == nil {
		return nil
	}
	itemVal := reflect.ValueOf(val)
	if itemVal.Kind() != reflect.Array && itemVal.Kind() != reflect.Slice {
		return nil
	}
	len := itemVal.Len()
	res := make([]string, len)
	for i := 0; i < len; i++ {
		res[i] = itemVal.Index(i).Interface().(string)
	}
	return res
}

func CastToStringMap(val interface{}) map[string]interface{} {
	if val == nil {
		return nil
	}
	itemVal := reflect.ValueOf(val)
	if itemVal.Kind() != reflect.Map {
		return nil
	}
	keys := itemVal.MapKeys()
	res := make(map[string]interface{}, len(keys))
	for _, key := range keys {
		res[key.String()] = itemVal.MapIndex(key).Interface()
	}
	return res
}

func MapKeys(mapToProcess map[string]interface{}) []string {
	maplen := len(mapToProcess)
	if maplen < 1 {
		return []string{}
	}
	retVal := make([]string, maplen)
	i := 0
	for k, _ := range mapToProcess {
		retVal[i] = k
		i++
	}
	return retVal
}

func MapValues(mapToProcess map[string]interface{}) interface{} {
	maplen := len(mapToProcess)
	if maplen < 1 {
		return []interface{}{}
	}
	var arr reflect.Value
	i := 0
	for _, v := range mapToProcess {
		if i == 0 {
			sliceType := reflect.SliceOf(reflect.TypeOf(v))
			arr = reflect.MakeSlice(sliceType, maplen, maplen)
		}
		arr.Index(i).Set(reflect.ValueOf(v))
		i++
	}
	return arr.Interface()
}

func ElementPtr(object interface{}) interface{} {
	return reflect.ValueOf(object).Elem().Interface()
}

func EncryptPassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func SetObjectFields(object interface{}, newVals map[string]interface{}) {
	entVal := reflect.ValueOf(object).Elem()
	for k, v := range newVals {
		f := entVal.FieldByName(k)
		if f.IsValid() {
			if f.CanSet() {
				kind := f.Kind()
				switch kind {
				case reflect.Slice:
					{
						if f.Type().String() == "[]string" {
							arr := CastToStringArray(v)
							f.Set(reflect.ValueOf(arr))
						}
						continue
					}
				default:
					{
						f.Set(reflect.ValueOf(v))
					}
				}
			}
		}
	}
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
