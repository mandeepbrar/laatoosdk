package ai

import "laatoo.io/sdk/server/core"

// Example shows concrete application patterns
type Example struct {
	Description    string                 `json:"description"`
	Request        core.RequestInfo `json:"request,omitempty"`
	Response       core.ResponseInfo `json:"response,omitempty"`
	Notes          string                 `json:"notes,omitempty"`
}

type Instruction interface {
	core.Storable
	GetText() string
	GetDescription() string
	IsTemplate() bool
	Params() map[string]core.Param
}

type Prompt Instruction