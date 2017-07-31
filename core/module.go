package core

import "laatoo/sdk/config"

type Module interface {
	Initialize(config.Config) error
	Factories() map[string]config.Config
	Services() map[string]config.Config
	Rules() map[string]config.Config
	Channels() map[string]config.Config
	Tasks() map[string]config.Config
}
