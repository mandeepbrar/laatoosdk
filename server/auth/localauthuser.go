package auth

type LocalAuthUser interface {
	User
	GetPassword() string
	ClearPassword()
}
