package components

import "laatoo.io/sdk/server/core"

type Notifier interface {
	GetSessionId() string
	GetUserId() string
	Notify(ctx core.RequestContext, notificaiton *core.Notification) error
}

type NotifiersRegistry interface {
	RegisterUserNotifier(ctx core.ServerContext, userId string, notifier Notifier)
}
