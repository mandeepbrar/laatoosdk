package components

import "laatoo/sdk/server/core"

type Communicator interface {
	SendCommunication(ctx core.RequestContext, communication map[interface{}]interface{}) error
}
