package config

import "laatoo/sdk/server/ctx"

//Config Interface used by Laatoo
type Config interface {
	GetString(ctx ctx.Context, configurationName string) (string, bool)
	GetBool(ctx ctx.Context, configurationName string) (bool, bool)
	GetStringArray(ctx ctx.Context, configurationName string) ([]string, bool)
	GetSubConfig(ctx ctx.Context, configurationName string) (Config, bool)
	GetStringsMap(ctx ctx.Context, configurationName string) (map[string]string, bool)
	GetStringMap(ctx ctx.Context, configurationName string) (map[string]interface{}, bool)
	GetConfigArray(ctx ctx.Context, configurationName string) ([]Config, bool)
	GetRoot(ctx ctx.Context) (string, Config, bool)
	Get(ctx ctx.Context, configurationName string) (interface{}, bool)
	SetString(ctx ctx.Context, configurationName string, configurationValue string)
	Set(ctx ctx.Context, configurationName string, configurationValue interface{})
	SetVals(ctx ctx.Context, vals map[string]interface{})
	Clone() Config
	ToMap() map[string]interface{}
	AllConfigurations(ctx ctx.Context) []string
}
