package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

type NotificationManager interface {
	core.ServerElement
	components.NotifiersRegistry
	SendNotification(ctx core.RequestContext, notification *core.Notification) error
	Broadcast(ctx core.RequestContext, notif *core.Notification) error
	RegisterNotificationHandler(ctx core.ServerContext, notifType core.NotificationType, queue string, reg components.NotifiersRegistry) error
}
