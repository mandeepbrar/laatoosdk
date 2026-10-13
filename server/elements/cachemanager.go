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
	// THE LOOKUP WALKS OUTWARD, nearest namespace first (laatooserver/src/core/cachemanagerproxy.go:
	// 14-23): a cache configured by an enclosing namespace is usable from every namespace beneath it,
	// and a namespace configuring the same name itself shadows the enclosing one.
	//
	// THERE IS NO ERROR RETURN AND NO DEFAULT. A name configured nowhere along that walk -- misspelt,
	// or never configured -- yields a nil interface that panics at the first method call, far from
	// the lookup that caused it. Nil-check the result.
	//
	// Corrected 2026-09-12: this comment said the lookup reads only this manager's own registrations
	// and does not walk to a parent level. The proxy walks.
	//
	// Which implementation comes back depends on what the deployment configured under that name,
	// and their behaviours differ materially — including whether the cache is shared between
	// replicas at all. See components.CacheComponent.
	GetCache(ctx core.ServerContext, name string) components.CacheComponent
}
