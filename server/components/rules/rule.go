package rules

import (
	"laatoo.io/sdk/server/core"
)

type TriggerType int

const (
	AsynchronousMessage TriggerType = iota
	SynchronousMessage
)

type Trigger struct {
	TriggerType TriggerType
	MessageType string
	Message     interface{}
}

// Rule is the platform's event-driven reaction unit — the mechanism behind entity lifecycle
// triggers (storable_presave, storable_postsave, ...) and any other message type a plugin declares
// a rule for. A rule is declared in src/server/registry/rules/<name>.yml naming a trigger type
// (async or sync), a message type and the Go object implementing this interface; the rules manager
// creates the object, calls Initialize on it, and subscribes it
// (laatooserver/src/core/rulesmanager.go:520-581).
//
// Condition is always evaluated first and Action runs only if it returned true.
type Rule interface {
	// Condition decides whether this trigger is one the rule cares about. It returns only a bool —
	// there is no way to report an error from here, so a rule that cannot evaluate its condition
	// must choose between false (skip, silently) and true (let Action fail loudly).
	Condition(ctx core.RequestContext, trigger *Trigger) bool

	// Action performs the rule's effect.
	//
	// WHETHER THE RETURNED ERROR IS SEEN BY ANYONE DEPENDS ENTIRELY ON THE TRIGGER TYPE, and the
	// rule cannot tell from inside:
	//
	//   - AsynchronousMessage (trigger: async): Condition and Action run on a NEW GOROUTINE, and
	//     the topic listener returns nil immediately. The error is logged and discarded; the
	//     publisher never learns the rule failed, and the rule cannot veto anything
	//     (rulesmanager.go:555-569).
	//   - SynchronousMessage (trigger: sync): Action runs inline on the caller's goroutine and the
	//     error PROPAGATES to whatever sent the message — which is how a presave rule rejects a
	//     save. It also ABORTS the loop, so rules for the same message type that had not yet run
	//     do not run at all (rulesmanager.go:592-610).
	//
	// Multiple sync rules on one message type execute in Go map iteration order, i.e. a different
	// order on every run — never rely on one rule seeing another's effect.
	Action(ctx core.RequestContext, trigger *Trigger) error
}
