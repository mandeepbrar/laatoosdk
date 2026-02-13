package core

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/utils"
)

// ActivityType defines the nature of the activity
type ActivityType string

const (
	ActivityTypeManual    ActivityType = "manual"
	ActivityTypeService   ActivityType = "service"
	ActivityTypeExecutor  ActivityType = "executor"
	ActivityTypeScript    ActivityType = "script"
)

// Activity represents the interface for a workflow step, treating it as a specialized service.
type Activity interface {
	UserInvokableService
	// GetDefinition returns the configuration definition of the activity.
	GetDefinition() *ActivityDefinition
}

// ActivityDefinition represents the configuration of a step in the workflow (YAML schema)
type ActivityDefinition struct {
	Name         string                 `json:"name" yaml:"name"`
	Activity     string                 `json:"activity" yaml:"activity"` // Function ID (e.g. "user.SetActive")
	ActivityType     ActivityType           `json:"activity_type" yaml:"activity_type"`
	AccessPermission string                 `json:"access_permission,omitempty" yaml:"access_permission,omitempty"`
	Condition        string                 `json:"condition,omitempty" yaml:"condition,omitempty"`
	Arguments    []string               `json:"arguments,omitempty" yaml:"arguments,omitempty"`
	Config       config.Config          `json:"config,omitempty" yaml:"config,omitempty"`
	Result       utils.StringsMap       `json:"result,omitempty" yaml:"result,omitempty"`
	Timeout      int                    `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Retry        *RetryPolicy           `json:"retry,omitempty" yaml:"retry,omitempty"`
	HumanTask    *HumanTaskConfig       `json:"human_task,omitempty" yaml:"human_task,omitempty"`
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

// ActivityExecutor defines the signature for executing an activity logic
type ActivityExecutor func(ctx RequestContext, activity Activity, params utils.StringMap) (interface{}, error)
