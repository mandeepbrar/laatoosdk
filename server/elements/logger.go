package elements

import (
	"laatoo/sdk/server/components"
	"laatoo/sdk/server/core"
)

type Logger interface {
	core.ServerElement
	components.Logger
}
