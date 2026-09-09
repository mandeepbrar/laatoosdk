package job

import (
	"time"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// TargetType defines the execution target of a job template.
type TargetType string

const (
	TargetTypeWorkflow TargetType = "workflow"
	TargetTypeActivity TargetType = "activity"
	TargetTypeService  TargetType = "service"
	TargetTypeScript   TargetType = "script"
)

// JobTemplate represents a reusable background job definition authored by plugins.
// It specifies WHAT to execute (workflow, activity, service, or script) and with what default parameters.
// Plugin authors define job templates in src/server/registry/jobs/<template>.yml.
type JobTemplate struct {
	Name       string                `json:"name" yaml:"name"`
	Title      string                `json:"title,omitempty" yaml:"title,omitempty"`
	TargetType TargetType            `json:"targetType" yaml:"targetType"`
	Target     string                `json:"target" yaml:"target"`
	Params     map[string]core.Param `json:"params,omitempty" yaml:"params,omitempty"`
	Timeout    time.Duration         `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	RetryLimit int                   `json:"retryLimit,omitempty" yaml:"retryLimit,omitempty"`
}

// ConcurrencyPolicy controls how tick dispatches behave when a previous run is still active.
type ConcurrencyPolicy string

const (
	// ConcurrencyForbid skips the new tick if a previous run of this job is still active.
	ConcurrencyForbid ConcurrencyPolicy = "Forbid"
	// ConcurrencyAllow permits concurrent executions with independent run IDs.
	ConcurrencyAllow ConcurrencyPolicy = "Allow"
	// ConcurrencyReplace marks the active run as cancelled and dispatches a new run.
	ConcurrencyReplace ConcurrencyPolicy = "Replace"
)

// Job represents a scheduled job instance declared in a namespace configuration.
// It specifies WHEN and IN WHAT CONTEXT a JobTemplate runs (cron expression, tenancy, user, concurrency).
// Declared in namespace configuration under config/schedules/*.yml.
type Job struct {
	Name              string            `json:"name" yaml:"name"`
	TemplateName      string            `json:"templateName" yaml:"templateName"`
	Cron              string            `json:"cron" yaml:"cron"`
	Timezone          string            `json:"timezone,omitempty" yaml:"timezone,omitempty"`
	Enabled           bool              `json:"enabled" yaml:"enabled"`
	ConcurrencyPolicy ConcurrencyPolicy `json:"concurrencyPolicy" yaml:"concurrencyPolicy"`
	Params            utils.StringMap   `json:"params,omitempty" yaml:"params,omitempty"`
	Tenant            string            `json:"tenant,omitempty" yaml:"tenant,omitempty"`
	User              string            `json:"user,omitempty" yaml:"user,omitempty"`
}

// Status represents the current lifecycle state of a job execution instance.
type Status string

const (
	StatusPending Status = "Pending"
	StatusRunning Status = "Running"
	StatusSuccess Status = "Success"
	StatusFailed  Status = "Failed"
	StatusSkipped Status = "Skipped"
)

// TriggerType indicates how a job run was initiated.
type TriggerType string

const (
	TriggerTypeCron   TriggerType = "cron"
	TriggerTypeManual TriggerType = "manual"
)

// JobRun captures the telemetry and audit history of a single execution of a Job.
// Live run leases are kept in NATS KV while persistent audit history is recorded via DataManager.
type JobRun struct {
	Id            string      `json:"id"`
	JobName       string      `json:"jobName"`
	TemplateName  string      `json:"templateName"`
	Status        Status      `json:"status"`
	TriggerType   TriggerType `json:"triggerType"`
	ScheduledTime time.Time   `json:"scheduledTime,omitempty"`
	StartTime     time.Time   `json:"startTime,omitempty"`
	EndTime       time.Time   `json:"endTime,omitempty"`
	DurationMs    int64       `json:"durationMs"`
	Error         string      `json:"error,omitempty"`
	TaskId        string      `json:"taskId,omitempty"`
	WorkflowId    string      `json:"workflowId,omitempty"`
}
