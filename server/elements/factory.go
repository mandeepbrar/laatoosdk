package elements

import (
	"laatoo.io/sdk/server/core"
)

type Factory interface {
	core.ServerElement
	Factory() core.ServiceFactory
	GetModule() core.Module
	Metadata() core.ServiceFactoryInfo
}
