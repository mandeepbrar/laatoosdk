package core

import (
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/utils"
)

// MessageListener receives a message published to a topic the listener subscribed to.
//
// On a DURABLE topic the return value is the acknowledgement: nil acknowledges the message and the
// subscriber advances past it, while an error declines it and the platform redelivers, up to the
// topic's configured redelivery bound, after which the message is dead-lettered. A listener that
// must not see a message twice therefore does its work before returning nil, not after.
//
// On a transient topic the return value is logged and nothing else — there is nothing to
// acknowledge, because there is nothing holding the message.
type MessageListener func(ctx RequestContext, message *Message, info utils.StringMap) error

type Message struct {
	// Id is the publisher's identifier for this message. Optional, and ignored entirely on a
	// transient topic.
	//
	// On a durable topic it is the IDEMPOTENCY KEY: republishing the same id within the topic's
	// deduplication window does not append a second copy, and the publish reports the sequence of
	// the original. A publisher that can retry — which, on a durable topic, is every publisher that
	// cares about the message — sets this to something stable across its own retries.
	Id string

	Data   interface{}
	Tenant auth.TenantInfo
	User   auth.User

	// Sequence is the message's position in a durable topic, assigned by the platform. It is
	// STAMPED, not supplied: Publish writes it on return, and delivery to a listener carries the
	// sequence the message was stored at.
	//
	// Zero on a transient topic, which has no ordering to report. A subscriber records this to
	// resume from where it stopped — see MessagingManager.SubscribeFrom.
	Sequence uint64
}
