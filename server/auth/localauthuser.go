package auth

type LocalAuthUser interface {
	GetId() string
	SetId(string)
	LoadClaims(map[string]interface{})
	PopulateClaims(map[string]interface{})
	GetPassword() string
	ClearPassword()
	GetUsernameField() string
	GetUserName() string
	GetRealm() string
}
