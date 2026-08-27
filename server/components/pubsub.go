package components

import (
	"laatoo.io/sdk/server/core"
)

// PubSubComponent is the transport a messaging provider supplies to the platform's messaging
// manager. Plugin code does NOT call it directly — it publishes through ctx.PublishMessage /
// the MessagingManager element, which applies the topic-declaration gate and the durable/transient
// split before reaching this interface.
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
	// There is no subscriber id and no Unsubscribe in this interface: a listener attached here
	// cannot be detached through it. The messaging manager's unsubscribe only removes the listener
	// from its own per-topic list; the provider-level subscription stays open
	// (messagingmanager.go:174-189).
	//
	// Listener errors are DISCARDED — the manager invokes each listener on its own goroutine and
	// returns nil (messagingmanager.go:375-377), so a failing subscriber is invisible to the
	// publisher and to the transport.
	Subscribe(ctx core.ServerContext, topics []string, lstnr core.MessageListener) error
}
