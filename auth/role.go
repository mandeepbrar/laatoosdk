package auth

type Role interface {
	GetId() string
	SetId(string)
	GetPermissions() []string
	SetPermissions([]string)
	GetRealm() string
	SetRealm(val string)
	GetName() string
	SetName(string)
}
