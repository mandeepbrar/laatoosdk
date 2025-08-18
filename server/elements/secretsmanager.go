package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

type SecretsManager interface {
	core.ServerElement
	components.SecretsManager
}
