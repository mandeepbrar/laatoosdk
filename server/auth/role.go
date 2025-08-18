package auth

import "laatoo.io/sdk/datatypes"

type Role interface {
	datatypes.Serializable
	GetId() string
	SetId(string)
	GetName() string
	SetName(string)
	GetPermissions() []string
	SetPermissions([]string)
	GetTenant() TenantInfo
}
