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
	// GetTask returns a previously pushed task by the id PushTask returned. The queue is
	// required because it selects which configured backend owns the task. Returns a
	// NotImplemented error when the backend serving that queue cannot address tasks by id.
	GetTask(ctx core.RequestContext, queue string, id string) (*components.Task, error)
	List(ctx core.ServerContext) utils.StringsMap
	CreateEmptyTaskObj(ctx core.ServerContext) *components.Task
}
