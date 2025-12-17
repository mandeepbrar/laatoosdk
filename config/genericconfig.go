package config

import (
	"strconv"

	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/utils"
)

type GenericConfig map[string]interface{}

func fillVariables(ctx ctx.Context, val interface{}) interface{} {
	/*expr, ok := val.(string)
	if !ok {
		return val
	}
	cont, err := utils.ProcessTemplate(ctx, []byte(expr), nil)
	if err != nil {
		return val
	}
	return string(cont)*/
	return val
}

// Get string configuration value
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

func (conf GenericConfig) Clone() Config {
	res := make(GenericConfig, len(conf))
	for k, v := range conf {
		mapV, ok := v.(GenericConfig)
		if ok {
			res[k] = mapV.Clone().(GenericConfig)
		} else {
			res[k] = v
		}

	}
	return res
}

func (conf GenericConfig) ToMap() map[string]interface{} {
	return map[string]interface{}(conf)
}

func (conf GenericConfig) GetRoot(ctx ctx.Context) (string, Config, bool) {
	confNames := conf.AllConfigurations(ctx)
	if len(confNames) == 1 {
		rootElem := confNames[0]
		rootConf, _ := conf.GetSubConfig(ctx, rootElem)
		return rootElem, rootConf, true
	}
	return "", nil, false
}

// Get string configuration value
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
		s, ok := val.(string)
		if ok {
			b, err := strconv.ParseBool(s)
			if err != nil {
				return false, false
			}
			return b, true
		}
	}
	return false, false
}

// Get string configuration value
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
		retVal, cok := val.([]Config)
		if cok {
			return retVal, true
		}
		confArr, cok := val.([]GenericConfig)
		if cok {
			retVal = make([]Config, len(confArr))
			for index, val := range confArr {
				retVal[index] = val
			}
			return retVal, true
		}
		cArr, cok := val.([]interface{})
		if !cok {
			return nil, false
		}
		retVal = make([]Config, len(cArr))
		for index, val := range cArr {
			var gc GenericConfig
			gc, ok := val.(GenericConfig)
			if !ok {
				gc, ok = val.(map[string]interface{})
				if !ok {
					return nil, false
				}
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

func (conf GenericConfig) checkConfig(ctx ctx.Context, val interface{}) (Config, bool) {
	var gc GenericConfig
	cf, ok := val.(map[string]interface{})
	if ok {
		gc = cf
		return gc, true
	} else {
		c, ok := val.(Config)
		if ok {
			return c, true
		} else {
			return nil, false
		}
	}
}

func (conf GenericConfig) GetSubConfig(ctx ctx.Context, configurationName string) (Config, bool) {
	val, found := conf[configurationName]
	if found {
		c, ok := conf.checkConfig(ctx, val)
		if ok {
			return c, true
		} else {
			/*			lookupVal := fillVariables(ctx, val)
						if lookupVal != val {
							c, ok := conf.checkConfig(ctx, lookupVal)
							if ok {
								return c, true
							}
						}*/
		}
	}
	return nil, false
}

func (conf GenericConfig) GetStringMap(ctx ctx.Context, configurationName string) (utils.StringMap, bool) {
	val, found := conf[configurationName]
	if found {
		pval, ok := val.(utils.StringMap)
		if ok {
			return pval, ok
		}
		cf, ok := val.(map[string]interface{})
		if ok {
			return cf, ok
		}
		conf, ok := val.(GenericConfig)
		if ok {
			return utils.StringMap(conf.ToMap()), ok
		}
	}
	return nil, false
}

func (conf GenericConfig) GetStringsMap(ctx ctx.Context, configurationName string) (utils.StringsMap, bool) {
	val, found := conf[configurationName]
	if found {
		cf, ok := val.(map[string]interface{})
		if ok {
			pval, ok := val.(utils.StringsMap)
			if ok {
				return pval, ok
			}
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

// Set string configuration value
func (conf GenericConfig) SetString(ctx ctx.Context, configurationName string, configurationValue string) {
	conf.Set(ctx, configurationName, configurationValue)
}

func (conf GenericConfig) Set(ctx ctx.Context, configurationName string, configurationValue interface{}) {
	conf[configurationName] = configurationValue
}

func (conf GenericConfig) SetVals(ctx ctx.Context, vals utils.StringMap) {
	if vals != nil {
		for k, v := range vals {
			conf.Set(ctx, k, v)
		}
	}
}
