# Sending Notifications

Laatoo provides a unified, asynchronous notification system for sending messages via Email and In-App channels. This guide explains how to configure and use notifications in your application.

## Architecture Overview

The notification system uses an asynchronous task-based architecture:

1. **Notification Manager**: Central entry point for sending notifications
2. **Task Queues**: Dedicated queues for each notification type (Email, In-App)
3. **Notification Plugins**: Background processors that deliver notifications

This design ensures your application isn't blocked by slow operations like sending emails.

## Configuration

### Email Notifications

Install and configure the `email` plugin to enable SMTP-based email delivery.

**File**: `config/server/modules/email.yml`

```yaml
plugin: email
settings:
  mailserver: smtp.gmail.com
  mailport: '587'
  mailsender: notifications@studentportal.com
  mailpass: "your-secure-password"
  emailqueue: emails
```

**Configuration options:**
- `mailserver`: SMTP server hostname
- `mailport`: SMTP port (`587` for TLS, `465` for SSL)
- `mailsender`: "From" email address
- `mailpass`: Sender account password (use secrets manager in production)
- `emailqueue`: Task queue name (default: `emails`)

> **Security**: Never commit passwords to version control. Use environment variables or a secrets manager.

### In-App Notifications

Configure the `inappnotifications` plugin to persist notifications to the database.

**File**: `config/server/modules/inappnotifications.yml`

```yaml
plugin: inappnotifications
settings:
  queue: inappnotifications
  dataservicefactory: database
```

**Configuration options:**
- `queue`: Task queue name (default: `inappnotifications`)
- `dataservicefactory`: Database factory for storing notification history

## Sending Notifications

### Getting the Notification Manager

Retrieve the NotificationManager from the server context in your service code.

```go
import (
    "laatoo.io/sdk/server/core"
    "laatoo.io/sdk/server/elements"
)

func SendNotification(ctx core.ServerContext) error {
    notifMgr := ctx.GetServerElement(core.ServerElementNotificationManager).(elements.NotificationManager)
    // ...use notification manager
    return nil
}
```

### Sending Email

Create a `core.Notification` with type `core.EMAIL` to send email notifications.

```go
emailNotif := &core.Notification{
    NotificationType: core.EMAIL,
    Subject:          "Welcome to Student Portal",
    Recipients: map[string]string{
        "student@example.com": "John Doe",
    },
    Message: []byte("<h1>Welcome!</h1><p>Your account is now active.</p>"),
    Mime:    "text/html",
}

// Send notification (returns immediately)
err := notifMgr.SendNotification(ctx, emailNotif)
```

### Sending In-App Notification

Create a `core.Notification` with type `core.INAPP` for in-app notifications.

```go
inAppNotif := &core.Notification{
    NotificationType: core.INAPP,
    Subject:          "New Grade Posted",
    Recipients: map[string]string{
        "user-id-123": "Student Name",
    },
    Message: []byte("Your grade for Mathematics has been posted."),
    Info: map[string]interface{}{
        "courseId": "course-101",
        "url":      "/grades/course-101",
    },
}

err := notifMgr.SendNotification(ctx, inAppNotif)
```

**Key differences for In-App:**
- `Recipients` keys are **User IDs**, not email addresses
- `Info` provides additional metadata for UI rendering
- Notifications are persisted to database
- Real-time push to connected clients

## Student Management Examples

### Example 1: Notify Students When Results Published

Send email notification when teacher publishes grades.

