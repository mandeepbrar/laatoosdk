package config

import (
	"fmt"
)

func FileAdapter(conf Config, configName string) (Config, error, bool) {
	var configToRet Config
	var err error
	confFileName, ok := conf.GetString(configName)
	if ok {
		configToRet, err = NewConfigFromFile(confFileName)
		if err != nil {
			return nil, fmt.Errorf("Could not read from file %s. Error:%s", confFileName, err), true
		}
	} else {
		configToRet, ok = conf.GetSubConfig(configName)
		if !ok {
			return nil, nil, false
		}
	}
	return configToRet, nil, true
}
