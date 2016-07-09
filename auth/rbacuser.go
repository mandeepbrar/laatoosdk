package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type RbacUser interface {
	GetId() string
	GetIdField() string
	LoadJWTClaims(*jwt.Token)
	PopulateJWTToken(*jwt.Token)
	GetUsernameField() string
	GetUserName() string
	GetPermissions() (permissions []string, err error)
	GetRoles() ([]string, error)
	GetRealm() string
	SetPermissions(permissions []string)
}
