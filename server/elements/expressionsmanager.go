package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

type ExpressionsManager interface {
	core.ServerElement
	components.ExpressionsManager
}
