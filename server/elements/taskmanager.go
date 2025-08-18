package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type TaskManager interface {
	core.ServerElement
	PushTask(ctx core.RequestContext, queue string, task interface{}) error
	SubscribeTaskCompletion(queue string, callback func(ctx core.RequestContext, invocationId string, result interface{})) error
	ProcessTask(ctx core.ServerContext, task *components.Task) (interface{}, error)
	List(ctx core.ServerContext) utils.StringsMap
	CreateEmptyTaskObj(ctx core.ServerContext) *components.Task
}
