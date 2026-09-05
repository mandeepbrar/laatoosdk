package components

import (
	"laatoo.io/sdk/server/core"
)

// PubSubComponent is the transport a messaging provider supplies to the platform's messaging
// manager. Plugin code does NOT call it directly — it publishes through ctx.PublishMessage /
// the MessagingManager element, which applies the topic-declaration gate and the durable/transient
// split before reaching this interface.
//
// THIS INTERFACE CHANGED SHAPE in the release carrying Unsubscribe: Subscribe gained a
// subscriberId and Unsubscribe was added. Both halves are a COMPILE-TIME break for implementors,
// and that is the point rather than a cost. Adding Unsubscribe alone would have been an additive
// change that broke implementors SILENTLY -- Go checks interface satisfaction at the assertion
// site, so a provider would compile and then fail at plugin.Open, far from the cause. Changing
// Subscribe in the same release makes every implementor fail at build, where the fix is obvious.
//
// TOPIC DECLARATION IS ENFORCED ONE LAYER ABOVE THIS. The messaging manager refuses to publish to
// a topic no module or configuration declared, returning MESSAGING_TOPIC_NOT_DECLARED
// (laatooserver/src/core/messagingmanager.go:236-253). A topic is declared either in a plugin's
// src/server/registry/topics/<topic>.yml or in solution/application configuration. That refusal
// replaced a silent nil return that told publishers their message had been delivered while the
// platform discarded it, so an implementation of this interface must never re-introduce that
// behaviour by swallowing an unroutable publish.
type PubSubComponent interface {
	// Publish sends one message on a topic.
	//
	// The message carries the PUBLISHER'S IDENTITY and the transport is expected to carry it, so a
	// subscriber can act as the publisher rather than as an anonymous system principal. The
	// embedded-NATS provider serialises User (id, name, roles) and Tenant onto the wire
	// (messagingmanager.go:501-531); the Redis provider stamps ctx.GetUser()/ctx.GetTenant() onto
	// the message when the caller left them unset (redispubsubservice.go:44-49).
	//
	// Delivery guarantees are the provider's: the NATS transient path is fire-and-forget core
	// NATS, so a nil return means "handed to the broker", not "a subscriber received it".
	Publish(ctx core.RequestContext, topic string, message *core.Message) error

	// Subscribe attaches ONE listener to a whole set of topics at once and is called once, at
	// Start, with every transient topic the level knows about. The listener the messaging manager
	// passes is a demultiplexer: it reads the topic back from the request context key
	// "messagetype" and fans out to the per-topic listeners registered with the manager
	// (messagingmanager.go:368-380). An implementation MUST therefore set that key.
	//
	// THE SUBSCRIBER RUNS AS THE PUBLISHER, not as a system principal. The provider rebuilds the
	// published identity and creates the request with it (messagingmanager.go:696-720), so an
	// authorization check or a multitenant read inside a listener sees the publishing user and
	// tenant. A message carrying NO identity still delivers as an ordinary system request, with
	// authorization DISABLED — but a message that asserted an identity which could not be rebuilt
	// is REFUSED and dropped, never downgraded to system.
	//
	// subscriberId NAMES THE SUBSCRIPTION so it can later be detached, and it is the parameter this
	// interface previously lacked. The platform passes the subscribing namespace's own address, so
	// a root and a child that both subscribe to one topic name -- reachable whenever a child
	// listens to a topic an ancestor declared -- are two distinct subscriptions the provider can
	// tell apart. A provider MUST key its subscription state by it: without that, Unsubscribe
	// below cannot say which of two subscriptions on the same topic to drop.
	//
	// It must be stable for the lifetime of the subscription. Re-subscribing with an id that is
	// already attached is the provider's call to define, but it must not leave two live
	// subscriptions the caller can only detach one of.
	//
	// Listener errors are DISCARDED — the manager invokes each listener on its own goroutine and
	// returns nil (messagingmanager.go:375-377), so a failing subscriber is invisible to the
	// publisher and to the transport.
	Subscribe(ctx core.ServerContext, topics []string, lstnr core.MessageListener, subscriberId string) error

	// Unsubscribe detaches everything subscriberId attached through Subscribe, so the provider stops
	// delivering to it.
	//
	// BY SUBSCRIBER RATHER THAN BY TOPIC, deliberately. The caller tearing a subscription down is a
	// namespace shutting down, and what it wants is "drop everything I attached" -- not a list of
	// topics it would have to reconstruct, and which would be ambiguous anyway when another
	// namespace holds its own subscription on the same topic.
	//
	// An unknown subscriberId is NOT an error. Teardown runs on paths that may not have subscribed
	// -- a namespace that declared no topics, a boot that failed before Start -- and a provider that
	// refused there would turn a clean shutdown into a reported failure. A provider that cannot
	// detach at all should say so with an error rather than returning nil, because a silent no-op
	// here is a subscription leaking for the life of the process with nothing reporting it.
	//
	// The messaging manager's own unsubscribe is a DIFFERENT operation and remains: it removes one
	// listener from one namespace's per-topic list. This detaches the provider-level subscription
	// underneath it.
	Unsubscribe(ctx core.ServerContext, subscriberId string) error
}
