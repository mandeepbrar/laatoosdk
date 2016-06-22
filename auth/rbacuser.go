package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type RbacUser interface {
	GetId() string
	SetId(string)
	GetIdField() string
	GetUsernameField() string
	GetUserName() string
	GetPermissions() (permissions []string, err error)
	SetPermissions(permissions []string)
	GetRoles() ([]string, error)
	SetRoles(roles []string) error
	AddRole(role string) error
	RemoveRole(role string) error
	SetJWTClaims(*jwt.Token)
	LoadJWTClaims(*jwt.Token)
}
