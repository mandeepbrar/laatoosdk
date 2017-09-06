package config

import (
	"laatoo/sdk/ctx"
	"laatoo/sdk/utils"
)

type GenericConfig map[string]interface{}

//Get string configuration value
func (conf GenericConfig) GetString(ctx ctx.Context, configurationName string) (string, bool) {
	val, found := conf[configurationName]
	if found {
		str, ok := fillVariables(ctx, val).(string)
		if ok {
			return str, true
		}
		return "", false
	}
	return "", false
}

//Get string configuration value
func (conf GenericConfig) GetBool(ctx ctx.Context, configurationName string) (bool, bool) {
	val, found := conf[configurationName]
	if found {
		b, ok := val.(bool)
		if ok {
			return b, true
		}
		val = fillVariables(ctx, val)
		b, ok = val.(bool)
		if ok {
			return b, true
		}
	}
	return false, false
}

//Get string configuration value
func (conf GenericConfig) Get(ctx ctx.Context, configurationName string) (interface{}, bool) {
	val, cok := conf[configurationName]
	if cok {
		return fillVariables(ctx, val), true
	}
	return nil, false
}

func (conf GenericConfig) GetStringArray(ctx ctx.Context, configurationName string) ([]string, bool) {
	val, found := conf[configurationName]
	if found {
		arr, cok := val.([]interface{})
		if !cok {
			return nil, false
		}
		retVal := make([]string, len(arr))
		var ok bool
		for index, val := range arr {
			retVal[index], ok = fillVariables(ctx, val).(string)
			if !ok {
				return nil, false
			}
		}
		return retVal, true
	}
	return nil, false
}

func (conf GenericConfig) GetConfigArray(ctx ctx.Context, configurationName string) ([]Config, bool) {
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

func (conf GenericConfig) AllConfigurations(ctx ctx.Context) []string {
	return utils.MapKeys(conf)
}

func (conf GenericConfig) GetSubConfig(ctx ctx.Context, configurationName string) (Config, bool) {
	val, found := conf[configurationName]
	if found {
		var gc GenericConfig
		cf, ok := val.(map[string]interface{})
		if ok {
			gc = cf
			return gc, true
		} else {
			c, ok := val.(Config)
			if ok {
				return c, true
			}
		}
		return nil, false
	}
	return nil, false
}

//Set string configuration value
func (conf GenericConfig) SetString(ctx ctx.Context, configurationName string, configurationValue string) {
	conf[configurationName] = configurationValue
}
