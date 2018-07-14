package elements

import (
	"laatoo/sdk/server/components"
	"laatoo/sdk/server/core"
)

type CacheManager interface {
	core.ServerElement
	GetCache(ctx core.ServerContext, name string) components.CacheComponent
}
