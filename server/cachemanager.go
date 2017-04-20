package server

import (
	"laatoo/sdk/components"
	"laatoo/sdk/core"
)

type CacheManager interface {
	core.ServerElement
	GetCache(ctx core.ServerContext, name string) components.CacheComponent
}
