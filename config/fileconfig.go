package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var confVariables map[string]string

func init() {
	confVariables = make(map[string]string, 0)
	fil := os.Getenv("LAATOO_CONF_VARS")
	if len(fil) == 0 {
		fil = "confvariables.json"
	} else {
		_, err := os.Stat(fil)
		if err != nil {
			fil = "confvariables.json"
		}
	}
	vardata, err := ioutil.ReadFile(fil)
	if err == nil {
		err = json.Unmarshal(vardata, &confVariables)
		if err != nil {
			fmt.Println("Could not read conf variables")
		}
	}
}

//Json based viper implementation of Config interface
type FileConfig struct {
	root GenericConfig
}

//creates a new config object from file provided to it
//only works for json configs
func NewConfigFromFile(file string) (Config, error) {
	rootElement := make(GenericConfig, 50)
	fileData, err := ioutil.ReadFile(file)
	if err == nil {
		fileStr := string(fileData)
		for variable, value := range confVariables {
			fileStr = strings.Replace(fileStr, variable, value, -1)
		}
		fileData = []byte(fileStr)
	} else {
		return nil, fmt.Errorf("Error opening config file %s. Error: %s", file, err.Error())
	}
	if err = json.Unmarshal(fileData, &rootElement); err != nil {
		return nil, fmt.Errorf("Error parsing config file %s. Error: %s", file, err.Error())
	}
	conf := &FileConfig{root: rootElement}

	return conf, nil
}

func (conf *FileConfig) AllConfigurations() []string {
	return conf.root.AllConfigurations()
}

//Get string configuration value
func (conf *FileConfig) GetString(configurationName string) (string, bool) {
	return conf.root.GetString(configurationName)
}

//Get string configuration value
func (conf *FileConfig) GetBool(configurationName string) (bool, bool) {
	return conf.root.GetBool(configurationName)
}

func (conf *FileConfig) GetStringArray(configurationName string) ([]string, bool) {
	return conf.root.GetStringArray(configurationName)
}

//Get string configuration value
func (conf *FileConfig) Get(configurationName string) (interface{}, bool) {
	return conf.root.Get(configurationName)
}

func (conf *FileConfig) GetConfigArray(configurationName string) ([]Config, bool) {
	return conf.root.GetConfigArray(configurationName)
}

func (conf *FileConfig) GetSubConfig(configurationName string) (Config, bool) {
	return conf.root.GetSubConfig(configurationName)
}

//Set string configuration value
func (conf *FileConfig) SetString(configurationName string, configurationValue string) {
	conf.root.SetString(configurationName, configurationValue)
}
