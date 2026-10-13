package components

import "laatoo.io/sdk/server/core"

/*
type Notifier interface {
	GetSessionId() string
	GetUserId() string
	Notify(ctx core.RequestContext, notificaiton *core.Notification) error
}*/

// NotificationChannel is a delivery mechanism — an email outbox, in-app delivery, a server-sent event
// stream, a webhook. A service implements it and registers itself under a NAME, serving one or more
// core.NotificationTypes, with elements.NotificationManager.RegisterNotificationChannel, normally from
// its own Initialize.
//
// Several channels may serve one type — two email outboxes are two channels with different names —
// and a notification picks one with core.Notification.Channel, or reaches the type's default. A
// duplicate NAME in one namespace is refused; see RegisterNotificationChannel.
type NotificationChannel interface {
	// SendNotification delivers one notification. It is reached only through
	// ctx.SendNotification / NotificationManager.SendNotification, which picks ONE channel for the
	// notification, or NotificationManager.Broadcast, which reaches every channel serving its type.
	//
	// Delivery may be asynchronous and usually is: the email channel does not send anything here,
	// it pushes the notification onto its configured task queue and returns
	// (laatoomodules/notifications/dev/plugins/email/src/server/go/emailservice.go:139-142). A nil
	// return therefore means "accepted", not "delivered".
	SendNotification(ctx core.RequestContext, notification *core.Notification) error
}
