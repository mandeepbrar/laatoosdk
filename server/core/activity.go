package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/utils"
)

// ActivityType defines the nature of the activity
type ActivityType string

// The built-in activity types. Any other value names a type served by a provider registered with
// elements.ActivityManager.RegisterActivityProvider -- ActivityTypeFlow is the one the platform ships
// a provider for.
//
// There is no executor type: a Go function is registered as an ACTION
// (elements.ActionManager.RegisterAction), and a workflow step names the action directly.
const (
	ActivityTypeManual  ActivityType = "manual"
	ActivityTypeService ActivityType = "service"
	ActivityTypeScript  ActivityType = "script"
	// ActivityTypeFlow is a low-code function: its Body is a workflow-DSL statement tree run
	// synchronously within the caller's request, returning its declared Output.
	ActivityTypeFlow ActivityType = "flow"
)

// Activity represents the interface for a workflow step, treating it as a specialized service.
type Activity interface {
	UserInvokableService
	// GetDefinition returns the configuration definition of the activity.
	GetDefinition() *ActivityDefinition
}

// ActivityDefinition represents the configuration of a step in the workflow (YAML schema)
type ActivityDefinition struct {
	Name             string           `json:"name" yaml:"name"`
	Activity         string           `json:"activity" yaml:"activity"` // Function ID (e.g. "user.SetActive")
	ActivityType     ActivityType     `json:"activity_type" yaml:"activity_type"`
	AccessPermission string           `json:"access_permission,omitempty" yaml:"access_permission,omitempty"`
	Condition        string           `json:"condition,omitempty" yaml:"condition,omitempty"`
	Arguments        []string         `json:"arguments,omitempty" yaml:"arguments,omitempty"`
	Config           config.Config    `json:"config,omitempty" yaml:"config,omitempty"`
	Result           utils.StringsMap `json:"result,omitempty" yaml:"result,omitempty"`
	Timeout          int              `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Retry            *RetryPolicy     `json:"retry,omitempty" yaml:"retry,omitempty"`
	HumanTask        *HumanTaskConfig `json:"human_task,omitempty" yaml:"human_task,omitempty"`
	// Streaming marks this activity as producing a streaming response.
	// When true, the service invoke path spawns a goroutine and the
	// activity manager drains the ResponseStream channel after execution.
	Streaming bool `json:"streaming,omitempty" yaml:"streaming,omitempty"`
	// Input declares what a caller passes, as a JSON-Schema object
	// ({type: object, properties: {...}, required: [...]}), the shape a workflow's inputSchema
	// takes. When set, a call missing a required property is refused before the activity runs.
	Input map[string]interface{} `json:"input,omitempty" yaml:"input,omitempty"`
	// Output declares what the activity returns, in the same shape. When set, the caller receives
	// only the declared properties of the activity's result.
	Output map[string]interface{} `json:"output,omitempty" yaml:"output,omitempty"`
	// Body is the definition a provider-backed type runs -- for ActivityTypeFlow, the workflow-DSL
	// statement tree (the value a workflow's inlineDefinition.root takes). The built-in types
	// ignore it.
	Body map[string]interface{} `json:"body,omitempty" yaml:"body,omitempty"`
}

type RetryPolicy struct {
	MaxAttempts     int    `json:"max_attempts" yaml:"max_attempts"`
	Backoff         string `json:"backoff" yaml:"backoff"`
	InitialInterval int    `json:"initial_interval" yaml:"initial_interval"`
	MaxInterval     int    `json:"max_interval" yaml:"max_interval"`
}

type HumanTaskConfig struct {
	Assignee       string                 `json:"assignee,omitempty" yaml:"assignee,omitempty"`
	CandidateRoles []string               `json:"candidate_roles,omitempty" yaml:"candidate_roles,omitempty"`
	CandidateUsers []string               `json:"candidate_users,omitempty" yaml:"candidate_users,omitempty"`
	DueDate        string                 `json:"due_date,omitempty" yaml:"due_date,omitempty"`
	Priority       string                 `json:"priority,omitempty" yaml:"priority,omitempty"`
	FormSchema     map[string]interface{} `json:"form_schema,omitempty" yaml:"form_schema,omitempty"`
	TaskQueue      string                 `json:"task_queue,omitempty" yaml:"task_queue,omitempty"`
	TaskManager    string                 `json:"task_manager,omitempty" yaml:"task_manager,omitempty"`
}
