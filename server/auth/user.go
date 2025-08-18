package auth

import "laatoo.io/sdk/datatypes"

type TenantInfo interface {
	datatypes.Serializable
	GetTenantId() string
	GetTenantName() string
}

type User interface {
	datatypes.Serializable
	GetId() string
	SetId(string)
	GetUsernameField() string
	GetUserName() string
	LoadClaims(map[string]interface{})
	PopulateClaims(map[string]interface{})
	GetEmail() string
	GetRealm() string
	GetTenant() TenantInfo
	GetUserAccount() UserAccount
}
