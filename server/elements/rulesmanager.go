package elements

import (
	"laatoo.io/sdk/server/components/rules"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// RulesManager owns the rules declared at one server level and dispatches the SYNCHRONOUS half of
// the rule system. Rules declared with trigger: async are not dispatched here at all — they are
// subscribed as ordinary pub/sub topic listeners and run on their own goroutines
// (laatooserver/src/core/rulesmanager.go:550-570).
//
// Rules are loaded from src/server/registry/rules/<name>.yml during Initialize. A rule name must be
// unique across the level: a second rule file with a name already taken fails startup with a Bad
// Conf error rather than replacing the first (rulesmanager.go:520-526).
type RulesManager interface {
	core.ServerElement

	// SendSynchronousMessage runs every rule subscribed to msgType, INLINE on the caller's
	// goroutine, in this order: Condition, then Action if the condition held
	// (rulesmanager.go:592-610). This is the path that lets an entity presave rule reject a save.
	//
	// A MESSAGE TYPE NO RULE SUBSCRIBED TO IS A SILENT NO-OP returning nil. There is no
	// declaration check here — unlike topic publishing, which now refuses an undeclared topic — so
	// a typo in msgType is indistinguishable from "no rule matched".
	//
	// The first Action error ABORTS the loop and is returned; rules that had not yet run do not
	// run. Rules for one message type execute in Go map iteration order, so which rules ran before
	// the failing one differs between invocations.
	//
	// data is wrapped into a rules.Trigger with TriggerType SynchronousMessage and handed to each
	// rule as Trigger.Message.
	SendSynchronousMessage(ctx core.RequestContext, msgType string, data interface{}) error

	// SubscribeSynchronousMessage registers a rule to receive msgType through
	// SendSynchronousMessage. Normally the rules manager calls this itself for every rule file
	// declaring trigger: sync; call it directly only to register a rule built in Go.
	//
	// RETURNS NOTHING, AND SILENTLY OVERWRITES. Rules are held per message type keyed by ruleName,
	// so re-registering an existing name for the same message type replaces the previous rule with
	// no error and no warning (rulesmanager.go:583-590). ruleName is also the handle unsubscribe
	// uses, so two distinct rules sharing a name cannot both be reached.
	SubscribeSynchronousMessage(ctx core.ServerContext, msgType string, rule rules.Rule, ruleName string)

	// List returns the rules loaded at this level.
	//
	// THE VALUES ARE ALWAYS THE EMPTY STRING. The map is built as ruleName -> "" with the module
	// name left as an unimplemented TODO (rulesmanager.go:508-514), so only the keys carry
	// information — do not read the value as a module or plugin name.
	List(ctx core.ServerContext) utils.StringsMap

	// Describe returns {"Name", "Conf"} for one rule — its name and the raw configuration it was
	// created from — or a NotFound error for an unknown rule (rulesmanager.go:645-653). It does
	// not report the rule's trigger type or subscribers separately; those are inside Conf.
	Describe(ctx core.ServerContext, rule string) (utils.StringMap, error)
}
