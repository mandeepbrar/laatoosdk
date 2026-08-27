package ai

import "laatoo.io/sdk/server/core"

// Skill represents a modular expertise package that agents can discover and use.
//
// A Skill is an ordinary invokable service; what makes it a skill is that the agent manager
// indexes it by descriptor so agents and the DelegateToSkill activity can find it by name or
// tag. Register one with AgentManager.RegisterSkill (pro-code) or let a skill YAML in a module
// produce a DefaultSkillService (low-code). Invoke it through AgentManager.InvokeSkill rather
// than calling Invoke directly — InvokeSkill installs the response handler that makes
// CompleteStream/HITL pauses work.
type Skill interface {
	core.UserInvokableService

	// GetSkillDescriptor returns the skill's full descriptor: identity, discovery guidance,
	// instructions, declared tools, knowledge refs and LLM defaults.
	//
	// It is called once, at registration time, and the result is cached by the agent manager;
	// returning a value that changes between calls has no effect after the first. An error
	// here aborts registration of the skill.
	//
	// The registry does not store what this returns verbatim. The agent manager normalizes the
	// descriptor first — filling Metadata.ID, Name, Description, Category, Discovery.Summary
	// and Tags from the service when the descriptor leaves them empty — and it is the
	// normalized Metadata.ID, not GetName(), that becomes the skill's canonical registry key.
	// See laatooserver/src/core/skillmanager.go:181.
	GetSkillDescriptor(ctx core.ServerContext) (*SkillDescriptor, error)

	// GetSkillType returns the skill's category, used as a fallback for
	// SkillDescriptor.Metadata.Category during normalization. It is NOT the skill-type key
	// that RegisterSkillType/createSkillService dispatch on — that one is read from the
	// skill YAML's `type` field before any Skill instance exists.
	GetSkillType() string

	// GetExamples returns worked examples of applying the skill.
	//
	// Every pro-code skill in the platform returns nil here (agentbuildingskill.go:217 and the
	// laatoodesigner *BuildingSkill family all do), so only a low-code DefaultSkillService —
	// which returns its YAML's instructions.examples — ever answers non-nil. Treat an empty
	// result as "not populated", never as "this skill has no applicable examples". Note also
	// that nothing in the platform calls this method; examples reach a consumer through
	// SkillDescriptor.Instructions.Examples instead.
	GetExamples() []Example
}
