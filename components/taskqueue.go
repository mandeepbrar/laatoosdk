package components

import (
	"laatoo/sdk/core"
)

type TaskQueue interface {
	PushTask(ctx core.RequestContext, queue string, task interface{}) error
}
