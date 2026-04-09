package ai

import "laatoo.io/sdk/server/core"

// Skill represents a modular expertise package that agents can discover and use
type Skill interface {
	core.UserInvokableService
	GetSkillDescriptor(ctx core.ServerContext) (*SkillDescriptor, error)
	GetSkillType() string
	GetExamples() []Example
}
