package config

import (
	"reflect"

	"laatoo.io/sdk/utils"
)

func CastToConfigArray(val interface{}) []Config {
	if val == nil {
		return nil
	}
	itemVal := reflect.ValueOf(val)
	if itemVal.Kind() != reflect.Array && itemVal.Kind() != reflect.Slice {
		return nil
	}
	len := itemVal.Len()
	res := make([]Config, len)
	for i := 0; i < len; i++ {
		res[i] = CastToConfig(itemVal.Index(i).Interface())
	}
	return res
}

func CastToConfig(val interface{}) Config {
	m := utils.CastToStringMap(val)
	if m != nil {
		return GenericConfig(m)
	}
	return nil
}
