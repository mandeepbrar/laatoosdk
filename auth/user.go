package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type User interface {
	GetId() string
	SetId(string)
	GetIdField() string
	GetUsernameField() string
	GetUserName() string
	SetJWTClaims(*jwt.Token)
	LoadJWTClaims(*jwt.Token)
}
