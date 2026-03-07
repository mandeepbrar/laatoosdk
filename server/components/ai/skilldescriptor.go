package ai

// SkillDescriptor is the shared runtime shape for both low-code and pro-code skills.
type SkillDescriptor struct {
	Metadata      SkillSummary          `json:"metadata"`
	Discovery     SkillDiscovery        `json:"discovery,omitempty"`
	Instructions  SkillInstructions     `json:"instructions,omitempty"`
	Knowledge     []SkillKnowledgeRef   `json:"knowledge,omitempty"`
	Tools         []SkillToolDefinition `json:"tools,omitempty"`
	Resources     []SkillResourceRef    `json:"resources,omitempty"`
	LLM           SkillLLMConfig        `json:"llm,omitempty"`
	Capabilities  []string              `json:"capabilities,omitempty"`
	Prerequisites SkillPrerequisites    `json:"prerequisites,omitempty"`
}

// SkillSummary contains the always-on identity and discovery metadata.
type SkillSummary struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Version     string   `json:"version,omitempty"`
	Description string   `json:"description,omitempty"`
	Category    string   `json:"category,omitempty"`
	Author      string   `json:"author,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// SkillDiscovery contains lightweight activation guidance used before full activation.
type SkillDiscovery struct {
	Summary            string   `json:"summary,omitempty"`
	WhenToUse          []string `json:"when_to_use,omitempty"`
	WhenNotToUse       []string `json:"when_not_to_use,omitempty"`
	ActivationExamples []string `json:"activation_examples,omitempty"`
}

// SkillInstructions contains activation-time procedural guidance.
type SkillInstructions struct {
	Title         string            `json:"title,omitempty"`
	Overview      string            `json:"overview,omitempty"`
	Steps         []string          `json:"steps,omitempty"`
	BestPractices []string          `json:"best_practices,omitempty"`
	ContextHints  string            `json:"context_hints,omitempty"`
	ErrorHandling map[string]string `json:"error_handling,omitempty"`
	Examples      []Example         `json:"examples,omitempty"`
	References    []SkillReference  `json:"references,omitempty"`
	FullText      *SkillFileRef     `json:"full_text,omitempty"`
}

// SkillReference is an inline human-authored reference entry.
type SkillReference struct {
	Title string `json:"title,omitempty"`
	Type  string `json:"type,omitempty"`
	URL   string `json:"url,omitempty"`
}

// SkillFileRef identifies file-backed content belonging to a skill.
type SkillFileRef struct {
	Path    string `json:"path,omitempty"`
	Summary string `json:"summary,omitempty"`
	Content string `json:"content,omitempty"`
}

// SkillKnowledgeRef identifies file-backed knowledge that may be activated on demand.
type SkillKnowledgeRef struct {
	ID      string `json:"id,omitempty"`
	Kind    string `json:"kind,omitempty"`
	Path    string `json:"path,omitempty"`
	Summary string `json:"summary,omitempty"`
	Load    string `json:"load,omitempty"`
	Content string `json:"content,omitempty"`
}

// SkillToolDefinition describes a tool exposed through a skill descriptor.
type SkillToolDefinition struct {
	ToolName    string                 `json:"tool_name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Required    []string               `json:"required,omitempty"`
	Annotations ToolAnnotations        `json:"annotations,omitempty"`
}

// SkillResourceRef describes non-prompt assets owned by a skill.
type SkillResourceRef struct {
	ID         string `json:"id,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Path       string `json:"path,omitempty"`
	Summary    string `json:"summary,omitempty"`
	Executable bool   `json:"executable,omitempty"`
	Timeout    int    `json:"timeout,omitempty"`
}

// SkillLLMConfig describes skill-level LLM behavior.
type SkillLLMConfig struct {
	Enabled           bool               `json:"enabled"`
	Profile           string             `json:"profile,omitempty"`
	InheritsFromAgent bool               `json:"inherits_from_agent,omitempty"`
	ToolCallMode      ToolCallMode       `json:"tool_call_mode,omitempty"`
	Overrides         *SkillLLMOverrides `json:"overrides,omitempty"`
}

// SkillLLMOverrides allows per-skill tuning of the selected LLM profile.
type SkillLLMOverrides struct {
	Model       *string  `json:"model,omitempty"`
	Temperature *float32 `json:"temperature,omitempty"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
	MaxCostUSD  *float64 `json:"max_cost_usd,omitempty"`
	Streaming   *bool    `json:"streaming,omitempty"`
}

// SkillPrerequisites captures optional tool, knowledge, and skill prerequisites.
type SkillPrerequisites struct {
	Tools     []string `json:"tools,omitempty"`
	Knowledge []string `json:"knowledge,omitempty"`
	Skills    []string `json:"skills,omitempty"`
}
