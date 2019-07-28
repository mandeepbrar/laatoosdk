package utils

import (
	"github.com/imdario/mergo"
)

type StringMap map[string]interface{}
type StringsMap map[string]string

//Get string configuration value
func (smap StringMap) GetString(key string) (string, bool) {
	val, found := smap[key]
	if found {
		str, ok := val.(string)
		if ok {
			return str, true
		}
	}
	return "", false
}

func (smap StringMap) Clone() StringMap {
	res := make(StringMap, len(smap))
	for k, v := range smap {
		mapV, ok := v.(StringMap)
		if ok {
			res[k] = mapV.Clone()
		} else {
			res[k] = v
		}
	}
	return res
}

//Get string configuration value
func (smap StringMap) GetBool(key string) (bool, bool) {
	val, found := smap[key]
	if found {
		b, ok := val.(bool)
		if ok {
			return b, true
		}
	}
	return false, false
}

func (smap StringMap) GetInt(key string) (int, bool) {
	val, found := smap[key]
	if found {
		b, ok := val.(int)
		if ok {
			return b, true
		}
	}
	return -1, false
}

func (smap StringMap) GetStringArray(key string) ([]string, bool) {
	val, found := smap[key]
	if found {
		strarr, sok := val.([]string)
		if sok {
			return strarr, true
		}

		arr, cok := val.([]interface{})
		if !cok {
			return nil, false
		}
		retVal := make([]string, len(arr))
		var ok bool
		for index, val := range arr {
			retVal[index], ok = val.(string)
			if !ok {
				return nil, false
			}
		}
		return retVal, true
	}
	return nil, false
}

func (smap StringMap) AllKeys() []string {
	return MapKeys(smap)
}

func (smap StringMap) GetStringMap(key string) (StringMap, bool) {
	val, found := smap[key]
	if found {
		cf, ok := val.(map[string]interface{})
		if ok {
			return cf, ok
		}
		imap, ok := val.(map[interface{}]interface{})
		if ok {
			res := make(map[string]interface{}, len(imap))
			for k, v := range imap {
				strkey, ok := k.(string)
				if !ok {
					return nil, false
				}
				res[strkey] = v
			}
		}
	}
	return nil, false
}

func (smap StringMap) GetStringsMap(key string) (StringsMap, bool) {
	val, found := smap[key]
	if found {
		cf, ok := val.(map[string]interface{})
		if ok {
			sm := make(map[string]string)
			for key, val := range cf {
				strval, ok := val.(string)
				if !ok {
					return nil, false
				}
				sm[key] = strval
			}
			return sm, true
		} else {
			res, ok := val.(map[string]string)
			if ok {
				return res, ok
			}
		}
	}
	return nil, false
}

func (smap StringMap) Set(key string, val interface{}) {
	smap[key] = val
}

func (smap StringMap) SetVals(vals StringMap) {
	if vals != nil {
		for k, v := range vals {
			smap.Set(k, v)
		}
	}
}

func MergeMaps(obj1, obj2 map[string]interface{}) map[string]interface{} {
	if obj1 == nil {
		return obj2
	}
	if obj2 == nil {
		return obj1
	}
	res := make(map[string]interface{})
	mergo.Merge(&res, obj1)
	mergo.Merge(&res, obj2)
	return res
}

func ShallowMergeMaps(obj1, obj2 map[string]interface{}) map[string]interface{} {
	if obj1 == nil {
		return obj2
	}
	if obj2 == nil {
		return obj1
	}
	res := make(map[string]interface{}, len(obj1))
	for k, v := range obj1 {
		res[k] = v
	}
	for k, v := range obj2 {
		res[k] = v
	}
	return res
}
