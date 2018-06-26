package auth

type OAuthUser interface {
	GetId() string
	SetId(string)
	LoadClaims(map[string]interface{})
	PopulateClaims(map[string]interface{})
	GetEmail() string
	GetUsernameField() string
	GetUserName() string
	GetRealm() string
}
