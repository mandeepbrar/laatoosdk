package auth

import "laatoo.io/sdk/datatypes"

type UserAccount interface {
	datatypes.Serializable
	GetId() string
	GetEmail() string
	GetFullName() string
	GetPicture() string
	GetGender() string
	GetUsernameField() string
	GetUserName() string
}
