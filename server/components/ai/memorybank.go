package ai

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ------------------------------
// Memory System Interfaces
// ------------------------------

type MemoryType string

const (
	MemoryTypeSession    MemoryType = "Session"
	MemoryTypeShared     MemoryType = "Shared"
	MemoryTypeReferences MemoryType = "References"
	MemoryTypeData       MemoryType = "Data"
)

type MemoryItem interface {
	core.Storable
	GetContent() any
	GetImportance() float64
	GetTimestamp() string
	GetTags() []string
	GetMetadata() utils.StringMap
}

// MemoryBank manages storage and retrieval of MemoryItems.
type MemoryBank interface {
	GetId() string

	// Add appends an item to the memory log. If the item has no ID, one is generated.
	// Use Set for keyed upsert semantics.
	Add(ctx core.RequestContext, item MemoryItem) error

	// Set stores an item under an explicit key, overwriting any existing item with the same key.
	// Intended for skill/agent state that must be read back by exact key (e.g. pending HITL state).
	Set(ctx core.RequestContext, key string, item MemoryItem) error

	// Get retrieves a single item by its exact key, as previously stored via Set.
	// Returns nil, nil if no item exists for the key.
	Get(ctx core.RequestContext, key string) (MemoryItem, error)

	// Retrieve fetches items from the append log based on a query or opts filter.
	// Does not return items stored via Set.
	Retrieve(ctx core.RequestContext, query string, opts utils.StringMap) ([]MemoryItem, error)

	Clear(ctx core.RequestContext) error
	Synthesize(ctx core.RequestContext) error
} // End of MemoryBank interface

