package components

import (
	"laatoo.io/sdk/server/core"
)

type PubSubComponent interface {
	Publish(ctx core.RequestContext, topic string, message *core.Message) error
	Subscribe(ctx core.ServerContext, topics []string, lstnr core.MessageListener) error
}
