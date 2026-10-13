package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

// The messaging layer's addressable element.
//
// Same membership rule as helddata.go and heldruntime.go: a topic is resolved BY NAME, so it is an
// element and carries an address. What differs is that there is nothing behind it to wrap.

// Topic is the server's handle on a declared pub/sub topic.
//
// A topic is DECLARED, not built: a plugin's registry/topics/<topic>.yml or a namespace's own
// configuration states that the topic exists and, optionally, how it is persisted. There is no
// provider object behind it — the messaging provider is per namespace and serves every topic —
// which is why this interface carries no X() accessor where DataComponent has Component() and
// Skill has Skill(). The declaration is the whole of the thing.
//
// Its ADDRESS is what the element exists for. Topic identity is inherited: a namespace may publish
// and subscribe to a topic an enclosing namespace declared, and two namespaces may declare the same
// name with different durability, the nearer shadowing the further. Asking which declaration a
// reference bound to is answered by asking the element where it lives.
type Topic interface {
	core.ServerElement

	// Durability returns how this topic is persisted, or nil when it is transient.
	//
	// NIL IS THE ORDINARY CASE and is not an error: most topics are transient. A caller that treats
	// nil as "unknown" rather than "transient" will provision storage nobody declared.
	Durability() *components.TopicDurability
}
