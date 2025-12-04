package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/utils"
)

type Configuration interface {
	GetName() string
	GetDescription() string
	IsRequired() bool
	GetDefaultValue() interface{}
	GetValue() interface{}
	GetType() datatypes.DataType
}

type ConfigurableObject interface {
	GetName() string
	GetConfigurations() map[string]Configuration
	AddStringConfigurations(ctx ServerContext, names []string, defaultValues []string)
	AddStringConfiguration(ctx ServerContext, name string)
	AddConfigurations(ctx ServerContext, requiredConfigTypeMap map[string]datatypes.DataType)
	AddOptionalConfigurations(ctx ServerContext, requiredConfigTypeMap map[string]datatypes.DataType, defaultValueMap utils.StringMap)
	GetConfiguration(ctx ServerContext, name string) (interface{}, bool)
	GetStringConfiguration(ctx ServerContext, name string) (string, bool)
	GetSecretConfiguration(ctx ServerContext, name string) ([]byte, bool, error)
	GetStringsMapConfiguration(ctx ServerContext, name string) (map[string]string, bool)
	GetStringArrayConfiguration(ctx ServerContext, name string) ([]string, bool)
	GetBoolConfiguration(ctx ServerContext, name string) (bool, bool)
	GetMapConfiguration(ctx ServerContext, name string) (config.Config, bool)
}
