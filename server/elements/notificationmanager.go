package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

// NotificationManager routes notifications to NAMED CHANNELS. It is a dispatcher only — it holds no
// queue, no retry and no delivery state of its own; everything past the routing decision belongs to
// the components.NotificationChannel.
//
// A channel is registered under a name and serves one or more notification types, so a namespace may
// hold several channels of one type — two email outboxes, say, "supportmail" and "billingmail". Which
// one a notification reaches is decided by SendNotification below; Broadcast reaches them all.
//
// Prefer ctx.SendNotification over reaching for this element directly.
type NotificationManager interface {
	core.ServerElement

	// SendNotification delivers one notification to ONE channel serving notification.NotificationType
	// and returns whatever that channel returns. The channel is:
	//
	//   1. notification.Channel, when set — exactly that channel, or an error if no channel of that
	//      name serves the type;
	//   2. otherwise the type's DEFAULT channel, declared in the namespace's notifications
	//      configuration;
	//   3. otherwise the ONLY channel serving the type, when there is exactly one;
	//   4. otherwise the send is REFUSED, naming every candidate — never delivered to whichever
	//      registered first.
	//
	// Channels resolve nearest-first: a namespace sees its own channels and those of every enclosing
	// namespace, and a channel it declares under an enclosing channel's name shadows it. A type with no
	// channel anywhere in reach is an error, not a silent drop.
	//
	// A nil return means the CHANNEL ACCEPTED the notification, not that it was delivered: the email
	// channel merely pushes onto a task queue and returns.
	SendNotification(ctx core.RequestContext, notification *core.Notification) error

	// Broadcast delivers one notification to EVERY channel serving notification.NotificationType that
	// the caller can reach, ignoring notification.Channel — the shape a server event wants, where each
	// live stream is its own channel and every one of them must receive it. Every channel is attempted
	// even when one fails; the failures are returned together.
	Broadcast(ctx core.RequestContext, notif *core.Notification) error

	// RegisterNotificationChannel registers a channel under a name, serving the given notification
	// types. Services call it from their own Initialize (see the email, in-app and SSE plugins).
	//
	// THE NAME IDENTIFIES THE CHANNEL, not the type: two channels may serve core.EMAIL as long as their
	// names differ, and a notification picks one with Notification.Channel. An EMPTY name registers
	// the channel under the registering service's own name, which is already unique in its namespace —
	// the right default for a plugin instantiated once per outbox.
	//
	// A DUPLICATE NAME IN ONE NAMESPACE IS REFUSED, naming both registrations; a name an enclosing
	// namespace holds is an override, not a duplicate. An empty types list, and a nil channel, are
	// refused rather than accepted as a registration of nothing.
	RegisterNotificationChannel(ctx core.ServerContext, name string, types []core.NotificationType, channel components.NotificationChannel) error
}
