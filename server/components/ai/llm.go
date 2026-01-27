package ai

import (
	"time"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/utils"
)

// LLMProvider is the core interface for all LLM providers
type LLMProvider interface {
	// Complete sends a prompt and gets a response
	Complete(ctx core.RequestContext, req *CompletionRequest) (*CompletionResponse, error)

	// Stream sends a prompt and streams back responses
	// Returns a channel of StreamEvent
	Stream(ctx core.RequestContext, req *CompletionRequest) (<-chan StreamEvent, error)

	// CountTokens estimates tokens for input text
	CountTokens(ctx core.RequestContext, text string, model string) (int, error)

	// GetConfig returns model configuration
	GetConfig(ctx core.RequestContext, model string) (*ModelConfig, error)

	// ListModels returns available models
	ListModels(ctx core.RequestContext) ([]string, error)

	// Health checks if provider is accessible
	Health(ctx core.RequestContext) error

	// Name returns provider name
	Name() string

}


// ModelConfig contains model-level configuration
type ModelConfig struct {
	Name            string
	Provider        string // "openai", "anthropic", "gemini", "ollama", etc.
	ContextWindow   int    // Max context tokens
	MaxOutput       int    // Max output tokens
	CostPer1KInput  float64 // USD per 1K input tokens
	CostPer1KOutput float64 // USD per 1K output tokens
	Capabilities    ModelCapabilities
	Metadata        map[string]interface{}
}

// ModelCapabilities describes what a model can do
type ModelCapabilities struct {
	SupportsVision        bool
	SupportsStreaming     bool
	SupportsServiceCall  bool
	SupportsJSON          bool
	SupportsSystemPrompt  bool
	MaxContextTokens      int
	ReleaseDate           time.Time
}


// CompletionRequest represents a request to complete a prompt
type CompletionRequest struct {
	Model              string              // Model identifier (e.g., "gpt-4o", "claude-3.5-sonnet")
	Messages           []ConversationMessage           // Conversation history
	Temperature        float32             // 0.0 to 2.0 (default: 1.0)
	MaxTokens          int                 // Max output tokens
	TopP               float32             // 0.0 to 1.0 (default: 1.0)
	TopK               int                 // Top-K sampling
	StopSequences      []string            // Stop generation at these
	SystemPrompt       string              // System-level instructions
	Tools          []Tool          // Available tools/functions
	ToolCallMode   ToolCallMode    // "none", "auto", "required"
	ObjectResponse     bool      // JSON schema, text, etc.
	ResponseObjectName datatypes.Serializable
	Metadata           utils.StringMap   // Custom metadata
}

// CompletionResponse represents a response from LLM
type CompletionResponse struct {
	Model            string              // Model used
	Content          string              // Generated text
	Tokens           TokenUsage          // Token breakdown
	Cost             Cost                // Request cost
	FinishReason     FinishReason        // Why generation stopped
	ToolRequests    []core.Request      // Tool calls made
	RequestID        string              // For tracking
	Latency          time.Duration       // Response time
	Timestamp        time.Time           // When request was made
	Metadata         utils.StringMap   // Response metadata
}

// TokenUsage tracks token consumption
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
	CachedTokens     int `json:"cached_tokens,omitempty"` // If using cache
}

// Cost tracks financial cost
type Cost struct {
	PromptCost      float64 `json:"prompt_cost_usd"`
	CompletionCost  float64 `json:"completion_cost_usd"`
	TotalCost       float64 `json:"total_cost_usd"`
	Currency        string  `json:"currency"` // "USD"
}

// FinishReason why generation stopped
type FinishReason string

const (
	FinishReasonStop       FinishReason = "stop"
	FinishReasonLength     FinishReason = "length"
	FinishReasonFunctionCall FinishReason = "function_call"
	FinishReasonError      FinishReason = "error"
	FinishReasonUnknown    FinishReason = "unknown"
)


// StreamEvent represents a single streaming event
type StreamEvent struct {
	Type      string    // "token", "function_call", "error", "done"
	Token     string    // For token events
	Delta     TokenUsage // Incremental tokens
	Cost      Cost      // Incremental cost
	ToolRequest core.Request // For function calls
	Error     error     // For errors
	Timestamp time.Time
}


// ContextWindowManager tracks and manages context window usage
type ContextWindowManager struct {
	TotalWindow      int
	UsedPrompt       int
	UsedCompletion   int
	ReservedForOutput int
}

// GetAvailableContext returns how many tokens are available for new input
func (cwm *ContextWindowManager) GetAvailableContext() int {
	used := cwm.UsedPrompt + cwm.UsedCompletion + cwm.ReservedForOutput
	return cwm.TotalWindow - used
}

// CanFitPrompt checks if prompt fits in context window
func (cwm *ContextWindowManager) CanFitPrompt(promptTokens int) bool {
	return cwm.GetAvailableContext() > promptTokens
}
