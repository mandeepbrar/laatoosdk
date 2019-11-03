package elements

import (
	"laatoo/sdk/server/components"
	"laatoo/sdk/server/core"
)

type SecretsManager interface {
	core.ServerElement
	components.SecretsManager
}
