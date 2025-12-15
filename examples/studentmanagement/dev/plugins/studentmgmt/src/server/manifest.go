package main

import (
	"laatoo.io/sdk/server/core"
)

func Manifest(provider core.MetaDataProvider) []core.PluginComponent {
	return []core.PluginComponent{
		core.PluginComponent{Object: StudentMgmtFactory{}},
		core.PluginComponent{Object: StudentService{}},
		core.PluginComponent{Object: CourseService{}},
		core.PluginComponent{Object: EnrollmentService{}},
		core.PluginComponent{Object: ResultsService{}},
	}
}
