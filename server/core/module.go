package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/utils"
)

type Module interface {
	ConfigurableObject
	Metadata() ModuleInfo
	MetaInfo(ctx ServerContext) utils.StringMap
	Describe(ServerContext) error
	Initialize(ctx ServerContext, conf config.Config) error
	Start(ctx ServerContext) error
	Factories(ctx ServerContext) map[string]config.Config
	Services(ctx ServerContext) map[string]config.Config
	Agents(ctx ServerContext) map[string]config.Config
	Rules(ctx ServerContext) map[string]config.Config
	Datasets(ctx ServerContext) map[string]config.Config
	Permissions(ctx ServerContext) utils.StringsMap
	Channels(ctx ServerContext) map[string]config.Config
	Tasks(ctx ServerContext) map[string]config.Config
	Workflows(ctx ServerContext) map[string]config.Config
	Activities(ctx ServerContext) map[string]config.Config
	GetContext() ServerContext
	ServerElement() ServerElement
	//	GetContext(ctx ServerContext, variable string) (interface{}, bool)
}
