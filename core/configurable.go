package core

import "laatoo/sdk/config"

type Configuration interface {
	GetName() string
	IsRequired() bool
	GetDefaultValue() interface{}
	GetValue() interface{}
	GetType() string
}

type ConfigurableObject interface {
	GetConfigurations() map[string]Configuration
	AddStringConfigurations(ctx ServerContext, names []string, defaultValues []string)
	AddStringConfiguration(ctx ServerContext, name string)
	AddConfigurations(ctx ServerContext, requiredConfigTypeMap map[string]string)
	AddOptionalConfigurations(ctx ServerContext, requiredConfigTypeMap map[string]string, defaultValueMap map[string]interface{})
	GetConfiguration(ctx ServerContext, name string) (interface{}, bool)
	GetStringConfiguration(ctx ServerContext, name string) (string, bool)
	GetBoolConfiguration(ctx ServerContext, name string) (bool, bool)
	GetMapConfiguration(ctx ServerContext, name string) (config.Config, bool)
}
