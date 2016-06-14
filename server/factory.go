package server

import (
	"laatoo/sdk/core"
)

type Factory interface {
	core.ServerElement
	Factory() core.ServiceFactory
}
