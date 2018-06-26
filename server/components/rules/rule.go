package rules

import (
	"laatoo/sdk/server/core"
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

type Rule interface {
	Condition(ctx core.RequestContext, trigger *Trigger) bool
	Action(ctx core.RequestContext, trigger *Trigger) error
}
