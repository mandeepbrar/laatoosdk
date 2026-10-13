package components

import (
	"time"

	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/errors"
)

// Durable topics.
//
// A topic is the platform's one messaging name: publishers publish to it, subscribers subscribe to
// it, and it is declared in configuration rather than invented at a call site. Durability does not
// add a second kind of name alongside it — there is no separate "stream" a caller has to know
// about. Durability is a PROPERTY OF THE TOPIC, declared where the topic is declared:
//
//	topics:
//	  workflow.execution.started: {}          # transient, unchanged
//	  kvdatabase.wal:
//	    durable: true
//	    storage: file
//	    replicas: 3
//	    maxage: 168h
//
// What that buys the topic: messages survive a restart, carry a sequence, are redelivered if a
// subscriber does not acknowledge them, and can be replayed from a past position. What it does not
// change: how anyone publishes to it or subscribes to it. Publish and Subscribe are the same calls
// on both kinds of topic, which is the point — a subscriber cannot opt into a durability its
// publisher did not provide, so making it a subscriber-side choice would only invite that mistake.
//
// The underlying broker's own vocabulary — JetStream streams, consumers, subjects — stays inside
// the provider. A plugin consuming a durable topic never names one.

// TopicStorage selects where a durable topic's messages are persisted.
//
// APPEND NEW VALUES AT THE END, NEVER INSERT. Plugins are compiled separately against their own
// pinned SDK, so renumbering an existing constant changes its meaning inside already-built plugins
// with no build error anywhere — the failure appears at runtime as a topic provisioned with the
// wrong storage class, far from the edit that caused it.
type TopicStorage int

const (
	// FileStorage persists to the provider's store directory, so the topic survives a restart.
	// This is the only storage class for which durability means anything, and it is the default.
	FileStorage TopicStorage = iota
	// MemoryStorage keeps messages in memory only. The topic is empty after a restart — useful for
	// redelivery and ordering without the disk cost, not for a write-ahead log.
	MemoryStorage
)

// TopicRetention selects when a durable topic drops a message.
//
// Append-only, for the same reason as TopicStorage.
type TopicRetention int

const (
	// LimitsRetention keeps every message until a configured limit evicts it, so a subscriber can
	// resume at any sequence still held and replay forward. This is what a write-ahead log needs,
	// and it is the default because replay is the reason to make a topic durable at all.
	LimitsRetention TopicRetention = iota
	// WorkQueueRetention drops a message once it has been acknowledged, so each message is handled
	// by exactly one subscriber and replay is not possible. For work distribution rather than for
	// an event log.
	WorkQueueRetention
)

// TopicDelivery selects how a durable topic's messages are spread across the subscribers that share
// a subscriber id -- which, in a multi-replica deployment, means across pods.
//
// Append-only, for the same reason as TopicStorage.
type TopicDelivery int

const (
	// BroadcastDelivery gives every subscriber its own position, so each one receives every
	// message. This matches what a transient topic does, and it is the default so that declaring a
	// topic durable changes its guarantees without quietly changing its fan-out.
	//
	// It requires each subscriber to be distinguishable across replicas; the platform derives that
	// from the server's identity, so a deployment whose replicas do not have stable identities
	// gets a fresh position -- and therefore a full replay -- whenever a replica is replaced.
	BroadcastDelivery TopicDelivery = iota
	// DistributedDelivery shares one position between every subscriber using the same subscriber
	// id, so each message is handled by exactly one of them. For work that must happen once no
	// matter how many replicas are running, and wrong for anything each replica needs to see --
	// a replica applying a change log needs every message, not a share of them.
	DistributedDelivery
)

// TopicDurability is a durable topic's declared configuration, parsed from the topic's config
// block. The zero value is a valid durable topic: file storage, limits retention, one replica, and
// the provider's own defaults for everything else.
//
// Every limit's zero value means "no limit from this field". A durable topic with no limit at all
// grows until the store fills, so a topic taking unbounded input is expected to set at least one.
type TopicDurability struct {
	// Storage selects persistence. Defaults to FileStorage, the zero value.
	Storage TopicStorage
	// Retention selects eviction. Defaults to LimitsRetention, the zero value.
	Retention TopicRetention
	// Delivery selects fan-out across replicas. Defaults to BroadcastDelivery, the zero value,
	// which matches transient behaviour.
	Delivery TopicDelivery
	// Replicas requested for the topic. 0 and 1 both mean a single replica. A value the cluster
	// cannot satisfy fails at startup rather than being quietly reduced to 1 — a topic that claims
	// replication it does not have is worse than one that admits it has none.
	Replicas int
	// MaxAge evicts a message once it is older than this. Zero means no age limit.
	MaxAge time.Duration
	// MaxMsgs caps the number of messages retained. Zero means no count limit.
	MaxMsgs int64
	// MaxBytes caps the total bytes retained. Zero means no size limit.
	MaxBytes int64
	// MaxDeliver bounds redelivery of a message a subscriber keeps declining. On exhaustion the
	// platform dead-letters it rather than redelivering forever. Zero takes the provider default.
	MaxDeliver int
	// AckWait is how long the platform waits for a subscriber to return before assuming it failed
	// and redelivering. Zero takes the provider default.
	//
	// Setting this very high disables redelivery for genuine crashes, so a subscriber whose work
	// takes an unbounded amount of time should return promptly and track completion elsewhere
	// rather than hold a message unacknowledged for hours.
	AckWait time.Duration
}

