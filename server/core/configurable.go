package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/datatypes"
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
	GetDescription() string
	GetVersion() string
	GetConfigurations() map[string]Configuration
	//AddStringConfigurations(ctx ServerContext, names []string, defaultValues []string)
	AddStringConfiguration(ctx ServerContext, name string, desc string, defaultValue string)
	AddConfiguration(ctx ServerContext, name string, desc string, dtype datatypes.DataType, defaultValue interface{})
	AddOptionalConfiguration(ctx ServerContext, name string, desc string, dtype datatypes.DataType, defaultValue interface{})
	GetConfiguration(ctx ServerContext, name string) (interface{}, bool)
	GetStringConfiguration(ctx ServerContext, name string) (string, bool)
	GetSecretConfiguration(ctx ServerContext, name string) ([]byte, bool, error)
	GetStringsMapConfiguration(ctx ServerContext, name string) (map[string]string, bool)
	GetStringArrayConfiguration(ctx ServerContext, name string) ([]string, bool)
	GetBoolConfiguration(ctx ServerContext, name string) (bool, bool)
	GetMapConfiguration(ctx ServerContext, name string) (config.Config, bool)
}
