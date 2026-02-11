package elements

import "laatoo.io/sdk/server/core"

type MessagingManager interface {
	core.ServerElement
	Publish(ctx core.RequestContext, topic string, message *core.Message) error

	Subscribe(ctx core.ServerContext, topics []string, lstnr core.MessageListener, lsnrid string) error
}
