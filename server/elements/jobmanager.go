package elements

import (
	"laatoo.io/sdk/server/components/job"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// JobManager manages background job templates, schedules, and clustered execution.
// Obtain with ctx.GetServerElement(core.ServerElementJobManager).(elements.JobManager).
//
// Plugin registries declare reusable templates (src/server/registry/jobs/*.yml), while
// namespace configurations declare schedules (config/schedules/*.yml). Execution is
// delegated to TaskManager via the "system_jobs" JetStream queue with MsgID deduplication
// ensuring clustered single-execution per tick.
type JobManager interface {
	core.ServerElement

	// RegisterJobTemplate registers a reusable job template under its name.
	RegisterJobTemplate(ctx core.ServerContext, template job.JobTemplate) error

	// GetJobTemplate returns the registered job template element by name and whether it was found.
	GetJobTemplate(ctx core.ServerContext, name string) (JobTemplate, bool)

	// ListJobTemplates returns all registered job template elements visible in this namespace.
	ListJobTemplates(ctx core.ServerContext) []JobTemplate

	// RegisterJob registers a scheduled job instance in this namespace.
	RegisterJob(ctx core.ServerContext, j job.Job) error

	// GetJob returns a scheduled job by name and whether it was found.
	GetJob(ctx core.ServerContext, name string) (job.Job, bool)

	// ListJobs returns all scheduled jobs declared in this namespace.
	ListJobs(ctx core.ServerContext) []job.Job

	// PauseJob disables schedule evaluation for the named job.
	PauseJob(ctx core.ServerContext, name string) error

	// ResumeJob re-enables schedule evaluation for the named job.
	ResumeJob(ctx core.ServerContext, name string) error

	// TriggerJob dispatches an immediate execution of the named job with optional parameter overrides,
	// bypassing the cron schedule and recording TriggerType: manual. Returns the runId.
	TriggerJob(ctx core.ServerContext, name string, overrideParams utils.StringMap) (string, error)

	// GetJobRun returns telemetry for a specific execution by run ID.
	GetJobRun(ctx core.ServerContext, runId string) (job.JobRun, bool)

	// ListJobRuns returns historical execution records for a job up to limit, ordered newest first.
	ListJobRuns(ctx core.ServerContext, jobName string, limit int) ([]job.JobRun, error)
}
