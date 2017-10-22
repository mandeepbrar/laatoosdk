package config

import "laatoo/sdk/ctx"

//Config Interface used by Laatoo
type Config interface {
	GetString(ctx ctx.Context, configurationName string) (string, bool)
	GetBool(ctx ctx.Context, configurationName string) (bool, bool)
	GetStringArray(ctx ctx.Context, configurationName string) ([]string, bool)
	GetSubConfig(ctx ctx.Context, configurationName string) (Config, bool)
	GetConfigArray(ctx ctx.Context, configurationName string) ([]Config, bool)
	Get(ctx ctx.Context, configurationName string) (interface{}, bool)
	SetString(ctx ctx.Context, configurationName string, configurationValue string)
	Set(ctx ctx.Context, configurationName string, configurationValue interface{})
	SetVals(ctx ctx.Context, vals map[string]interface{})
	Clone() Config
	AllConfigurations(ctx ctx.Context) []string
}
