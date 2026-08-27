package components

import "laatoo.io/sdk/server/core"

/*
type Notifier interface {
	GetSessionId() string
	GetUserId() string
	Notify(ctx core.RequestContext, notificaiton *core.Notification) error
}*/

// NotificationChannel is the delivery mechanism for ONE core.NotificationType — email, in-app,
// server-sent event, webhook. A service implements it and registers itself against a type with
// elements.NotificationManager.RegisterNotificationHandler, normally from its own Initialize.
//
// ONE HANDLER PER NOTIFICATION TYPE, LAST REGISTRATION WINS, SILENTLY. The manager stores handlers
// in a plain map keyed by type with no duplicate check
// (laatooserver/src/core/notificationmanager.go:95-102), so two plugins registering for core.EMAIL
// leave only the one that started later — with no error and no warning.
type NotificationChannel interface {
	// SendNotification delivers one notification. It is reached only through
	// ctx.SendNotification / NotificationManager.SendNotification, which routes on
	// notification.NotificationType; a type with no registered channel is a Bad Conf error rather
	// than a silent drop (notificationmanager.go:70-81).
	//
	// Delivery may be asynchronous and usually is: the email channel does not send anything here,
	// it pushes the notification onto its configured task queue and returns
	// (laatoomodules/notifications/dev/plugins/email/src/server/go/emailservice.go:139-142). A nil
	// return therefore means "accepted", not "delivered".
	SendNotification(ctx core.RequestContext, notification *core.Notification) error
}
