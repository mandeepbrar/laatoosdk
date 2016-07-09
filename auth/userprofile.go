package auth

type UserProfile interface {
	GetId() string
	GetIdField() string
	GetEmail() string
	GetName() string
	GetPicture() string
	GetGender() string
	GetUsernameField() string
	GetUserName() string
}
