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
	GetImportance() float64
	GetTimestamp() string
	GetTags() []string
	GetMetadata() utils.StringMap
}

// MemoryBank manages storage and retrieval of MemoryItems.
type MemoryBank interface {
	GetId() string
	// Add stores an item.
	Add(ctx core.RequestContext, item MemoryItem) error
	
	// Retrieve fetches items based on a query.
	Retrieve(ctx core.RequestContext, query string, opts utils.StringMap) ([]MemoryItem, error)
	
	Clear(ctx core.RequestContext) error
	Synthesize(ctx core.RequestContext) error
} // End of MemoryBank interface

