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
