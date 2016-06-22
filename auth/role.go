package auth

type Role interface {
	GetId() string
	SetId(string)
	GetIdField() string
	GetPermissions() []string
	SetPermissions([]string)
	GetName() string
	SetName(string)
}
