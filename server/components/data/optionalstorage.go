package data

import (
	"time"

	"laatoo.io/sdk/server/core"
)

// Optional storage capabilities.
//
// These are capabilities only SOME providers can offer, declared the way this package already
// declares one: an optional interface next to DataComponent rather than a method on it. The
// reasoning is written out in full at ExpandingComponent (query.go) and applies unchanged here.
//
// Never on DataComponent itself. Adding a method to that interface breaks every implementor
// SILENTLY, because Go checks satisfaction at the assertion site — the providers still compile
// and fail at startup instead. Around 130 plugin go.mod files pin this module, so the blast
// radius of getting that wrong is the whole ecosystem.
//
// Deliberately NOT a QueryCapability and NOT a Feature, either:
//
//   - QueryCapability names something a provider can COMPILE inside a query. Expiry and prefix
//     enumeration are neither predicates nor projections; they are storage operations, and
//     query.go says so where it separates the two types.
//   - Feature (data.go) is a self-declared claim, which is exactly the shape ExpandingComponent
//     rejected: a provider can answer true and then ignore the request, and nothing downstream
//     can tell. An interface makes "claims support but ignores it" unrepresentable instead of
//     merely discouraged.
//
// Nor do these carry the granular sub-capabilities expansion has. That granularity exists
// because one $expand carries options a provider may support only partly — filter, select,
// orderby, paging, count — and a single flag would let it accept a query it answers wrongly.
// Expiry and prefix enumeration have no sub-options to be partly right about, so a second layer
// would be ceremony rather than safety.
//
// THE REFUSAL PROTOCOL IS PART OF EACH CONTRACT. Implementing one of these interfaces means the
// provider can serve the operation in general, not that it can serve every call. A provider that
// cannot serve a particular call returns errors.NotImplemented, and callers must distinguish
// that from a genuine failure — the same split the expansion caller makes between "this provider
// declined, use the fallback" and "this read failed". A provider that cannot serve the operation
// at all must not implement the interface, so the caller never offers it the work.

// ExpiringComponent is implemented by providers that can expire a record after a duration.
//
// A record past its expiry must read as absent — by id AND by query — from the moment it
// expires, not from whenever a background sweep happens to reclaim it. Providers differ in when
// they physically reclaim the storage and that difference must not be observable.
//
// Candidate implementors beyond the embedded key-value engines: aerospike (per-record TTL),
// couchbase (document expiry) and mongo (TTL index) all have a native mechanism. None implements
// this today, which is the argument for the interface being optional rather than required.
//
// Two limits a caller cannot see from the signature and which an implementor must therefore
// state in its own documentation, because both are real today:
//
//   - Granularity is the provider's, not the caller's. A store whose expiry resolution is one
//     second cannot honour a 100ms ttl and should refuse rather than round.
//   - Replication may make expiry unsafe. A record expiring independently on each replica
//     diverges them, so a replicated provider may refuse TTL entirely — which is a refusal at
//     call time, not a reason to withhold the interface.
type ExpiringComponent interface {
	// PutWithTTL stores item at id, creating or replacing, and expires it after ttl.
	// A ttl of zero or less stores the record with no expiry.
	PutWithTTL(ctx core.RequestContext, id string, item core.Storable, ttl time.Duration) error
	// SaveWithTTL stores item under its own id and expires it after ttl.
	SaveWithTTL(ctx core.RequestContext, item core.Storable, ttl time.Duration) error
}

// PrefixScanComponent is implemented by providers whose store can enumerate keys by prefix.
//
// This is a KEY-space operation, not a query: it answers "which records are filed under this
// prefix" without decoding any of them, which is why it returns ids rather than records and why
// it is not expressible as a filter. A provider whose keys carry no meaningful ordering — or
// whose only prefix answer would be a full scan with a string comparison per record — should not
// implement this, because the whole point of the interface is that the caller can rely on it
// being cheaper than reading everything.
//
// limit bounds the result; zero or less means unbounded, which a caller should reserve for a
// prefix it knows to be small. An unbounded prefix scan over a large store is a full scan.
type PrefixScanComponent interface {
	// GetPrefix returns the ids of records whose key begins with prefix, up to limit.
	GetPrefix(ctx core.RequestContext, prefix string, limit int) ([]string, error)
	// DeletePrefix removes every record whose key begins with prefix, returning how many were
	// removed. Implementors must remove secondary index entries with the records; a prefix
	// delete that leaves an index pointing at nothing is the failure this method most invites.
	DeletePrefix(ctx core.RequestContext, prefix string) (int, error)
}
