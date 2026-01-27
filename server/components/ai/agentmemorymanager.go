package ai

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// AgentMemoryManager manages the lifecycle and access to named MemoryBanks.
// It acts as a central registry and facade for memory operations.
type AgentMemoryManager interface {

	// CreateMemory initializes a new named memory bank of a specific type.
	CreateMemory(ctx core.RequestContext, id string, config map[string]interface{}) (MemoryBank, error)

	// GetMemory retrieves an existing memory bank by name.
	GetMemory(ctx core.RequestContext, id string) (MemoryBank, error)

	// DeleteMemory removes a memory bank.
	DeleteMemory(ctx core.RequestContext, id string) error

	// AddItem adds a memory item to the specified memory bank.
	AddItem(ctx core.RequestContext, memoryid string, item MemoryItem) error

	// Retrieve searches for items in the specified memory bank.
	// Uses the core.RetrievalConfig which allows specifying limits, scoring, filters, etc.
	Retrieve(ctx core.RequestContext, memoryid string, query string, config utils.StringMap) ([]MemoryItem, error)
}
