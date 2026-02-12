package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type ScriptManager interface {
	core.ServerElement
	RegisterScript(ctx core.ServerContext, alias string, act components.Script) error
	RegisterProvider(ctx core.ServerContext, extension string, provider components.ScriptManager) error
}