// DurablePubSubComponent is the OPTIONAL durable capability of a pub/sub provider.
//
// Deliberately separate from PubSubComponent rather than folded into it. PubSubComponent has
// implementors outside the server — redispubsub asserts it at compile time — and widening it would
// break every one of them to demand a capability MOST OF THEM HAVE NO WAY TO PROVIDE. A provider
// that can offer durability implements this as well; the messaging manager type asserts for it when
// a topic is declared durable, and a declaration the provider cannot honour fails at startup.
//
// THE EMPHASIS IS LOAD-BEARING, and the release that added PubSubComponent.Unsubscribe is what
// makes it so. That change widened the base interface, which reads as a contradiction of this
// paragraph and is not one: the test is not "does widening break implementors" — it always does —
// but "is this a capability an implementor could reasonably lack". Durability is: a transport with
// no persistence cannot invent it, and demanding it would exclude providers the platform wants.
// Detaching a subscription is not: anything that can attach one can stop delivering to it, so a
// provider that cannot is broken rather than merely simpler. That is why unsubscribe went into the
// base interface and durability stays out here.
//
// This interface is the PROVIDER contract, not the caller-facing one. Callers use
// elements.MessagingManager, whose Publish and Subscribe are the same for both kinds of topic.
type DurablePubSubComponent interface {
	PubSubComponent

	// EnsureDurableTopic provisions topic to match cfg, or reconciles an existing one to it. Called
	// once per durable topic at startup, so it must be idempotent — and a limit changed in config
	// takes effect on the next start rather than needing anything torn down by hand.
	EnsureDurableTopic(ctx core.ServerContext, topic string, cfg *TopicDurability) error

	// PublishDurable appends message to topic and stamps message.Sequence with the position it was
	// stored at. It fails when the message was not persisted — the acknowledgement IS the
	// durability guarantee, so a publish that was not stored is an error rather than the silent
	// success the transient path allows.
	//
	// A non-empty message.Id makes the append idempotent within the provider's deduplication
	// window: republishing the same id does not append a second copy, and Sequence is stamped with
	// the original's position.
	PublishDurable(ctx core.RequestContext, topic string, message *core.Message) error

	// SubscribeDurable attaches the subscriber named lsnrid to topic. The name is what lets the
	// subscriber resume rather than replay after a restart, so it must be stable across restarts.
	//
	// fromSequence positions the subscriber: 0 resumes from its recorded position, or from the
	// oldest retained message if it has never run; any other value starts there, so 1 means from
	// the beginning of what the topic still holds.
	SubscribeDurable(ctx core.ServerContext, topic string, lstnr core.MessageListener, lsnrid string, fromSequence uint64) error

	// UnsubscribeDurable detaches lsnrid from topic, leaving its recorded position intact so a
	// later subscribe resumes rather than replays.
	UnsubscribeDurable(ctx core.ServerContext, topic string, lsnrid string) error
}

// DurableNotSupported builds the error returned when a topic is declared durable but the configured
// messaging provider has no durable capability.
//
// Ordinarily this surfaces at STARTUP, when the manager provisions declared topics — a deployment
// that asks for durability the provider cannot give should be told before it takes traffic, not on
// the first publish. It also covers SubscribeFrom against a transient topic, which is a caller
// error rather than a deployment one.
//
// The code it carries, errors.CORE_ERROR_DURABLE_NOT_SUPPORTED, is registered in the errors package
// alongside every other platform code — registration is what gives it a message and what
// RegisterErrorHandler keys against, and an unregistered code panics if it ever reaches the
// standard-error path.
//
// CALLERS MUST STILL CHECK FOR NIL. errors.ThrowError returns nil when a handler registered against
// the internal code claims to have handled it (errorregistry.go:120-131), so a deployment that
// registers a handler for this code turns this into a nil error. Nothing registers one today, and
// this is not a reason to avoid the error registry — but a nil error here must not be read as
// success.
func DurableNotSupported(c ctx.Context, topic string) error {
	// ThrowError stamps the internal error code that IsDurableNotSupported reads back
	return errors.ThrowError(c,
		"topic is declared durable but the configured messaging provider has no durable capability: "+topic,
		errors.CORE_ERROR_DURABLE_NOT_SUPPORTED)
}

// IsDurableNotSupported reports whether err means durability was asked for and is unavailable.
//
// This is the detection half of the contract: a caller that wants to degrade rather than fail asks
// this rather than inspecting the message. A direct type assertion is sufficient because
// errors.WrapError returns a *errors.Error unchanged rather than nesting it, so the internal code
// survives being wrapped on the way up.
func IsDurableNotSupported(err error) bool {
	if err == nil {
		return false
	}
	laatooErr, ok := err.(*errors.Error)
	if !ok {
		return false
	}
	return laatooErr.InternalErrorCode == errors.CORE_ERROR_DURABLE_NOT_SUPPORTED
}
