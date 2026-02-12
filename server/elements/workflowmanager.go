package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

type WorkflowManager interface {
	core.ServerElement
	components.WorkflowManager
	RegisterProvider(ctx core.ServerContext, wfType components.WorkflowType, mgr components.WorkflowManager) error
}
