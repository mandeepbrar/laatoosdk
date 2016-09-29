package auth

type RbacUser interface {
	GetId() string
	SetId(string)
	LoadClaims(map[string]interface{})
	PopulateClaims(map[string]interface{})
	GetUsernameField() string
	GetUserName() string
	GetPermissions() (permissions []string, err error)
	GetRoles() ([]string, error)
	GetRealm() string
	SetPermissions(permissions []string)
}
