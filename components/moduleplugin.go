package components

import (
	"laatoo/sdk/config"
	"laatoo/sdk/core"
)

type ModuleManagerPlugin interface {
	Load(ctx core.ServerContext, name, dir string, mod core.Module, conf config.Config, settings config.Config) error
	Loaded(ctx core.ServerContext) error
}
