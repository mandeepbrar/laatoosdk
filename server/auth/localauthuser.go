package auth

type LocalAuthUser interface {
	User
	GetPassword() string
	ClearPassword()
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
