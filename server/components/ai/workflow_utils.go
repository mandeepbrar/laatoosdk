package ai

import "strings"

// NormaliseWorkflowInstanceID strips any _sub_* suffix from a workflow instance ID
// so that signals and resume calls always target the root go-workflows instance.
// It handles any nesting depth: abc_sub_A_sub_B_sub_C → abc.
func NormaliseWorkflowInstanceID(id string) string {
	if idx := strings.Index(id, "_sub_"); idx > 0 {
		return id[:idx]
	}
	return id
}
