package elements_test

import (
	"testing"

	"laatoo.io/sdk/server/components/job"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/elements"
	"laatoo.io/sdk/utils"
)

// mockJobManager validates interface satisfaction at compile time.
type mockJobManager struct {
	core.ServerElement
}

func (m *mockJobManager) RegisterJobTemplate(ctx core.ServerContext, template job.JobTemplate) error {
	return nil
}
func (m *mockJobManager) GetJobTemplate(ctx core.ServerContext, name string) (job.JobTemplate, bool) {
	return job.JobTemplate{}, false
}
func (m *mockJobManager) ListJobTemplates(ctx core.ServerContext) []job.JobTemplate {
	return nil
}
func (m *mockJobManager) RegisterJob(ctx core.ServerContext, j job.Job) error {
	return nil
}
func (m *mockJobManager) GetJob(ctx core.ServerContext, name string) (job.Job, bool) {
	return job.Job{}, false
}
func (m *mockJobManager) ListJobs(ctx core.ServerContext) []job.Job {
	return nil
}
func (m *mockJobManager) PauseJob(ctx core.ServerContext, name string) error {
	return nil
}
func (m *mockJobManager) ResumeJob(ctx core.ServerContext, name string) error {
	return nil
}
func (m *mockJobManager) TriggerJob(ctx core.ServerContext, name string, overrideParams utils.StringMap) (string, error) {
	return "run-1", nil
}
func (m *mockJobManager) GetJobRun(ctx core.ServerContext, runId string) (job.JobRun, bool) {
	return job.JobRun{}, false
}
func (m *mockJobManager) ListJobRuns(ctx core.ServerContext, jobName string, limit int) ([]job.JobRun, error) {
	return nil, nil
}

func TestJobManager_InterfaceSatisfaction(t *testing.T) {
	var _ elements.JobManager = (*mockJobManager)(nil)
}
