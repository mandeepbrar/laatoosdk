package server

import (
	"laatoo/sdk/core"
)

type TaskManager interface {
	core.ServerElement
	PushTask(ctx core.RequestContext, queue string, task interface{}) error
}
