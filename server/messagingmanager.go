package server

import (
	"laatoo/sdk/core"
)

type MessagingManager interface {
	core.ServerElement
	Publish(ctx core.RequestContext, topic string, message interface{}) error
	Subscribe(ctx core.ServerContext, topics []string, lstnr core.ServiceFunc) error
}
