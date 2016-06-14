package config

import (
	"laatoo/sdk/utils"
)

type GenericConfig map[string]interface{}

//Get string configuration value
func (conf GenericConfig) GetString(configurationName string) (string, bool) {
	val, found := conf[configurationName]
	if found {
		str, ok := val.(string)
		if ok {
			return str, true
		}
		return "", false
	}
	return "", false
}

//Get string configuration value
func (conf GenericConfig) GetBool(configurationName string) (bool, bool) {
	val, found := conf[configurationName]
	if found {
		return val.(bool), true
	}
	return false, false
}

//Get string configuration value
func (conf GenericConfig) Get(configurationName string) (interface{}, bool) {
	val, cok := conf[configurationName]
	if cok {
		return val, true
	}
	return nil, false
}

func (conf GenericConfig) GetStringArray(configurationName string) ([]string, bool) {
	val, found := conf[configurationName]
	if found {
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

func (conf GenericConfig) GetConfigArray(configurationName string) ([]Config, bool) {
	val, found := conf[configurationName]
	if found {
		confArr, cok := val.([]interface{})
		if !cok {
			return nil, false
		}
		retVal := make([]Config, len(confArr))
		for index, val := range confArr {
			var gc GenericConfig
			gc, ok := val.(map[string]interface{})
			if !ok {
				return nil, false
			}
			retVal[index] = gc
		}
		return retVal, true
	}
	return nil, false
}

func (conf GenericConfig) AllConfigurations() []string {
	return utils.MapKeys(conf)
}

func (conf GenericConfig) GetSubConfig(configurationName string) (Config, bool) {
	val, found := conf[configurationName]
	if found {
		var gc GenericConfig
		cf, ok := val.(map[string]interface{})
		if ok {
			gc = cf
			return gc, true
		}
		return nil, false
	}
	return nil, false
}

//Set string configuration value
func (conf GenericConfig) SetString(configurationName string, configurationValue string) {
	conf[configurationName] = configurationValue
}
