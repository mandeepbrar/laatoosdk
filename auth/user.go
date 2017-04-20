package auth

type User interface {
	GetId() string
	SetId(string)
	GetUsernameField() string
	GetUserName() string
	LoadClaims(map[string]interface{})
	PopulateClaims(map[string]interface{})
	GetRealm() string
}
