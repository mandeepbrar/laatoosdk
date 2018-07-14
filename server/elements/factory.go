package elements

import (
	"laatoo/sdk/server/core"
)

type Factory interface {
	core.ServerElement
	Factory() core.ServiceFactory
}
