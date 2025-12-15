package main

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/errors"
)

type StudentMgmtFactory struct {
	core.ServiceFactory
}

func (f *StudentMgmtFactory) CreateService(
	ctx core.ServerContext,
	serviceName string,
	method string,
	conf config.Config,
) (core.Service, error) {

	switch serviceName {
	case "StudentService":
		service := &StudentService{}
		return service, nil
	case "CourseService":
		service := &CourseService{}
		return service, nil
	case "EnrollmentService":
		service := &EnrollmentService{}
		return service, nil
	case "ResultsService":
		service := &ResultsService{}
		return service, nil
	default:
		return nil, errors.NotFound(ctx, serviceName)
	}
}

// Export factory for plugin system
var Factory core.ServiceFactory = &StudentMgmtFactory{}
