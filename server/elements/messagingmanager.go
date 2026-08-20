package elements

import "laatoo.io/sdk/server/core"

// MessagingManager is the platform's messaging surface.
//
// There is ONE kind of name here: a topic, declared in configuration. Publishers publish to a
// topic, subscribers subscribe to a topic, and that is true whether or not the topic is durable —
// durability is declared on the topic itself (see components.TopicDurability), not chosen at a call
// site. So Publish and Subscribe below are the whole surface for both kinds of topic, and a plugin
// moving a topic from transient to durable changes configuration, not code.
//
// What a durable topic adds, without changing these signatures:
//   - the message survives a restart, and Publish stamps message.Sequence with where it was stored
//   - message.Id, if set, makes the publish idempotent within the topic's deduplication window
//   - the listener's return value acknowledges: nil advances past the message, an error redelivers
//     it up to the topic's configured bound, after which it is dead-lettered
//   - lsnrid names the subscriber durably, so a restart resumes rather than replays
//
// The one thing durability enables that has no transient counterpart is replaying from a past
// position, which is what SubscribeFrom is for.
type MessagingManager interface {
	core.ServerElement

	// Publish sends message to topic. On a durable topic it returns only once the message is
	// persisted — a publish that was not stored is an error rather than the silent success the
	// transient path allows — and message.Sequence is stamped with its position.
	Publish(ctx core.RequestContext, topic string, message *core.Message) error

	// Subscribe attaches lstnr to topics under the subscriber name lsnrid.
	//
	// On a durable topic, lsnrid is the subscriber's durable identity: it must be stable across
	// restarts, because it is what lets the platform resume delivery from where this subscriber
	// stopped instead of replaying everything. A subscriber that has never run starts from the
	// oldest message the topic still retains.
	Subscribe(ctx core.ServerContext, topics []string, lstnr core.MessageListener, lsnrid string) error

	// SubscribeFrom is Subscribe positioned at a past point in a durable topic, for a subscriber
	// that knows where it left off and cannot rely on the recorded position — a replica catching up
	// from the last change it applied, a consumer rebuilding state from an event log.
	//
	// fromSequence is a value previously seen as message.Sequence; 1 means the oldest message the
	// topic still retains, and 0 behaves as Subscribe. Delivery is in sequence order from there.
	//
	// ONE topic, not a list, unlike Subscribe: a sequence is a position within a single topic, and
	// the same number means something unrelated in each of several topics. A subscriber catching up
	// on more than one topic is catching up to a different point in each, so it calls this once per
	// topic with that topic's own position.
	//
	// Returns an error for which components.IsDurableNotSupported reports true when the topic is
	// not durable — a transient topic has no position to start from, and silently ignoring the
	// argument would strand a subscriber that believes it is catching up.
	SubscribeFrom(ctx core.ServerContext, topic string, lstnr core.MessageListener, lsnrid string, fromSequence uint64) error
}
