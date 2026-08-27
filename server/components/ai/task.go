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

// TaskProvider is a factory for one task TYPE. A plugin implements it and the laatooai module
// picks it up by type-asserting the plugin's module object during Load, or the plugin calls
// RegisterTaskProvider directly; either way the provider is stored under GetName() and looked
// up by a task's `type` key in the agent YAML. Registration OVERWRITES silently, so two
// plugins claiming the same provider name leave whichever loaded last in place with no warning.
// See laatoomodules/ai/dev/plugins/laatooai/src/server/go/laatooaimodule.go:49.
type TaskProvider interface {
	// GetName returns the task TYPE this provider builds — "default", "activity", and so on —
	// not the name of any task instance. It is the registry key, so it must match the `type`
	// field of the task entries in an agent YAML.
	GetName() string

	// GetAgenticTaskExecutor builds one Task instance from its config block.
	//
	// WHETHER THE RETURNED TASK IS INITIALIZED IS NOT SETTLED. The two shipped providers
	// disagree: ActivityTaskProvider calls Initialize inside this method
	// (laatooai/src/sdk/go/tasks/activity_task.go:92), while DefaultTaskProvider only parses
	// the config into fields and leaves Initialize to the agent, which calls it later over
	// every task (goalagents/agent.go:261). So a task built by the default provider gets
	// Initialize twice via the activity provider's path and once otherwise. Write Initialize
	// to be idempotent.
	//
	// The returned Task is NOT nil when err is non-nil — ActivityTaskProvider returns
	// `(task, err)` unconditionally — so branch on err, never on a nil task.
	GetAgenticTaskExecutor(ctx core.ServerContext, conf config.Config, agt Agent) (Task, error)
}

// Task represents an executable unit of work inside an agent: one step of a goal agent's plan,
// scheduled by an ExecutionStrategy according to GetDependencies and run through Execute.
//
// The lifecycle is construct (TaskProvider.GetAgenticTaskExecutor) → Initialize → Execute, and
// only Execute receives a RequestContext. Anything needing per-request state must be read
// inside Execute or out of the MemoryBank handed to it, never cached at Initialize time.
type Task interface {
	// GetName returns the task's id — for the default provider, the `id` key of its config
	// block, NOT the YAML filename. It is the key an ExecutionStrategy uses for the dependency
	// graph, so a task whose config omits `id` gets the empty string and collides with every
	// other such task in the same agent.
	GetName() string

	// GetPurpose returns what the task is for. The default implementation returns the config's
	// `description`; it is fed to the planner LLM as the task's rationale, so an empty
	// description leaves the planner nothing to select on.
	GetPurpose() string

	// GetConfig returns the config block the task was built from. The goal agent re-passes
	// this same value back into Initialize (`task.Initialize(ctx, task.GetConfig(), s)`), so
	// an implementation that returns nil here breaks its own initialization.
	GetConfig() config.Config

	// Tools returns the tools this task may call, keyed by tool name.
	//
	// TOOLS NAMED IN CONFIG BUT NOT RESOLVABLE ARE DROPPED WITHOUT A WORD. DefaultTask.Initialize
	// looks each configured tool name up and appends it only `if err == nil`, with no else and
	// no log (goalagents/defaulttask.go:157-163), so a typo'd or not-yet-loaded tool leaves a
	// map that is quietly smaller than the YAML declares and the task runs without it.
	Tools(ctx core.ServerContext) map[string]Tool

	// Skills returns the skills this task may delegate to, keyed by skill name. It has exactly
	// the same silent-drop behaviour as Tools — an unresolvable skill name in the YAML is
	// skipped with no error and no log (goalagents/defaulttask.go:147-153).
	Skills(ctx core.ServerContext) map[string]Skill

	// GetDependencies returns the names of tasks that must complete before this one, read from
	// the config's `dependsOn`. The ExecutionStrategy builds its graph from these; a name that
	// matches no task in the agent is not validated here.
	GetDependencies(ctx core.ServerContext) []string

	// GetStructuredOutput reports whether the task's LLM output should be parsed as structured
	// data rather than taken as prose. Despite the RequestContext parameter, every shipped
	// implementation returns a value fixed at config time (`structured_output`) and ignores
	// the context entirely — it cannot vary per request.
	GetStructuredOutput(ctx core.RequestContext) bool

	// Initialize wires the task to its owning agent and resolves its declared tools, skills and
	// dependencies. Called by the agent over all its tasks after the agent itself is
	// configured, and possibly a second time by the provider — see GetAgenticTaskExecutor.
	//
	// agt is passed as the ai.Agent interface but DefaultTask asserts it to the concrete
	// *GoalAgentService without checking (goalagents/defaulttask.go:138), so handing a
	// DefaultTask to any other agent type panics rather than erroring. Returning an error here
	// aborts agent startup.
	Initialize(ctx core.ServerContext, conf config.Config, agt Agent) error

	// Execute runs the task against the supplied memory bank and returns its result.
	//
	// Return a *TaskResult even on failure, with Success=false and Error set — that is what
	// the shipped implementations do, and strategies read Success rather than only the error.
	//
	// BEWARE THE EMBEDDED BASE. tasks.BaseTask, which custom tasks are meant to embed, has an
	// Execute that does no work and returns `{TaskName, Success: true}`
	// (laatooai/src/sdk/go/tasks/base_task.go:165). A task that embeds BaseTask and forgets to
	// override Execute reports success on every run without executing anything, and no
	// validator catches it.
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
