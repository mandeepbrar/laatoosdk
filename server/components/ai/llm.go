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
	CountTokens(ctx core.ServerContext, text string, model string) (int, error)

	// GetConfig returns model configuration
	GetConfig(ctx core.ServerContext, model string) (*ModelConfig, error)

	// ListModels returns available models
	ListModels(ctx core.ServerContext) ([]string, error)

	// Health checks if provider is accessible
	Health(ctx core.ServerContext) error

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

	// NEW: Cost Optimization
	CostBudget              *CostBudget
	TokenCostEstimate       *TokenCostEstimate
	EnableCostTracking      bool
	CostAlert               *CostAlert

	// NEW: Streaming Configuration
	StreamingConfig         *StreamingConfig
	EnableThinkingStream    bool
	TokenStreamingGranularity int
	EnableCostStream        bool

	// NEW: Cache & Performance
	CacheControl            *CacheControl
	PreRequestCacheCheck    bool
	AllowCachedResponse     bool
	CacheMaxAge             int

	// NEW: Priority & Routing
	Priority                int
	PreferredModels         []string
	FallbackModels          []string
	RouteByComplexity       bool

	// NEW: Request Tracking
	RequestID               string
	CorrelationID           string
	AgentName               string
	SessionID               string
	UserID                  string

	// NEW: Advanced Reasoning
	ThinkingConfig          *ThinkingConfig

	// NEW: Batch Operations
	BatchID                 string
	IsBatchItem             bool

	// NEW: Timeout & Reliability
	RequestTimeoutMs        int
	FirstTokenTimeoutMs     int
	RetryStrategy           *RetryStrategy
}

// CostBudget Control
type CostBudget struct {
	MaxCostUSD            float64 // Hard limit
	AlertThresholdUSD     float64 // Alert when exceeded
	BudgetExceededAction  string  // "fail", "truncate", "fallback_model"
	AllowPartialExecution bool    // Allow returning partial response
}

type TokenCostEstimate struct {
	EstimatedInputTokens  int
	EstimatedOutputTokens int
	EstimatedCostUSD      float64
	Confidence            float32 // 0.0-1.0
	LastUpdated           time.Time
}

type CostAlert struct {
	Enabled        bool
	ThresholdUSD   float64
	AlertHandler   func(alert *CostAlertEvent)
	IncludeMetrics bool
}

type CostAlertEvent struct {
	RequestID         string
	EstimatedCost     float64
	ThresholdUSD      float64
	Percentage        float32
	Timestamp         time.Time
	Model             string
	TokenEstimate     *TokenCostEstimate
	RecommendedAction string
}

type StreamingConfig struct {
	Enabled            bool
	BufferSize         int
	ChunkSize          int
	TimeoutMs          int
	EventFilters       []string
	CompressionEnabled bool
	PingIntervalMs     int
}

type CacheControl struct {
	Mode             CacheMode // "auto", "force", "disable", "refresh"
	TTLSeconds       int
	CacheKeyCustom   string
	EphemeralCache   bool
	ReturnCacheStats bool
}

type ThinkingConfig struct {
	Enabled                 bool
	BudgetTokens           int    // 1024-8000
	Quality                ThinkingQuality // "fast", "standard", "thorough"
	Type                   string
	StreamThinking         bool
	IncludeThinkingInOutput bool
}

type RetryStrategy struct {
	MaxRetries        int
	InitialBackoffMs  int
	MaxBackoffMs      int
	BackoffMultiplier float64
	JitterEnabled     bool
	RetryableErrors   []string
}

// CompletionResponse represents a response from LLM
type CompletionResponse struct {
	Model         string            // Model used
	Content       string            // Generated text
	Tokens        TokenUsage        // Token breakdown
	Cost          Cost              // Request cost
	FinishReason  FinishReason      // Why generation stopped
	FunctionCalls []FunctionCall    // Simple function calls
	RequestID     string            // For tracking
	Latency       time.Duration     // Response time
	CacheUsed     bool              // If cache was used
	CacheMetadata utils.StringMap   // Cache details
	Timestamp     time.Time         // When request was made
	Metadata      utils.StringMap   // Response metadata
}

// FunctionCall represents a function call made by LLM
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON string
	ID        string `json:"id,omitempty"`
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


// StreamEventType why generation stopped
type StreamEventType string

const (
	StreamEventToken        StreamEventType = "token"
	StreamEventFunctionCall StreamEventType = "function_call"
	StreamEventError        StreamEventType = "error"
	StreamEventDone         StreamEventType = "done"
	StreamEventThinking     StreamEventType = "thinking"
	StreamEventCost         StreamEventType = "cost"
)

// CacheMode defines caching behavior
type CacheMode string

const (
	CacheModeAuto    CacheMode = "auto"
	CacheModeForce   CacheMode = "force"
	CacheModeDisable CacheMode = "disable"
	CacheModeRefresh CacheMode = "refresh"
)

// ThinkingQuality defines reasoning quality levels
type ThinkingQuality string

const (
	ThinkingQualityFast     ThinkingQuality = "fast"
	ThinkingQualityStandard ThinkingQuality = "standard"
	ThinkingQualityThorough ThinkingQuality = "thorough"
)

// StreamEvent represents a single streaming event
type StreamEvent struct {
	Type            StreamEventType // "token", "function_call", "error", "done", "thinking", "cost"
	Token           string          // For token events
	ThinkingContent string          // For thinking events
	Delta           TokenUsage      // Incremental tokens
	Cost            Cost            // Incremental cost
	ToolRequest     core.Request    // For function calls
	Error           error           // For errors
	Timestamp       time.Time
	CostDelta       *StreamingCostEvent
	TokenDelta      *StreamingTokenEvent
	ThinkingDelta   *StreamingThinkingEvent
	Index           int
	Metadata        map[string]interface{}
}

type StreamingCostEvent struct {
	DeltaCostUSD           float64
	CumulativeCostUSD      float64
	ChunkInputTokens       int
	ChunkOutputTokens      int
	CumulativeInputTokens  int
	CumulativeOutputTokens int
	BudgetPercentage       float32
	ProjectedTotalCost     float64
}

type StreamingTokenEvent struct {
	Text             string
	TokenCount       int
	CumulativeTokens int
	TokensPerSecond  float64
	TimeToFirstToken time.Duration
}

type StreamingThinkingEvent struct {
	Content                 string
	ThinkingTokensUsed      int
	ThinkingBudgetRemaining int
	QualityIndicator        string
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
