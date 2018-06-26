package elements

import "laatoo/sdk/server/core"

type Module interface {
	core.ServerElement
	GetObject() core.Module
	GetModuleProperties() map[string]interface{}
}