```go
func (s *ResultsService) PublishGrades(ctx core.RequestContext) error {
    courseId, _ := ctx.GetStringParam("courseId")
    semester, _ := ctx.GetStringParam("semester")
    
    // Get all students in the course
    students, err := s.GetEnrolledStudents(ctx, courseId)
    if err != nil {
        return err
    }
    
    // Get notification manager
    notifMgr := ctx.GetServerElement(core.ServerElementNotificationManager).(elements.NotificationManager)
    
    // Build recipient map
    recipients := make(map[string]string)
    for _, student := range students {
        recipients[student.Email] = student.Name
    }
    
    // Create email notification
    emailNotif := &core.Notification{
        NotificationType: core.EMAIL,
        Subject:          fmt.Sprintf("Grades Published - %s", semester),
        Recipients:       recipients,
        Message:          s.BuildGradeEmailHtml(courseId, semester),
        Mime:             "text/html",
    }
    
    // Send notification asynchronously
    return notifMgr.SendNotification(ctx, emailNotif)
}

func (s *ResultsService) BuildGradeEmailHtml(courseId, semester string) []byte {
    html := `
        <html>
        <body>
            <h2>Grades Published</h2>
            <p>Your grades for %s have been published.</p>
            <p>Log in to the Student Portal to view your results.</p>
            <a href="https://portal.example.com/grades">View Grades</a>
        </body>
        </html>
    `
    return []byte(fmt.Sprintf(html, semester))
}
```

### Example 2: In-App Notification for Grade Updates

Send real-time in-app notification when a grade is updated.

```go
func (s *GradesService) UpdateGrade(ctx core.RequestContext) error {
    gradeId, _ := ctx.GetStringParam("gradeId")
    newScore, _ := ctx.GetFloatParam("score")
    
    // Update the grade in database
    grade, err := s.UpdateGradeScore(ctx, gradeId, newScore)
    if err != nil {
        return err
    }
    
    // Send in-app notification to student
    notifMgr := ctx.GetServerElement(core.ServerElementNotificationManager).(elements.NotificationManager)
    
    inAppNotif := &core.Notification{
        NotificationType: core.INAPP,
        Subject:          "Grade Updated",
        Recipients: map[string]string{
            grade.Student.Id: grade.Student.Name,
        },
        Message: []byte(fmt.Sprintf(
            "Your grade for %s has been updated to %s (%.1f%%)",
            grade.CourseName,
            grade.LetterGrade,
            grade.Score,
        )),
        Info: map[string]interface{}{
            "gradeId":    gradeId,
            "courseId":   grade.CourseId,
            "courseName": grade.CourseName,
            "url":        fmt.Sprintf("/grades/%s", gradeId),
        },
    }
    
    return notifMgr.SendNotification(ctx, inAppNotif)
}
```

### Example 3: Enrollment Approval Notification

Send email when a student's course enrollment is approved.

```go
func (s *EnrollmentService) ApproveEnrollment(ctx core.RequestContext) error {
    enrollmentId, _ := ctx.GetStringParam("enrollmentId")
    
    // Approve enrollment in database
    enrollment, err := s.ApproveInDatabase(ctx, enrollmentId)
    if err != nil {
        return err
    }
    
    // Send email notification
    notifMgr := ctx.GetServerElement(core.ServerElementNotificationManager).(elements.NotificationManager)
    
    emailNotif := &core.Notification{
        NotificationType: core.EMAIL,
        Subject:          "Course Enrollment Approved",
        Recipients: map[string]string{
            enrollment.Student.Email: enrollment.Student.Name,
        },
        Message: []byte(fmt.Sprintf(`
            <html>
            <body>
                <h2>Enrollment Approved</h2>
                <p>Dear %s,</p>
                <p>Your enrollment in <strong>%s</strong> has been approved.</p>
                <p><strong>Course Details:</strong></p>
                <ul>
                    <li>Course Code: %s</li>
                    <li>Credits: %d</li>
                    <li>Instructor: %s</li>
                </ul>
                <p>Classes begin on %s.</p>
            </body>
            </html>
        `,
            enrollment.Student.Name,
            enrollment.Course.Name,
            enrollment.Course.Code,
            enrollment.Course.Credits,
            enrollment.Course.Instructor,
            enrollment.Course.StartDate.Format("January 2, 2006"),
        )),
        Mime: "text/html",
    }
    
    return notifMgr.SendNotification(ctx, emailNotif)
}
```

