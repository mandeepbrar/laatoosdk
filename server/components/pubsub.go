package components

import (
	"laatoo/sdk/server/core"
)

type PubSubComponent interface {
	Publish(ctx core.RequestContext, topic string, message interface{}) error
	Subscribe(ctx core.ServerContext, topics []string, lstnr core.MessageListener) error
}
