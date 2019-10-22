package components

import "laatoo/sdk/server/ctx"

type Communicator interface {
	SendMessage(ctx ctx.Context, recipients []string, msg interface{}) error
}
