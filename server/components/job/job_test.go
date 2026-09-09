package job_test

import (
	"testing"
	"time"

	"laatoo.io/sdk/server/components/job"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

func TestJobComponents_Contracts(t *testing.T) {
	tmpl := job.JobTemplate{
		Name:       "cleanup",
		Title:      "Nightly Cleanup",
		TargetType: job.TargetTypeActivity,
		Target:     "system.maintenance.cleanup",
		Timeout:    10 * time.Minute,
		RetryLimit: 3,
	}
	if tmpl.Name != "cleanup" || tmpl.TargetType != "activity" {
		t.Errorf("JobTemplate unexpected field values: %+v", tmpl)
	}

	scriptTmpl := job.JobTemplate{
		Name:       "script_job",
		TargetType: job.TargetTypeScript,
		Target:     "maintenance.purge_old_records",
	}
	if scriptTmpl.TargetType != "script" {
		t.Errorf("expected TargetType == script, got %s", scriptTmpl.TargetType)
	}

	j := job.Job{
		Name:              "nightly_cleanup",
		TemplateName:      "cleanup",
		Cron:              "0 2 * * *",
		Timezone:          "UTC",
		Enabled:           true,
		ConcurrencyPolicy: job.ConcurrencyForbid,
		Params:            utils.StringMap{"retentionDays": "30"},
		Tenant:            "system",
	}
	if j.ConcurrencyPolicy != job.ConcurrencyForbid || !j.Enabled {
		t.Errorf("Job unexpected field values: %+v", j)
	}

	now := time.Now()
	run := job.JobRun{
		Id:            "run-123",
		JobName:       "nightly_cleanup",
		TemplateName:  "cleanup",
		Status:        job.StatusSuccess,
		TriggerType:   job.TriggerTypeCron,
		ScheduledTime: now,
		StartTime:     now,
		EndTime:       now.Add(5 * time.Second),
		DurationMs:    5000,
	}
	if run.Status != job.StatusSuccess || run.TriggerType != job.TriggerTypeCron {
		t.Errorf("JobRun unexpected field values: %+v", run)
	}
}

func TestServerElementJobManager_ConstantValue(t *testing.T) {
	if core.ServerElementJobManager != 42 {
		t.Fatalf("expected ServerElementJobManager = 42, got %d", core.ServerElementJobManager)
	}
}
