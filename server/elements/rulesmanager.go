package elements

import (
	"laatoo/sdk/server/components/rules"
	"laatoo/sdk/server/core"
)

type RulesManager interface {
	core.ServerElement
	SendSynchronousMessage(ctx core.RequestContext, msgType string, data interface{}) error
	SubscribeSynchronousMessage(ctx core.ServerContext, msgType string, rule rules.Rule)
}
