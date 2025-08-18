package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

type CacheManager interface {
	core.ServerElement
	GetCache(ctx core.ServerContext, name string) components.CacheComponent
}
