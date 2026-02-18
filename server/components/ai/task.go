package ai

import (
	"time"

	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// TaskType defines different types of tasks
type TaskType string

const (
	TaskTypeAutomated TaskType = "automated" // LLM-driven task
	TaskTypeHuman     TaskType = "human"     // Human-in-the-loop
	TaskTypeSaga      TaskType = "saga"      // Saga step with compensation
	TaskTypeService   TaskType = "service"   // External service call via MCP
	TaskTypeActivity  TaskType = "activity"  // Backend activity execution
)

// TaskProvider interface for registering and creating tasks
type TaskProvider interface {
	GetName() string
	GetAgenticTaskExecutor(ctx core.ServerContext, conf config.Config, agt Agent) (Task, error)
}

// Task represents an executable unit of work
type Task interface {
	GetName() string
	GetPurpose() string
	GetConfig() config.Config
	Tools(ctx core.ServerContext) map[string]Tool
	Skills(ctx core.ServerContext) map[string]Skill
	GetDependencies(ctx core.ServerContext) []string
	GetStructuredOutput(ctx core.RequestContext) bool
	Initialize(ctx core.ServerContext, conf config.Config, agt Agent) error
	Execute(ctx core.RequestContext, memory MemoryBank) (*TaskResult, error)
}


// ToolCall represents a function call made by the LLM
type ToolCall struct {
	Name      string
	Arguments string
	Output    string
}

// TaskResult tracks per-task metrics
type TaskResult struct {
	TaskName  string
	Success   bool
	Duration  time.Duration
	Cost      float64
	Output    string
	Error     string
	Data      utils.StringMap
	ToolCalls []ToolCall
}
