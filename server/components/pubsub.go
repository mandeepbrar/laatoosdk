package components

import (
	"laatoo.io/sdk/server/core"
)

type PubSubComponent interface {
	Publish(ctx core.RequestContext, topic string, message interface{}) error
	Subscribe(ctx core.ServerContext, topics []string, lstnr core.MessageListener) error
}
