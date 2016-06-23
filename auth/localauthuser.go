package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type LocalAuthUser interface {
	GetId() string
	SetId(string)
	GetIdField() string
	GetPassword() string
	SetPassword(string)
	SetJWTClaims(*jwt.Token)
	LoadJWTClaims(*jwt.Token)
	GetUsernameField() string
	GetUserName() string
	GetRealm() string
}
