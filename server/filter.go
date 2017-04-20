package server

import (
	"laatoo/sdk/core"
)

type Filter interface {
	Allowed(ctx core.ServerContext, objectName string) bool
}