### Example 4: Bulk Notification for Semester Results

Send notifications to all students when semester results are published.

```go
func (s *SemesterService) PublishAllResults(ctx core.RequestContext) error {
    semester, _ := ctx.GetStringParam("semester")
    
    // Get all students with results this semester
    students, err := s.GetStudentsWithResults(ctx, semester)
    if err != nil {
        return err
    }
    
    notifMgr := ctx.GetServerElement(core.ServerElementNotificationManager).(elements.NotificationManager)
    
    // Send both email and in-app notifications
    for _, student := range students {
        // Email notification
        emailNotif := &core.Notification{
            NotificationType: core.EMAIL,
            Subject:          fmt.Sprintf("Semester Results - %s", semester),
            Recipients: map[string]string{
                student.Email: student.Name,
            },
            Message: s.BuildResultsEmail(student, semester),
            Mime:    "text/html",
        }
        
        // In-app notification
        inAppNotif := &core.Notification{
            NotificationType: core.INAPP,
            Subject:          "Semester Results Available",
            Recipients: map[string]string{
                student.Id: student.Name,
            },
            Message: []byte(fmt.Sprintf(
                "Your results for %s are now available. GPA: %.2f",
                semester,
                student.GPA,
            )),
            Info: map[string]interface{}{
                "semester": semester,
                "gpa":      student.GPA,
                "url":      fmt.Sprintf("/results/%s", semester),
            },
        }
        
        // Send both notifications
        notifMgr.SendNotification(ctx, emailNotif)
        notifMgr.SendNotification(ctx, inAppNotif)
    }
    
    return nil
}
```

## How It Works

### Asynchronous Processing

When you call `SendNotification`, the system:

1. Serializes the notification object
2. Pushes it as a task to the configured queue
3. Returns immediately (non-blocking)
4. Background plugin picks up the task
5. Plugin delivers the notification

This ensures your service responds quickly without waiting for email servers.

### In-App Notification Storage

In-App notifications are persistent:

1. Saved to database via `dataservicefactory`
2. Checked if user has active session
3. Pushed to client if connected (real-time)
4. Available in history when user logs in later

## Best Practices

### Email Notifications

1. **Use HTML Templates**: Create reusable email templates
2. **Include Plain Text**: Provide plain text alternative
3. **Personalize Content**: Use recipient name and context
4. **Add Unsubscribe**: Include unsubscribe link for bulk emails
5. **Monitor Delivery**: Log notification send events

### In-App Notifications

1. **Provide Context**: Include `Info` object with relevant data
2. **Add URLs**: Link to relevant pages for quick navigation
3. **Keep Brief**: Notification messages should be concise
4. **Use Icons**: Provide icon hints in `Info` for UI rendering
5. **Mark as Read**: Implement read/unread tracking

### Security

1. **Validate Recipients**: Ensure only authorized users receive notifications
2. **Sanitize Content**: Prevent XSS in email HTML
3. **Use Secrets Manager**: Never hardcode email passwords
4. **Rate Limiting**: Prevent notification spam
5. **Audit Trail**: Log all notification activities

### Performance

1. **Batch Notifications**: Group similar notifications together
2. **Queue Management**: Monitor queue health and processing rates
3. **Error Handling**: Implement retry logic for failed sends
4. **Cache Templates**: Cache frequently used email templates

## Error Handling

Monitor notification errors and implement retry logic.

```go
err := notifMgr.SendNotification(ctx, emailNotif)
if err != nil {
    // Log error
    log.Error(ctx, "Failed to send notification", "error", err)
    
    // Could implement retry or fallback logic
    return err
}
```

## See Also

- [Services Guide](services.md) - Implement notification logic in services
- [Workflows Guide](workflows.md) - Trigger notifications from workflows
- [Tasks Guide](tasks.md) - Understand task queue mechanism
