package elements

import (
	"laatoo.io/sdk/server/core"
)

// Solution is the level between the server and its applications: one deployed product, owning the
// applications, the tenants, and -- when clustering is configured -- the embedded NATS server that
// backs pub/sub, tasks, the cache and the cross-pod service registry.
type Solution interface {
	core.ServerElement

	// GetPeers returns the other pods of this solution, as node names.
	//
	// THE filter ARGUMENT IS IGNORED. No implementation reads it; the full peer set is returned
	// whatever is passed.
	//
	// The result is live only when the embedded NATS server is running, in which case it is that
	// server's currently active peers. Otherwise it is the snapshot taken once when the solution
	// started -- which is nil on a solution with no embedded NATS at all, i.e. any single-pod
	// deployment. The error is never non-nil.
	GetPeers(filter string) ([]string, error)
}
