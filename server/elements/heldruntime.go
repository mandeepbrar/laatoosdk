package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/components/ai"
	"laatoo.io/sdk/server/core"
)

// The runtime layer's addressable elements — the named things the activity, script, agent, cache
// and workflow managers hold.
//
// Same split as the data elements: each wraps an implementation and is not one. See helddata.go
// for why that separation is load-bearing rather than stylistic.

// Activity is the server's handle on a registered workflow step. The ActivityManager stores these
// under an "activity.<name>" alias, which the activity itself does not carry -- one reason the
// element is a separate handle rather than the activity itself.
type Activity interface {
	core.ServerElement

	// Activity returns the implementation behind this element.
	Activity() core.Activity

	// GetDefinition returns the activity's configuration definition.
	GetDefinition() *core.ActivityDefinition
}

// Script is the server's handle on a registered script. A script is resolved by alias by workflow
// activities and by the expression layer.
type Script interface {
	core.ServerElement

	// Script returns the implementation behind this element.
	Script() components.Script
}

// Agent is the server's handle on a registered AI agent.
type Agent interface {
	core.ServerElement

	// Agent returns the implementation behind this element.
	Agent() ai.Agent
}

// Skill is the server's handle on a registered skill.
//
// A skill is reachable by its canonical name AND by any alias registered alongside it, so two
// addresses may resolve to one element. That is deliberate and is why the element reports its own
// address: asking which name a reference bound to is answered by asking the element.
type Skill interface {
	core.ServerElement

	// Skill returns the implementation behind this element.
	Skill() ai.Skill

	// GetDescriptor returns the skill's descriptor -- metadata, tags and declared parameters.
	GetDescriptor() *ai.SkillDescriptor
}

// Cache is the server's handle on a built cache component. The CacheManager declares cache names
// during Initialize and builds the components behind them during Start, so this element exists
// only from Start onward.
type Cache interface {
	core.ServerElement

	// Cache returns the implementation behind this element.
	Cache() components.CacheComponent
}

// LLMProvider is the server's handle on a registered LLM provider. A provider registers under its
// own name and again under each model it lists, so several addresses may resolve to one element.
type LLMProvider interface {
	core.ServerElement

	// Provider returns the implementation behind this element.
	Provider() ai.LLMProvider
}

// Workflow is the server's handle on a registered workflow definition.
type Workflow interface {
	core.ServerElement

	// Workflow returns the definition behind this element.
	Workflow() components.Workflow
}
