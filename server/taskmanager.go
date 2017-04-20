package server

import (
	"laatoo/sdk/components"
	"laatoo/sdk/core"
)

type TaskManager interface {
	core.ServerElement
	PushTask(ctx core.RequestContext, queue string, task interface{}) error
	ProcessTask(ctx core.RequestContext, task *components.Task) error
}
