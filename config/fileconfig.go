package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"laatoo/sdk/ctx"
)

//creates a new config object from file provided to it
//only works for json configs
func NewConfigFromFile(ctx ctx.Context, file string) (Config, error) {
	conf := make(GenericConfig, 50)
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Error opening config file %s. Error: %s", file, err.Error())
	}
	if err = json.Unmarshal(fileData, &conf); err != nil {
		return nil, fmt.Errorf("Error parsing config file %s. Error: %s", file, err.Error())
	}
	return conf, nil
}
