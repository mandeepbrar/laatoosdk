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
	// Topics returns the messaging topics this module declares.
	//
	// A topic belongs to the code that publishes and subscribes to it, not to whichever solution
	// happens to host that code. Before modules could declare their own, a topic existed only if a
	// solution author had named it in configuration — and publishing to a topic nobody named
	// silently dropped the message, so a plugin's events depended on someone else remembering them.
	Topics(ctx ServerContext) map[string]config.Config
	Workflows(ctx ServerContext) map[string]config.Config
	Activities(ctx ServerContext) map[string]config.Config
	GetContext() ServerContext
	ServerElement() ServerElement
	//	GetContext(ctx ServerContext, variable string) (interface{}, bool)
}
