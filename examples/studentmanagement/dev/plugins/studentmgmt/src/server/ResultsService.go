package main

import (
	"fmt"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/elements"
	"laatoo.io/sdk/server/errors"
)

type ResultsService struct {
	core.Service
}

func (s *ResultsService) Invoke(ctx core.RequestContext) error {
	method, _ := ctx.GetStringParam("method")

	switch method {
	case "publish":
		return s.publishGrades(ctx)
	default:
		return errors.BadArg(ctx, "method")
	}
}

func (s *ResultsService) publishGrades(ctx core.RequestContext) error {
	courseId, _ := ctx.GetStringParam("courseId")

	// 1. Logic to publish grades for the course
	// ...

	// 2. Send notifications to students
	notifMgr := ctx.GetServerElement(core.ServerElementNotificationManager).(elements.NotificationManager)

	// Mock recipients
	recipients := map[string]string{
		"student1@example.com": "Student One",
		"student2@example.com": "Student Two",
	}

	emailNotif := &core.Notification{
		NotificationType: core.EMAIL,
		Subject:          "Grades Published",
		Recipients:       recipients,
		Message:          []byte(fmt.Sprintf("Grades for course %s have been published.", courseId)),
		Mime:             "text/plain",
	}

	return notifMgr.SendNotification(ctx, emailNotif)
}
