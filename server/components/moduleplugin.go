package components

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type ModInfo struct {
	ModName            string
	PluginName         string
	PluginDir          string
	ParentModName      string
	Mod                core.Module
	UserObj            core.Module
	PluginConf         config.Config
	ModConf            config.Config
	ModSettings        config.Config
	Configurations     map[string]core.Configuration
	ModProps           utils.StringMap
	IsExtended         bool
	ExtendedPluginName string
	ExtendedPluginConf config.Config
	ExtendedPluginDir  string
	Hot                bool
}

func (info *ModInfo) GetContext(ctx core.ServerContext, variable string) (interface{}, bool) {
	return info.Mod.GetContext().Get(variable)
}

type ModuleManagerPlugin interface {
	GetName() string
	Load(ctx core.ServerContext, modInfo *ModInfo) error
	Loaded(ctx core.ServerContext) error
	Unloaded(ctx core.ServerContext, insName, modName string) error
	Unloading(ctx core.ServerContext, insName, modName string) error
}
