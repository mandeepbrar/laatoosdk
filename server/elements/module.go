package elements

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type Module interface {
	core.ServerElement
	GetObject() core.Module
	GetModuleProperties() utils.StringMap
}
