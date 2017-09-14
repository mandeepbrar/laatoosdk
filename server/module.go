package server

import "laatoo/sdk/core"

type Module interface {
	core.ServerElement
	GetObject() core.Module
	GetModuleProperties() map[string]interface{}
}
