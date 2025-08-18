package elements

import (
	"laatoo.io/sdk/server/components/rules"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

type RulesManager interface {
	core.ServerElement
	SendSynchronousMessage(ctx core.RequestContext, msgType string, data interface{}) error
	SubscribeSynchronousMessage(ctx core.ServerContext, msgType string, rule rules.Rule, ruleName string)
	List(ctx core.ServerContext) utils.StringsMap
	Describe(ctx core.ServerContext, rule string) (utils.StringMap, error)
}
