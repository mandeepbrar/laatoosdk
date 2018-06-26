package elements

import (
	"laatoo/sdk/server/components"
	"laatoo/sdk/server/core"
)

type TaskManager interface {
	core.ServerElement
	PushTask(ctx core.RequestContext, queue string, task interface{}) error
	ProcessTask(ctx core.RequestContext, task *components.Task) error
}
