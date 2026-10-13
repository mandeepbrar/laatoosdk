package elements

import (
	"laatoo.io/sdk/server/components/job"
	"laatoo.io/sdk/server/core"
)

// JobTemplate is the server's handle on a declared background job template.
//
// A job template is DECLARED, not built: a plugin's registry/jobs/<name>.yml
// states that the job template exists and its target type, target, parameters, timeout, and retry limit.
//
// Its ADDRESS is what the element exists for. Job template identity is inherited:
// a namespace may reference a template an enclosing namespace declared. Asking which declaration
// a reference bound to is answered by asking the element where it lives (e.g. "::myapp::jobmanager::<name>").
type JobTemplate interface {
	core.ServerElement

	// Template returns the job template declaration.
	Template() job.JobTemplate
}
