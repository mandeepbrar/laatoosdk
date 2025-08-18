package rbac

import (
	"laatoo.io/sdk/server/auth"
	"laatoo.io/sdk/server/components/data"
)

type RbacUser interface {
	auth.User
	GetRoles() ([]data.StorableRef, error)
	SetRoles([]data.StorableRef) error
}
