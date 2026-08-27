package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

// NotificationManager routes notifications to the channel registered for their type. It is a
// dispatcher only — it holds no queue, no retry and no delivery state of its own; everything past
// the routing decision belongs to the components.NotificationChannel.
//
// Prefer ctx.SendNotification over reaching for this element directly.
type NotificationManager interface {
	core.ServerElement

	// SendNotification routes one notification to the channel registered for
	// notification.NotificationType and returns whatever that channel returns
	// (laatooserver/src/core/notificationmanager.go:70-81).
	//
	// A type with NO registered channel is a Bad Conf error, not a silent drop — but note the
	// registration itself can be lost silently, since a later RegisterNotificationHandler for the
	// same type replaces an earlier one without complaint.
	//
	// A nil return means the CHANNEL ACCEPTED the notification, not that it was delivered: the
	// email channel merely pushes onto a task queue and returns.
	SendNotification(ctx core.RequestContext, notification *core.Notification) error

	// Broadcast is NOT IMPLEMENTED. Every call returns errors.NotImplemented regardless of
	// arguments or configuration (notificationmanager.go:91-93). Do not build on it; send to
	// recipients explicitly instead.
	Broadcast(ctx core.RequestContext, notif *core.Notification) error

	// RegisterNotificationHandler binds a channel to one notification type. Services call it from
	// their own Initialize (see the email, in-app and SSE plugins).
	//
	// ONE HANDLER PER TYPE, LAST WRITER WINS, SILENTLY — the manager assigns into a plain map with
	// no duplicate check, so a second plugin registering for the same type displaces the first
	// with no error (notificationmanager.go:95-102).
	//
	// A nil reg is IGNORED and still returns nil: nothing is registered, and the caller has no way
	// to tell that from success.
	RegisterNotificationHandler(ctx core.ServerContext, notifType core.NotificationType, reg components.NotificationChannel) error
}
