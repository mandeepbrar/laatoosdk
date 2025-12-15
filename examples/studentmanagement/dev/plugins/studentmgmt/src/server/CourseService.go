package main

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/errors"
)

type CourseService struct {
	core.Service
}

func (s *CourseService) Invoke(ctx core.RequestContext) error {
	method, _ := ctx.GetStringParam("method")

	switch method {
	case "create":
		return s.createCourse(ctx)
	case "list":
		return s.listCourses(ctx)
	default:
		return errors.BadArg(ctx, "method")
	}
}

func (s *CourseService) createCourse(ctx core.RequestContext) error {
	// Logic to create a new course
	return nil
}

func (s *CourseService) listCourses(ctx core.RequestContext) error {
	// Logic to list available courses
	return nil
}
