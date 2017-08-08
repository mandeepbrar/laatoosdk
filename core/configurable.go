package core

import "laatoo/sdk/config"

type ConfigurableObject interface {
	GetConfigurations() map[string]interface{}
	AddStringConfigurations(ctx ServerContext, names []string, defaultValues []string)
	AddStringConfiguration(ctx ServerContext, name string)
	AddConfigurations(ServerContext, map[string]string)
	AddOptionalConfigurations(ServerContext, map[string]string, map[string]interface{})
	GetConfiguration(ServerContext, string) (interface{}, bool)
	GetStringConfiguration(ServerContext, string) (string, bool)
	GetBoolConfiguration(ServerContext, string) (bool, bool)
	GetMapConfiguration(ServerContext, string) (config.Config, bool)
}
