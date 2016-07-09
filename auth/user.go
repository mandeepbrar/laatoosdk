package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type User interface {
	GetId() string
	GetIdField() string
	GetUsernameField() string
	GetUserName() string
	LoadJWTClaims(*jwt.Token)
	PopulateJWTToken(*jwt.Token)
	GetRealm() string
}
