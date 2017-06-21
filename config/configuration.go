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
