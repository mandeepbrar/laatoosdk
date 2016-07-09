package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type LocalAuthUser interface {
	GetId() string
	GetIdField() string
	LoadJWTClaims(*jwt.Token)
	PopulateJWTToken(*jwt.Token)
	GetPassword() string
	ClearPassword()
	GetUsernameField() string
	GetUserName() string
	GetRealm() string
}
