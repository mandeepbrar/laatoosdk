package elements

import (
	"laatoo.io/sdk/server/components"
	"laatoo.io/sdk/server/core"
)

// CacheManager is the server element that owns the deployment's configured caches. Reach it with
// ctx.GetServerElement(core.ServerElementCacheManager).
type CacheManager interface {
	core.ServerElement

	// GetCache returns the cache registered under name, or NIL when nothing is registered under
	// that name.
	//
	// THERE IS NO ERROR RETURN AND NO FALLBACK. The lookup is a plain map read over this manager's
	// own registrations (laatooserver/src/core/cachemanagerproxy.go:12-18): it does not walk up to
	// a parent level and does not substitute a default. A misspelt or unconfigured name therefore
	// yields a nil interface that panics at the first method call, far from the lookup that caused
	// it. Nil-check the result.
	//
	// Which implementation comes back depends on what the deployment configured under that name,
	// and their behaviours differ materially — including whether the cache is shared between
	// replicas at all. See components.CacheComponent.
	GetCache(ctx core.ServerContext, name string) components.CacheComponent
}
