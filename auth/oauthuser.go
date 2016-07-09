package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type OAuthUser interface {
	GetId() string
	GetIdField() string
	LoadJWTClaims(*jwt.Token)
	PopulateJWTToken(*jwt.Token)
	GetEmail() string
	GetUsernameField() string
	GetUserName() string
	GetRealm() string
}
