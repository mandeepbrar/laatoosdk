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

// MemoryItem is one unit stored in a MemoryBank.
//
// In practice this interface is CLOSED: the shipped session bank type-asserts every item to
// the concrete *AIMemoryItem without checking (SessionMemory.Add,
// laatoosessionmemory/sessionmemoryservice.go:149), so passing your own MemoryItem
// implementation to a session bank PANICS rather than returning an error. Always construct
// &ai.AIMemoryItem{...}. The other banks are more forgiving but still lossy — ReferenceMemory
// and ChromemMemory copy selected fields into a fresh *AIMemoryItem and discard the rest,
// including the item's Id and Type.
type MemoryItem interface {
	core.Storable

	// GetContent returns the item's payload. Declared as any, but the only implementation
	// (AIMemoryItem) returns its []byte Content field — usually JSON. A caller doing
	// item.GetContent().(string) therefore always fails the assertion; use []byte.
	// Both vector banks re-marshal whatever comes back through the json codec before
	// embedding it, so a []byte payload is embedded as its quoted base64 JSON form.
	GetContent() any

	// GetImportance returns a 0..1 relevance weight used by the session bank's smart-lattice
	// pruning. A session item left at 0 is scored on save — 0.2/0.5/0.8 by content length —
	// so 0 means "unset", not "worthless". The vector banks store it as metadata only.
	GetImportance() float64

	// GetTimestamp returns an ISO8601/RFC3339 STRING, not a time.Time. It is the session
	// bank's sole ordering key (string-compared, which is why the format must be RFC3339), and
	// both SessionMemory.Add and .Set stamp time.Now() when it is empty. Nothing else does:
	// an item handed to a vector bank with an empty timestamp keeps an empty timestamp, and
	// ChromemMemory's Prune index then has no time to prune it by.
	GetTimestamp() string

	// GetTags returns free-form labels. Stored by the session and reference banks; ChromemMemory
	// flattens them into its string-only document metadata. Nothing in the platform filters or
	// retrieves by tag.
	GetTags() []string

	// GetMetadata returns arbitrary side-channel data carried with the item. The agent
	// manager writes {"role": "<AgentStakeholder>"} here for conversation turns.
	// ChromemMemory can only store string values, so non-string metadata is dropped on the
	// way into that bank.
	GetMetadata() utils.StringMap
}

