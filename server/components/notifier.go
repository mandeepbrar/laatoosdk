package components

import "laatoo/sdk/server/core"

type Notifier interface {
	Notify(ctx core.ServerContext, identifier interface{}, msg interface{}, dataid string) error
	Broadcast(ctx core.ServerContext, msg interface{}, dataid string) error
}
