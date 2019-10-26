package components

import "laatoo/sdk/server/ctx"

type Communicator interface {
	SendCommunication(ctx ctx.Context, communication map[interface{}]interface{}) error
}