// MemoryBank manages storage and retrieval of MemoryItems.
//
// FOUR IMPLEMENTATIONS SHIP AND THEY DO NOT AGREE ON THE CONTRACT. Obtain one through
// AgentManager.GetMemory/CreateMemory with a MemoryType, and read the per-method notes for the
// type you asked for:
//
//   - MemoryTypeSession (laatoosessionmemory) — cache-backed append log, keyed items in a
//     separate bucket. The only bank where Add/Retrieve behave as an append log.
//   - MemoryTypeShared (laatoosharedmemory) — cache bucket plus a pub/sub notification;
//     Add and Set write to the SAME bucket, so the "log vs keyed" split does not exist at all.
//   - MemoryTypeReferences (laatooreferencememory) — vector store over a DataComponent;
//     Set/Get are an in-PROCESS map that is lost on restart and invisible to other replicas.
//   - MemoryTypeData (chromemmemory) — chromem vector store; Set writes to BOTH the vector
//     index and an in-process map, so Retrieve does return Set items here.
//
// ONLY THE FIRST TWO ARE REACHABLE. Session and Shared call
// AgentManager.RegisterAgentMemoryManager in their Initialize; chromemmemory and
// laatooreferencememory implement AgentMemoryManager but never register it, so
// GetMemory/CreateMemory with MemoryTypeData or MemoryTypeReferences always fails NotFound
// ("Memory Manager") no matter which modules the solution loads. Reaching those two banks
// today means holding the service object directly.
//
// Method errors are unreliable across the board: several implementations swallow the
// underlying failure and report success or emptiness instead. Each case is called out below.
type MemoryBank interface {
	// GetId returns the bank's identifier — the id passed to CreateMemory/GetMemory. For a
	// session bank that is the sessionId, and it is what namespaces the underlying cache
	// bucket or collection.
	GetId() string

	// Add records an item in the bank. Use Set for keyed upsert semantics.
	//
	// "Append to a log" is only true for MemoryTypeSession, which stamps a UUID when the item
	// has no ID, stamps a timestamp when it has none, scores importance when it is 0, then
	// prunes the window and deletes what fell out. It PANICS on any MemoryItem that is not a
	// *ai.AIMemoryItem (unchecked assertion, sessionmemoryservice.go:149).
	//
	// The others diverge sharply and none of them generates an ID:
	//   - Shared keys the cache entry on item.GetId(), making Add a keyed UPSERT into the same
	//     bucket Set uses — an item with an empty Id overwrites the empty-key entry every
	//     time — and then publishes the CONTENT to topic "memory:update:<id>".
	//   - Reference and Data embed the content and store a vector document, DISCARDING the
	//     item's Id (Reference also forces Type to "reference"). Both call the json codec on
	//     GetContent() first, so an embedding failure — usually a missing/misconfigured
	//     embedder — surfaces here rather than at retrieval.
	Add(ctx core.RequestContext, item MemoryItem) error

	// Set stores an item under an explicit key, overwriting any existing item with the same key.
	// Intended for skill/agent state that must be read back by exact key (e.g. pending HITL state).
	//
	// Durability differs per bank and this is the trap:
	//   - Session writes to a separate cache KV bucket — survives restart, shared across replicas.
	//   - Shared writes to the ordinary cache bucket — shared, but the same one Add uses.
	//   - Reference stores it in a plain in-process map guarded by a mutex — NOT persisted, NOT
	//     visible to another pod, gone on restart (referencememorybank.go:192).
	//   - Data (chromem) stores it in an in-process sync.Map AND as a vector document with a
	//     deterministic "key:<key>" document ID.
	// Only Session returns an error for a non-*AIMemoryItem; the others accept any item.
	Set(ctx core.RequestContext, key string, item MemoryItem) error

	// Get retrieves a single item by its exact key, as previously stored via Set.
	// Returns nil, nil if no item exists for the key.
	//
	// nil, nil is also what you get when the read genuinely FAILED: SessionMemory.Get treats a
	// cache decode error as absence (sessionmemoryservice.go:257), so an `if err != nil` guard
	// passes and the caller nil-derefs the item. Always test the item for nil, not just the
	// error. ChromemMemory.Get reads only its in-process map, so after a process restart — or
	// on a second replica — a key written by Set on another instance reads back as absent even
	// though its vector document is still there. ChromemMemory.Clear does not clear that map,
	// so Get can also return an item the bank has otherwise deleted.
	Get(ctx core.RequestContext, key string) (MemoryItem, error)

	// Retrieve fetches items matching a query or opts filter.
	//
	// The claim that it does not return items stored via Set holds ONLY for the session and
	// reference banks. ChromemMemory.Set indexes into the same vector collection Retrieve
	// searches, and SharedMemory has one bucket for both — so on those two banks Set items do
	// come back.
	//
	// Per-bank semantics:
	//   - Session IGNORES both query and opts entirely and returns the whole window sorted by
	//     timestamp ascending; if the cache component cannot be obtained it returns (nil, nil)
	//     — an empty history and no error (sessionmemoryservice.go:406).
	//   - Reference and Data embed query as a vector and do similarity search, reading opts
	//     keys "limit" (int, default 5), "min_score" (float64, default 0), "use_mmr" (bool,
	//     default false) and "mmr_lambda" (float64, default 0.5). Any other opts key is
	//     ignored silently, and a "limit" supplied as a string is ignored rather than parsed.
	//   - Shared treats query as an exact cache KEY, not as a search, and returns a
	//     zero-length slice with no error when the key is absent.
	Retrieve(ctx core.RequestContext, query string, opts utils.StringMap) ([]MemoryItem, error)

	// Clear removes the bank's contents.
	//
	// Not uniformly complete or uniformly honest about failure. SessionMemory.Clear returns
	// nil when it cannot even reach the cache (sessionmemoryservice.go:432), so a nil error is
	// not evidence anything was deleted. ChromemMemory drops and recreates the collection —
	// wiping its gob files from disk — but leaves its in-process keyed map populated, so Get
	// still answers for keys written before the Clear. ReferenceMemory implements it as a
	// prune of everything older than now+100 years.
	Clear(ctx core.RequestContext) error

	// Synthesize is meant to compact or summarise the bank's contents.
	//
	// IT DOES NOTHING. All four shipped implementations are `return nil` with an empty body —
	// session, shared, reference and chromem alike. Calling it never compacts, never
	// summarises, and never reports that it did not: the nil error is indistinguishable from
	// success. Do not build a memory-growth strategy on this method; use Clear, or the session
	// bank's automatic window pruning inside Add.
	Synthesize(ctx core.RequestContext) error
} // End of MemoryBank interface
