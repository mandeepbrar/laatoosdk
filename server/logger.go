package server

import (
	"laatoo/sdk/components"
	"laatoo/sdk/core"
)

type Logger interface {
	core.ServerElement
	components.Logger
}
