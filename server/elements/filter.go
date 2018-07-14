package elements

import (
	"laatoo/sdk/server/core"
)

type Filter interface {
	Allowed(ctx core.ServerContext, objectName string) bool
}
