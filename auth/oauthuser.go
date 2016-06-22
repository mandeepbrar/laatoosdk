package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type OAuthUser interface {
	GetId() string
	SetId(string)
	GetIdField() string
	SetJWTClaims(*jwt.Token)
	LoadJWTClaims(*jwt.Token)
	GetEmail() string
	GetUsernameField() string
	GetUserName() string
}
