package components

import "laatoo/sdk/core"

type Notifier interface {
	Notify(ctx core.RequestContext, identifier interface{}, msg interface{}) error
	Broadcast(ctx core.RequestContext, msg interface{}) error
}
