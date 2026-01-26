package core

import (
	"laatoo.io/sdk/utils"
)

type InformationBucket interface {
	ConfigurableObjectInfo
}

type Agent interface {
	Service
	GetAgentType() string
	GetVersion() string
	GetDescription() string
	GetAgentProperties() utils.StringMap
	
	// Memory Management
	GetMemory(ctx RequestContext, id string) (MemoryBank, error)
	AddMemory(ctx RequestContext,  bank MemoryBank) error
}

type AgentConversation interface {
	GetId() string
	AssignAgent(Agent)
	GetPresentAgent() Agent
	GetHistory() utils.StringsMap
	AddHistory(ctx RequestContext, actor string, input utils.StringMap)
}

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

// MemoryItem is a logical unit of memory (DTO).
// It acts as a data transfer object between Agents and MemoryBanks.
type MemoryItem struct {
	Type       string                 `json:"type"` // e.g., "message", "artifact", "datarecord"
	Content    interface{}            `json:"content"`
	Importance float64                `json:"importance"`
	Timestamp  string                 `json:"timestamp"` // ISO8601 string
	Tags       []string               `json:"tags"`
	Metadata   utils.StringMap `json:"metadata"`
}

// MemoryBank manages storage and retrieval of MemoryItems.
type MemoryBank interface {
	GetId() string
	// Add stores an item.
	Add(ctx RequestContext, item MemoryItem) error
	
	// Retrieve fetches items based on a query.
	Retrieve(ctx RequestContext, query string, opts utils.StringMap) ([]MemoryItem, error)
	
	Clear(ctx RequestContext) error
	Synthesize(ctx RequestContext) error
} // End of MemoryBank interface

// Skill represents a modular expertise package that agents can discover and use
type Skill interface {
	Service
	GetSkillMetadata() SkillMetadata
	GetInstructions() SkillInstruction
	GetTools() []ToolDefinition
	GetResources() []SkillResource
}

// SkillInfo contains the complete skill definition (Level 1 + 2 + 3) (rename of previous Skill struct)
type SkillInfo struct {
	Metadata     SkillMetadata     `json:"metadata"`
	Instructions SkillInstruction  `json:"instructions"`
	Tools        []ToolDefinition  `json:"tools"`
	Resources    []SkillResource   `json:"resources,omitempty"`
	Type         string            `json:"type,omitempty"`
	CreatedAt    string            `json:"created_at,omitempty"`
	UpdatedAt    string            `json:"updated_at,omitempty"`
}

// SkillMetadata provides discovery information (Level 1)
type SkillMetadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Author      string   `json:"author,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// SkillInstruction contains procedural knowledge (Level 2)
type SkillInstruction struct {
	Title          string         `json:"title"`
	Overview       string         `json:"overview"`
	Steps          []string       `json:"steps"`
	BestPractices  []string       `json:"best_practices,omitempty"`
	Examples       []SkillExample `json:"examples,omitempty"`
	References     []Reference    `json:"references,omitempty"`
	ContextHints   string         `json:"context_hints,omitempty"`
	ErrorHandling  map[string]string `json:"error_handling,omitempty"`
}

// SkillExample shows concrete application patterns
type SkillExample struct {
	Description    string                 `json:"description"`
	Input          map[string]interface{} `json:"input,omitempty"`
	ExpectedOutput string                 `json:"expected_output,omitempty"`
	Notes          string                 `json:"notes,omitempty"`
}

// Reference provides external knowledge
type Reference struct {
	Title string `json:"title"`
	URL   string `json:"url,omitempty"`
	Type  string `json:"type"` // "documentation", "best-practice", "template"
}

// ToolDefinition describes what tools this skill can invoke
type ToolDefinition struct {
	ToolName    string                 `json:"tool_name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Required    []string               `json:"required"`
	Annotations ToolAnnotations        `json:"annotations"`
}

// ToolAnnotations provide hints about tool behavior
type ToolAnnotations struct {
	ReadOnly     bool `json:"read_only"`
	Destructive  bool `json:"destructive"`
	Idempotent   bool `json:"idempotent"`
	OpenWorld    bool `json:"open_world"`     // accesses external resources
	RequiresAuth bool `json:"requires_auth"`
}

// SkillResource represents Level 3: Executable tools and templates
type SkillResource struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"` // "script", "template", "reference"
	ContentPath  string   `json:"content_path"`
	Executable   bool     `json:"executable"`
	Dependencies []string `json:"dependencies,omitempty"`
	Timeout      int      `json:"timeout,omitempty"` // seconds
}


