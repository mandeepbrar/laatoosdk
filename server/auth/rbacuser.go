package auth

type RbacUser interface {
	User
	GetPermissions() (permissions []string, err error)
	GetRoles() ([]string, error)
	SetRoles([]string) error
	SetPermissions(permissions []string)
}

/*GetId() string
SetId(string)
GetUsernameField() string
GetUserName() string
LoadClaims(map[string]interface{})
PopulateClaims(map[string]interface{})
GetRealm() string
GetTenant() string
*/
