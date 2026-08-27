package elements

import (
	"laatoo.io/sdk/server/components/ai"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// ============================================================
// HANDOFF MANAGER - Coordinates agent handoffs across the system
// ============================================================

// HandoffManager coordinates agent handoffs across the system.
//
// One instance is created and started by the agent manager during its Initialize; obtain it
// with AgentManager.GetHandoffManager(). Agents are enrolled automatically — every agent that
// satisfies ai.HandoffCapableAgent is registered as it is built — so a plugin rarely calls
// RegisterAgent itself. Its errors are plain fmt.Errorf values, not platform errors, so
// errors.Is/Core_* code matching does not work on them.
type HandoffManager interface {
	// RegisterAgent enrolls an agent and indexes it under each of its GetCapabilities strings.
	//
	// The registry is keyed on agent.GetName(), and registering a second agent under an
	// existing name OVERWRITES it with no error — while the capability index keeps BOTH
	// entries, so the displaced agent's capabilities still resolve, now to its replacement.
	// Always returns nil.
	RegisterAgent(ctx core.ServerContext, agent ai.HandoffCapableAgent) error

	// UnregisterAgent removes an agent and unwinds its capability-index entries, erroring if
	// no agent is registered under agentID. It re-reads GetCapabilities to know what to
	// unwind, so an agent whose capabilities changed since registration leaves stale index
	// entries pointing at a name that is no longer in the registry — and
	// FindAgentByCapability then returns a nil agent with a nil error for those capabilities.
	UnregisterAgent(ctx core.ServerContext, agentID string) error

	// GetAgent returns the agent registered under agentID, or an error if none is.
	GetAgent(ctx core.ServerContext, agentID string) (ai.HandoffCapableAgent, error)

	// ExecuteHandoff routes and runs a handoff, returning its result.
	//
	// When req.TargetAgentID is empty the target is resolved by the router and WRITTEN BACK
	// into req — this method mutates the request you pass it, and also appends it to the
	// per-session history when req.Context["sessionId"] is a string. Result and error are
	// independent: a non-nil result can accompany a non-nil error, and the result carries the
	// authoritative Success flag, so check both.
	ExecuteHandoff(ctx core.RequestContext, req *ai.HandoffRequest) (*ai.HandoffResult, error)

	// FindAgentByCapability returns the single agent matching the most of the requested
	// capabilities.
	//
	// It is a best-match, not a filter: an agent matching one of five requested capabilities
	// is returned as happily as one matching all five, and the caller cannot tell which
	// happened. Matching is exact and case-sensitive. TIES ARE BROKEN BY GO MAP ITERATION
	// ORDER, so with two equally-qualified agents the winner differs between calls in the same
	// process. An empty capabilities slice is an error; no match is an error. It can also
	// return (nil, nil) when the capability index holds a stale name — see UnregisterAgent.
	FindAgentByCapability(ctx core.ServerContext, capabilities []string) (ai.HandoffCapableAgent, error)

	// FindAgentByID returns the agent registered under agentID. It is a straight alias for
	// GetAgent with identical behaviour, not a broader search.
	FindAgentByID(ctx core.ServerContext, agentID string) (ai.HandoffCapableAgent, error)

	// Start starts the handoff executor, the agent event listener and the background
	// expired-handoff cleanup goroutine. Idempotent — a second call while running returns nil
	// without restarting anything. Called for you by the agent manager at startup.
	Start(ctx core.ServerContext) error

	// Stop DOES NOT STOP ANYTHING. The implementation only sets an internal isRunning flag to
	// false and returns nil (laatooserver/src/ai/handoffmanager.go:69): the executor and the
	// event listener are never stopped, the CleanupExpiredHandoffs goroutine keeps running,
	// and ExecuteHandoff continues to work afterwards because nothing consults the flag. The
	// only observable effect is that GetStatistics reports is_running=false and a later Start
	// will re-run the startup sequence.
	Stop(ctx core.ServerContext) error

	// GetStatistics returns a snapshot with exactly three keys: "registered_agents" (int),
	// "history_count" (int — the number of sessions with recorded handoffs, not the number of
	// handoffs) and "is_running" (bool, and see Stop for why it may lie). There are no
	// success, failure or latency counters.
	GetStatistics(ctx core.ServerContext) map[string]interface{}
}

// ============================================================
// AGENT MANAGER - Manages AI agents and their capabilities
// ============================================================

// AgentManager is the server element that owns AI agents, skills, LLM providers, agent memory
// and the MCP server registry. Reach it with
// ctx.GetServerElement(core.ServerElementAgentManager).(elements.AgentManager).
//
// Most of it is a set of registries populated at startup by module loading, plus a thin façade
// over the registered LLM providers. Register* methods are startup-time and take a
// ServerContext; the request-time surface is LLMRequest, LLMStreamingRequest, InvokeSkill and
// the memory methods.
type AgentManager interface {
	core.ServerElement

	// GetAgent returns the agent registered under alias, or a NotFound error.
	//
	// The alias is the agent YAML's metadata.id when it declares one, otherwise the config
	// filename — the two differ often enough to be the usual cause of a NotFound here.
	// See laatooserver/src/core/agentmanager.go:142.
	GetAgent(ctx core.ServerContext, alias string) (ai.Agent, error)
	//	GetEngine(ctx core.ServerContext, name string) (AgentEngine, error)

	// ListAgents returns every registered agent keyed by alias.
	//
	// It returns the manager's LIVE internal map, NOT a copy (agentmanager.go:256) — writing
	// to or deleting from the returned map mutates the agent registry itself, and iterating it
	// while modules are still loading is a data race. Copy before you keep it. Note the
	// inconsistency with ListSkills, which does return a copy.
	ListAgents(ctx core.ServerContext) map[string]ai.Agent

	// RegisterAgentType binds an agent type to the ServiceFactory that builds agents of that
	// type. Call it from a plugin's Initialize, before any agent YAML naming the type is
	// loaded; an agent whose metadata.agentType has no factory fails startup with NotFound.
	// Registration OVERWRITES an existing type silently and always returns nil.
	RegisterAgentType(ctx core.ServerContext, agenttype ai.AgentType, factory core.ServiceFactory) error

	// RegisterAgent allows direct registration of pro-code agents (e.g. golangagent), skipping
	// the YAML + factory path.
	//
	// Keyed on agent.GetName(): an empty name is a Core_Bad_Conf, and unlike RegisterAgentType
	// a DUPLICATE name is rejected rather than overwritten. If the agent satisfies
	// ai.HandoffCapableAgent it is also enrolled with the HandoffManager; a failure to enroll
	// is logged as a warning and does NOT fail this call, so a nil error does not prove the
	// agent is reachable for handoffs.
	RegisterAgent(ctx core.ServerContext, agent ai.Agent) error

	// LLMRequest sends a completion request to the provider owning req.Model and returns the
	// response.
	//
	// Dispatch is BY MODEL NAME ONLY — the manager looks req.Model up in a flat map built from
	// every registered provider's ListModels — so req.Model is required, provider selection is
	// not something the caller controls, and an unknown or empty model is a NotFound
	// ("LLM Model"). Two providers advertising the same model name silently collide: the one
	// registered later wins for every request. req.PreferredModels / FallbackModels are NOT
	// consulted here; any fallback is the provider's own affair.
	LLMRequest(ctx core.RequestContext, req *ai.CompletionRequest) (*ai.CompletionResponse, error)

	// LLMStreamingRequest sends a prompt and streams back responses as a channel of
	// ai.StreamEvent. Model dispatch works exactly as in LLMRequest, including the
	// same-model-name collision. The caller must drain the channel to completion; abandoning
	// it leaks the provider's producer goroutine.
	LLMStreamingRequest(ctx core.RequestContext, req *ai.CompletionRequest) (<-chan ai.StreamEvent, error)

	// GetMCPServer returns the MCP engine wrapper registered at rootpath, or a NotFound error.
	// rootpath is matched EXACTLY as a map key — there is no prefix or longest-match logic —
	// so it must be the same string the engine registered with.
	GetMCPServer(ctx core.ServerContext, rootpath string) (ai.Mcp, error)

	// RegisterMCPServer registers an MCP engine wrapper under rootpath. Called by the MCP
	// engine at startup. Overwrites any existing entry for the same path silently and always
	// returns nil.
	RegisterMCPServer(ctx core.ServerContext, rootpath string, mcpsvr ai.Mcp) error

	// GetLLMProvider returns the provider registered under name, or a NotFound error. This is
	// the provider registration name (e.g. "openai"), not a model name.
	GetLLMProvider(ctx core.ServerContext, name string) (ai.LLMProvider, error)

	// RegisterLLMProvider registers an LLM provider and indexes every model its ListModels
	// reports, so those models become dispatchable by LLMRequest.
	//
	// The model index is flat and shared: a model name already claimed by an earlier provider
	// is REBOUND to this one with no error or warning. An error from ListModels fails the
	// registration. Because the index is built once here, models a provider gains later are
	// invisible until it is registered again.
	RegisterLLMProvider(ctx core.ServerContext, name string, llmprovider ai.LLMProvider) error

	// HasModel reports whether any registered provider advertises modelName.
	//
	// The error return is vestigial — the implementation is a map lookup that returns
	// (true, nil) or (false, nil) and can never produce a non-nil error
	// (agentmanager.go:395). Do not write recovery logic for it. A false here means no
	// provider was registered advertising that name, which at startup usually means the
	// provider module has not loaded yet rather than that the model does not exist.
	HasModel(ctx core.ServerContext, modelName string) (bool, error)

	// ListSkills returns every registered skill keyed by canonical name (the skill descriptor's
	// Metadata.ID, which is not necessarily the service name). Unlike ListAgents this returns a
	// COPY of the map, so it is safe to hold and mutate; the ai.Skill values are still shared.
	ListSkills(ctx core.ServerContext) map[string]ai.Skill

	// ListSkillSummaries returns each registered skill's discovery metadata (id, name, version,
	// description, category, author, tags) keyed by canonical name — the cheap catalogue to
	// show an LLM before activating a skill. Skills whose descriptor is nil are omitted, so
	// this map can be smaller than ListSkills.
	ListSkillSummaries(ctx core.ServerContext) map[string]ai.SkillSummary

	// RegisterSkill allows direct registration of pro-code skills (e.g. golangskill).
	//
	// It calls skill.GetSkillDescriptor immediately (an error there aborts registration),
	// normalizes the descriptor, and files the skill under descriptor.Metadata.ID — so THE
	// CANONICAL NAME IS THE DESCRIPTOR ID, NOT GetName(). The service name is additionally
	// recorded as an alias, so GetSkill resolves either. Duplicate canonical names and
	// conflicting aliases are both rejected with Core_Bad_Conf rather than overwritten.
	RegisterSkill(ctx core.ServerContext, skill ai.Skill) error

	// RegisterSkillType binds a skill-type key to the ServiceFactory that builds skills of that
	// type, for the low-code path where a skill YAML declares `type`. A YAML with no `type`
	// gets "default". Overwrites silently and always returns nil — contrast RegisterSkill,
	// which rejects duplicates.
	RegisterSkillType(ctx core.ServerContext, skillType string, factory core.ServiceFactory) error

	// GetSkill returns a skill by canonical name or by any registered alias (the service name),
	// or a NotFound error. Names are trimmed but matched case-sensitively.
	GetSkill(ctx core.ServerContext, name string) (ai.Skill, error)

	// GetSkillDescriptor returns the normalized descriptor for a skill, resolved by canonical
	// name or alias, or a NotFound error. The result is a DEEP JSON CLONE made per call — safe
	// to mutate, but not cheap in a loop, and mutations never reach the registry.
	GetSkillDescriptor(ctx core.ServerContext, name string) (*ai.SkillDescriptor, error)

	// GetSkillsByTag returns the skills carrying tag, matching the tag's own name and every
	// name up its ParentTag chain, case-insensitively, against both the descriptor's tags and
	// the service's tags.
	//
	// A NIL TAG RETURNS EVERY REGISTERED SKILL, not none — so an unset tag variable silently
	// widens the query to everything instead of failing. There is no error return: an unknown
	// tag and a genuine empty result are indistinguishable.
	GetSkillsByTag(ctx core.ServerContext, tag *core.Tag) []ai.Skill

	// InvokeSkill is the canonical way to invoke a named skill from any context.
	// It creates a fresh RequestContext backed by SkillResponseHandler so that
	// skills can call CompleteStream for HITL pauses or final responses without
	// interference from the caller's response stream (e.g. lazySessionStream).
	//
	// params must contain "sessionId" for streaming to work. THE KEY IS "sessionId" AND ONLY
	// "sessionId" — the implementation reads params.GetString("sessionId") and nothing else
	// (laatooserver/src/core/skillmanager.go:361), so a caller passing "session_id" gets a
	// silently non-streaming invocation, not an error. (This doc previously claimed either
	// spelling worked; it does not.)
	//
	// The returned map is ALWAYS non-nil, on the error path too — a failure returns
	// {"success": false, "error": ..., "skillName": ..., "duration_ms": ...} alongside the
	// error, so branch on err, not on a nil map. On success the map is
	// {"success": true, "skillName", "duration_ms"} merged with the skill's
	// ai.SkillResponse.Output, and the skill's own output keys can OVERWRITE those three. A
	// skill that responds with anything other than a *ai.SkillResponse contributes no output
	// at all and still reports success.
	InvokeSkill(ctx core.RequestContext, skillName string, params utils.StringMap) (utils.StringMap, error)

	// CreateMemory creates a memory bank of the given type under id, delegating to the
	// AgentMemoryManager registered for that MemoryType.
	//
	// Only MemoryTypeSession and MemoryTypeShared have a registered manager in the shipped
	// platform — chromemmemory and laatooreferencememory implement AgentMemoryManager but
	// never call RegisterAgentMemoryManager — so MemoryTypeData and MemoryTypeReferences
	// always fail here with NotFound ("Memory Manager").
	CreateMemory(ctx core.RequestContext, memorytype ai.MemoryType, id string, config map[string]interface{}) (ai.MemoryBank, error)

	// GetMemory returns an existing memory bank of the given type by id, or an error. Same
	// registered-type limitation as CreateMemory. Callers generally want the
	// get-then-create-on-miss pattern the manager itself uses for session banks.
	GetMemory(ctx core.RequestContext, memorytype ai.MemoryType, id string) (ai.MemoryBank, error)

	// RegisterAgentMemoryManager registers the manager that creates banks for a MemoryType.
	// Call it from the memory plugin's Initialize — a plugin that implements
	// ai.AgentMemoryManager but omits this call is unreachable through CreateMemory/GetMemory,
	// which is exactly the state chromemmemory and laatooreferencememory are in today.
	// Overwrites silently and always returns nil.
	RegisterAgentMemoryManager(ctx core.ServerContext, memorytype ai.MemoryType, mgr ai.AgentMemoryManager) error

	// WriteMessageToMemory records a conversation turn in the session memory bank.
	//
	// The session id is read from ctx.GetString("sessionId") — NOT from ctx.GetSession(), as
	// this doc previously claimed (agentmanager.go:298). If that context value is absent the
	// call does nothing.
	//
	// It has no error return and swallows every failure: an empty content, a missing session
	// id, a memory bank that cannot be created, a JSON marshal failure and a failed
	// bank.Add(...) all return quietly, at most a warning in the log. A caller can never tell
	// whether the turn was recorded. The turn is stored as an ai.AIMessage carrying only Role
	// and Content — no actor, no id — inside an AIMemoryItem of Type "message" stamped with
	// the current time.
	WriteMessageToMemory(ctx core.RequestContext, role ai.AgentStakeholder, content string)

	// GetMessagesFromMemory returns all conversation messages for the session in chronological
	// order. Session ID is read from ctx via ctx.GetString("sessionId").
	//
	// Returns nil when there is no sessionId on ctx, when the bank cannot be reached, or when
	// it holds no messages — the three are indistinguishable, and none of them is an error.
	// Only AIMemoryItems of Type "message" whose Content unmarshals as an ai.AIMessage are
	// included; anything else in the bank is skipped silently, so a short result is not
	// evidence of a short conversation. As a side effect it CREATES the session bank if one
	// does not exist.
	GetMessagesFromMemory(ctx core.RequestContext) []ai.ConversationMessage

	// ============================================================
	// HANDOFF MANAGEMENT
	// ============================================================

	// GetHandoffManager returns the centralized handoff manager. The agent manager constructs
	// and starts one during its own Initialize, so this is non-nil on a running server; it can
	// still be nil if called before the agent manager has initialized.
	GetHandoffManager() HandoffManager

	// FindHandoffAgent discovers agents by capability for handoff targeting. A convenience
	// wrapper over HandoffManager.FindAgentByCapability, so it inherits that method's
	// best-match semantics and its NON-DETERMINISTIC tie-breaking — read the notes there
	// before relying on which agent comes back. Errors with InternalError if the handoff
	// manager has not been initialized.
	FindHandoffAgent(ctx core.ServerContext, capabilities []string) (ai.HandoffCapableAgent, error)

	// ============================================================
	// HITL MANAGEMENT
	// ============================================================

	// GetHITLManager returns the server-level HITL coordinator.
	// Returns nil if no HITLManager has been registered.
	GetHITLManager() ai.HITLManager

	// RegisterHITLManager sets the server-level HITL coordinator.
	// Called once at startup by the server or a plugin that owns the implementation.
	//
	// The server already installs a default HITL manager during agent-manager Initialize, so
	// calling this REPLACES a working coordinator rather than filling an empty slot. There is
	// no duplicate check and no error path — it assigns, logs, and returns nil — so two
	// plugins registering leaves whichever ran last in charge, silently, along with any
	// pending pauses recorded against the one it displaced.
	RegisterHITLManager(ctx core.ServerContext, mgr ai.HITLManager) error

	// ============================================================
	// COMPLETION REQUEST FACTORIES
	// Returns pre-configured CompletionRequest instances for common use-cases.
	//
	// EVERY FACTORY HARDCODES A MODEL NAME — "gpt-5-mini" for most, "gemini-2-5-pro" for the
	// high-quality and research variants — with no reference to what the solution has actually
	// configured. If the matching provider is not registered, the request these produce fails
	// LLMRequest with NotFound ("LLM Model"). Treat them as starting points and set .Model, or
	// use WithModel, unless you know that model is available.
	//
	// Each returns a FRESH struct per call, so mutating one is safe. All of them enable cost
	// tracking with BudgetExceededAction "fail", so a request that exceeds MaxCostUSD is
	// refused rather than truncated.
	// See laatooserver/src/core/agentmanager_defaults.go.
	// ============================================================

	// DefaultCompletionRequest returns a CompletionRequest with sensible production defaults:
	// gpt-5-mini, temperature 0.7, 1500 max tokens, no tool calling, streaming off, a $0.10
	// budget and 2 retries. This is the base every other factory below clones and adjusts.
	DefaultCompletionRequest() *ai.CompletionRequest

	// DefaultCompletionRequestCostSensitive returns defaults optimized for minimum cost:
	// 500 max tokens, temperature 0.3, a $0.02 budget, falling back to gemini-1.5-flash.
	DefaultCompletionRequestCostSensitive() *ai.CompletionRequest

	// DefaultCompletionRequestHighQuality returns defaults optimized for quality & reasoning:
	// gemini-2-5-pro, 4000 max tokens, a $0.50 budget and NO fallback models — so if that
	// model is unavailable the request simply fails.
	DefaultCompletionRequestHighQuality() *ai.CompletionRequest

	// DefaultCompletionRequestFast returns defaults optimized for lowest time-to-first-token:
	// 500 max tokens, temperature 0.3, streaming on with a 15s timeout, a $0.05 budget.
	DefaultCompletionRequestFast() *ai.CompletionRequest

	// DefaultCompletionRequestBatching returns defaults optimized for batch processing: the
	// cost-sensitive profile with streaming off, forced caching, IsBatchItem set, lowered
	// priority and a 60s request timeout.
	DefaultCompletionRequestBatching() *ai.CompletionRequest

	// DefaultCompletionRequestResearch returns defaults for research & analysis workloads: the
	// high-quality profile widened to 8000 max tokens, temperature 0.8, streaming on and a
	// $1.00 budget — the most expensive profile here by an order of magnitude.
	DefaultCompletionRequestResearch() *ai.CompletionRequest

	// DefaultCompletionRequestModerationSafe returns defaults for safe/moderated content:
	// temperature 0, TopP 0.9, TopK 10 and a $0.01 budget. The name describes determinism and
	// cost only — NOTHING here enables a moderation or safety filter, and no content check is
	// performed anywhere in this path.
	DefaultCompletionRequestModerationSafe() *ai.CompletionRequest

	// ============================================================
	// COMPLETION REQUEST BUILDERS
	// Immutable-style helpers that clone a base request with targeted overrides.
	// ============================================================

	// CloneWithOverrides copies a CompletionRequest and applies the provided override function
	// to the copy.
	//
	// IT IS NOT A DEEP COPY, despite the older wording of this comment. Exactly four
	// pointer-typed sub-structs are cloned — CostBudget, StreamingConfig, CacheControl and
	// RetryStrategy (agentmanager_defaults.go:168-192). Everything else is shallow-copied and
	// stays SHARED with base: the CostAlert and ThinkingConfig pointers, and the Metadata,
	// Messages, Tools, StopSequences, PreferredModels and FallbackModels maps/slices. Mutating
	// any of those through the clone — appending a message, adding a metadata key — mutates
	// the original too. Replace such a field wholesale in the override function rather than
	// mutating it in place.
	//
	// A nil overrides function is allowed and yields a plain copy. base is dereferenced
	// unconditionally, so a nil base PANICS.
	CloneWithOverrides(base *ai.CompletionRequest, overrides func(*ai.CompletionRequest)) *ai.CompletionRequest

	// SetRequestMetadata sets common tracking fields (RequestID, UserID, AgentName) on the
	// request IN-PLACE and returns the same pointer for chaining — it does not clone, so it
	// mutates the request every other holder of that pointer sees. It also sets CorrelationID
	// to requestID when CorrelationID is currently empty, which the name does not suggest.
	SetRequestMetadata(req *ai.CompletionRequest, requestID, userID, agentName string) *ai.CompletionRequest

	// WithModel returns a clone of req with the given model selected. Shares req's maps and
	// slices — see CloneWithOverrides.
	WithModel(req *ai.CompletionRequest, model string) *ai.CompletionRequest

	// WithTemperature returns a clone of req with the given temperature. Shares req's maps and
	// slices — see CloneWithOverrides.
	WithTemperature(req *ai.CompletionRequest, temp float32) *ai.CompletionRequest

	// WithMaxTokens returns a clone of req with the given max-token limit. Shares req's maps
	// and slices — see CloneWithOverrides.
	WithMaxTokens(req *ai.CompletionRequest, tokens int) *ai.CompletionRequest

	// WithBudget returns a clone of req with the given USD cost cap.
	// The alert threshold is automatically set to 70 % of maxCostUSD.
	//
	// PANICS IF req.CostBudget IS NIL. The override dereferences r.CostBudget unconditionally
	// (agentmanager_defaults.go:225) and CloneWithOverrides leaves a nil CostBudget nil, so
	// this works on a request from the Default*CompletionRequest factories (which always set
	// one) and nil-panics on a hand-built &ai.CompletionRequest{}.
	WithBudget(req *ai.CompletionRequest, maxCostUSD float64) *ai.CompletionRequest

	// WithStreaming returns a clone of req with streaming enabled, and also switches on
	// EnableCostStream.
	//
	// PANICS IF req.StreamingConfig IS NIL, for the same reason as WithBudget
	// (agentmanager_defaults.go:232) — safe on a factory-produced request, a nil deref on a
	// hand-built one.
	WithStreaming(req *ai.CompletionRequest) *ai.CompletionRequest
}

// AgentEngine is RETIRED — the type below is commented out and does not exist. It described a
// pluggable per-vendor agent runtime that owned its own tool and model registries; the
// AgentManager's GetTools/GetModel role was absorbed into the MCP tool registry and the
// LLM-provider model index instead, and AgentManager.GetEngine went with it (see the commented
// line in AgentManager above). The only implementors, adkengine and agentsdkgoengine, now live
// under nolongerreq/. Do not revive this without deciding how it relates to RegisterAgentType,
// which is the surviving extension point for a new kind of agent.
/*
type AgentEngine interface {
	core.ServerElement
	GetTools(ctx core.ServerContext, toolNames []string) (utils.StringMap, error)
	RegisterTool(ctx core.ServerContext, toolName string, svc Service) error
	GetModel(ctx core.ServerContext, modelName string) (interface{}, error)
	RegisterModel(ctx core.ServerContext, modelName string, model interface{}) error
}
*/
