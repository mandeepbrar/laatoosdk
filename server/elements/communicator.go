package elements

import (
	"laatoo/sdk/server/components"
	"laatoo/sdk/server/core"
)

type Communicator interface {
	core.ServerElement
	components.Communicator
}
