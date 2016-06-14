package components

import (
	"laatoo/sdk/auth"
	"laatoo/sdk/core"
)

type AuthenticationComponent interface {
	SetTokenGenerator(core.ServerContext, func(auth.User) (string, auth.User, error))
}
