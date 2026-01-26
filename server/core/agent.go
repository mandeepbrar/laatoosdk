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
}

type AgentConversation interface {
	GetId() string
	AssignAgent(Agent)
	GetPresentAgent() Agent
	GetHistory() utils.StringsMap
	AddHistory(ctx RequestContext, actor string, input utils.StringMap)
}

// Skill represents a modular expertise package that agents can discover and use
type Skill struct {
	Metadata     SkillMetadata     `json:"metadata"`
	Instructions SkillInstruction  `json:"instructions"`
	Service      *SkillServiceRef  `json:"service,omitempty"`
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
	Title         string         `json:"title"`
	Overview      string         `json:"overview"`
	Steps         []string       `json:"steps"`
	BestPractices []string       `json:"best_practices,omitempty"`
	Examples      []SkillExample `json:"examples,omitempty"`
	References    []Reference    `json:"references,omitempty"`
	ContextHints  string         `json:"context_hints,omitempty"`
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

// SkillServiceRef references a Laatoo service for execution
type SkillServiceRef struct {
	ServiceName string `json:"service_name"`
	Method      string `json:"method,omitempty"` // defaults to "Invoke"
}

