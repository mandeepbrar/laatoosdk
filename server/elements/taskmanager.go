package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type TaskManager interface {
	core.ServerElement
	PushTask(ctx core.RequestContext, queue string, task interface{}, metadata utils.StringMap) (string, error)
	SubscribeTaskCompletion(ctx core.ServerContext, topic string, handler core.MessageListener, subscriberId string) error
	CompleteTask(ctx core.RequestContext, queue string, invocationId string, result interface{}, metadata utils.StringMap, err error) error
	ProcessTask(ctx core.ServerContext, task *components.Task) (interface{}, error)
	List(ctx core.ServerContext) utils.StringsMap
	CreateEmptyTaskObj(ctx core.ServerContext) *components.Task
}
