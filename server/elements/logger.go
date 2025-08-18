package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

type Logger interface {
	core.ServerElement
	components.Logger
}
