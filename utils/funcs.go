package utils

import (
	"math/rand"
	"reflect"
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
	if itemVal.Kind() != reflect.Array {
		return nil
	}
	len := itemVal.Len()
	res := make([]interface{}, len)
	for i := 0; i < len; i++ {
		res[i] = itemVal.Index(i).Interface()
	}
	return res
}

func CastToStringArray(val interface{}) []string {
	if val == nil {
		return nil
	}
	itemVal := reflect.ValueOf(val)
	if itemVal.Kind() != reflect.Array {
		return nil
	}
	len := itemVal.Len()
	res := make([]string, len)
	for i := 0; i < len; i++ {
		res[i] = itemVal.Index(i).String()
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

func SetObjectFields(object interface{}, newVals map[string]interface{}) {
	entVal := reflect.ValueOf(object).Elem()
	for k, v := range newVals {
		f := entVal.FieldByName(k)
		if f.IsValid() {
			// A Value can be changed only if it is
			// addressable and was not obtained by
			// the use of unexported struct fields.
			if f.CanSet() {
				f.Set(reflect.ValueOf(v))
			}
		}
	}
}
