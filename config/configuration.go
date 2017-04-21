package config

//Config Interface used by Laatoo
type Config interface {
	GetString(configurationName string) (string, bool)
	GetBool(configurationName string) (bool, bool)
	GetStringArray(configurationName string) ([]string, bool)
	GetSubConfig(configurationName string) (Config, bool)
	GetConfigArray(configurationName string) ([]Config, bool)
	Get(configurationName string) (interface{}, bool)
	SetString(configurationName string, configurationValue string)
	AllConfigurations() []string
}

func Cast(conf interface{}) (Config, bool) {
	var gc GenericConfig
	cf, ok := conf.(map[string]interface{})
	if ok {
		gc = cf
		return gc, true
	}
	return nil, false
}

func Merge(conf1 Config, conf2 Config) Config {
	mergedConf := make(GenericConfig)
	copyConfs := func(conf Config) {
		if conf == nil {
			return
		}
		confNames := conf.AllConfigurations()
		for _, confName := range confNames {
			val, _ := conf.Get(confName)
			subConf, ok := val.(Config)
			if ok {
				existingVal, eok := mergedConf[confName]
				if eok {
					existingConf, cok := existingVal.(Config)
					if cok {
						mergedSubConf := Merge(existingConf, subConf)
						mergedConf[confName] = mergedSubConf
					} else {
						mergedConf[confName] = val
					}
				} else {
					mergedConf[confName] = val
				}
			} else {
				mergedConf[confName] = val
			}
		}
	}
	copyConfs(conf1)
	copyConfs(conf2)
	return mergedConf
}
