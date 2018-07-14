package components

import (
	"laatoo/sdk/server/auth"
	"laatoo/sdk/server/core"
)

type AuthenticationComponent interface {
	SetTokenGenerator(core.ServerContext, func(auth.User, string) (string, auth.User, error))
}
