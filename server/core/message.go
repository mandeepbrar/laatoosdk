package core

import (
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/utils"
)

type MessageListener func(ctx RequestContext, message *Message, info utils.StringMap) error

type Message struct {
	Data   interface{}
	Tenant auth.TenantInfo
	User   auth.User
}
