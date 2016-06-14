package config

import (
	"fmt"
	"github.com/spf13/viper"
)

//Json based viper implementation of Config interface
type FileConfig struct {
	vpr *viper.Viper
}

//creates a new config object from file provided to it
//only works for json configs
func NewConfigFromFile(file string) (Config, error) {
	conf := &FileConfig{vpr: viper.New()}
	//set config type
	conf.vpr.SetConfigType("json")
	//set file name to read
	conf.vpr.SetConfigFile(file)
	//read the file
	err := conf.vpr.ReadInConfig() // Find and read the config file
	if err != nil {                // Handle errors reading the config file
		return nil, fmt.Errorf("Fatal error reading config file: %s Error: %s \n", file, err)
	}
	return conf, nil
}

func (conf *FileConfig) AllConfigurations() []string {
	return conf.vpr.AllKeys()
}

//Get string configuration value
func (conf *FileConfig) GetString(configurationName string) (string, bool) {
	if conf.vpr.InConfig(configurationName) {
		val := conf.vpr.Get(configurationName)
		str, ok := val.(string)
		if ok {
			return str, true
		}
		return "", false
	}
	return "", false
}

//Get string configuration value
func (conf *FileConfig) GetBool(configurationName string) (bool, bool) {
	if conf.vpr.InConfig(configurationName) {
		return conf.vpr.GetBool(configurationName), true
	}
	return false, false
}

func (conf *FileConfig) GetStringArray(configurationName string) ([]string, bool) {
	if conf.vpr.InConfig(configurationName) {
		return conf.vpr.GetStringSlice(configurationName), true
	}
	return nil, false
}

//Get string configuration value
func (conf *FileConfig) Get(configurationName string) (interface{}, bool) {
	if conf.vpr.InConfig(configurationName) {
		return conf.vpr.Get(configurationName), true
	}
	return nil, false
}

func (conf *FileConfig) GetConfigArray(configurationName string) ([]Config, bool) {
	if conf.vpr.InConfig(configurationName) {
		confInt := conf.vpr.Get(configurationName)
		confArr, cok := confInt.([]interface{})
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

func (conf *FileConfig) GetSubConfig(configurationName string) (Config, bool) {
	if conf.vpr.InConfig(configurationName) {
		var val GenericConfig
		val = conf.vpr.GetStringMap(configurationName)
		if val != nil {
			return val, true
		}
		return nil, false
	}
	return nil, false
}

//Set string configuration value
func (conf *FileConfig) SetString(configurationName string, configurationValue string) {
	conf.vpr.Set(configurationName, configurationValue)
}
