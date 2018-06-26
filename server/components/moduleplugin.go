package components

import (
	"laatoo/sdk/common/config"
	"laatoo/sdk/server/core"
)

type ModuleManagerPlugin interface {
	Load(ctx core.ServerContext, name, moduleName, dir, parentMod string, mod core.Module, conf config.Config, settings config.Config, props map[string]interface{}) error
	Loaded(ctx core.ServerContext) error
}
