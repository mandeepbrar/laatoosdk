package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

type WorkflowType string

const (
	WorkflowTypeFunction WorkflowType = "function"
	WorkflowTypeProcess  WorkflowType = "process"
	WorkflowTypeDurable  WorkflowType = "durable"
	WorkflowTypeAgent    WorkflowType = "agent"
)

type WorkflowManager interface {
	core.ServerElement
	components.WorkflowManager
	RegisterProvider(wfType WorkflowType, mgr components.WorkflowManager)
}
