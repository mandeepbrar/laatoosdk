package main

import (
	"fmt"

	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/elements"
	"laatoo.io/sdk/server/errors"
)

type EnrollmentService struct {
	core.Service
}

func (s *EnrollmentService) Invoke(ctx core.RequestContext) error {
	method, _ := ctx.GetStringParam("method")

	switch method {
	case "enroll":
		return s.enrollStudent(ctx)
	case "approve":
		return s.approveEnrollment(ctx)
	default:
		return errors.BadArg(ctx, "method")
	}
}

func (s *EnrollmentService) enrollStudent(ctx core.RequestContext) error {
	// Logic to enroll a student in a course
	return nil
}

func (s *EnrollmentService) approveEnrollment(ctx core.RequestContext) error {
	enrollmentId, _ := ctx.GetStringParam("enrollmentId")
	fmt.Printf("Approving enrollment: %s\n", enrollmentId)

	// 1. Update enrollment status in DB
	// ... (DB update logic)

	// 2. Send notification
	notifMgr := ctx.GetServerElement(core.ServerElementNotificationManager).(elements.NotificationManager)

	// Mock student email for example
	studentEmail := "student@example.com"
	studentName := "John Doe"

	emailNotif := &core.Notification{
		NotificationType: core.EMAIL,
		Subject:          "Enrollment Approved",
		Recipients: map[string]string{
			studentEmail: studentName,
		},
		Message: []byte("Your enrollment has been approved."),
		Mime:    "text/plain",
	}

	return notifMgr.SendNotification(ctx, emailNotif)
}
