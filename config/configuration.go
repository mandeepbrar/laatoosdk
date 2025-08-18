package config

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/utils"
)

// Config Interface used by Laatoo
type Config interface {
	GetString(ctx ctx.Context, configurationName string) (string, bool)
	GetBool(ctx ctx.Context, configurationName string) (bool, bool)
	GetStringArray(ctx ctx.Context, configurationName string) ([]string, bool)
	GetSubConfig(ctx ctx.Context, configurationName string) (Config, bool)
	GetStringsMap(ctx ctx.Context, configurationName string) (utils.StringsMap, bool)
	GetStringMap(ctx ctx.Context, configurationName string) (utils.StringMap, bool)
	GetConfigArray(ctx ctx.Context, configurationName string) ([]Config, bool)
	GetRoot(ctx ctx.Context) (string, Config, bool)
	Get(ctx ctx.Context, configurationName string) (interface{}, bool)
	SetString(ctx ctx.Context, configurationName string, configurationValue string)
	Set(ctx ctx.Context, configurationName string, configurationValue interface{})
	SetVals(ctx ctx.Context, vals utils.StringMap)
	Clone() Config
	ToMap() map[string]interface{}
	AllConfigurations(ctx ctx.Context) []string
}
