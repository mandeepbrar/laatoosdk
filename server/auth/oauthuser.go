package auth

type OAuthUser interface {
	User
	GetEmail() string
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
