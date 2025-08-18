package elements

import (
	"laatoo.io/sdk/server/core"
)

type Filter interface {
	Allowed(ctx core.ServerContext, objectName string) bool
}
