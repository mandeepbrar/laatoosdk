package core

import "laatoo/sdk/config"

type Module interface {
	ConfigurableObject
	Initialize(ctx ServerContext) error
	Start(ctx ServerContext) error
	Factories(ctx ServerContext) map[string]config.Config
	Services(ctx ServerContext) map[string]config.Config
	Rules(ctx ServerContext) map[string]config.Config
	Channels(ctx ServerContext) map[string]config.Config
	Tasks(ctx ServerContext) map[string]config.Config
}
