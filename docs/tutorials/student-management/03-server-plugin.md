# Chapter 3: Server Plugin - Services

Create custom services for Student Management System business logic.

## Overview

While entities provide automatic CRUD operations, you'll need custom services for:
- Complex business logic
- Multi-step operations
- Validation beyond simple constraints
- Integration with other services

In this chapter, we'll create:
- **EnrollmentService**: Register students with validation
- **ExamSchedulingService**: Schedule and manage exams
- **ExamSubmissionService**: Handle exam attempts
- **ResultProcessingService**: Calculate and store results

## Step 1: Create Service Definitions Directory

```bash
cd dev/plugins/studentmgmt
mkdir -p config/services
mkdir -p src/server/go
```

## Step 2: Create Enrollment Service

### Service Configuration

Create `config/services/EnrollmentService.yml`:

```yaml
# Service definition for student enrollment
servicemethod: studentmgmt.EnrollmentService.Invoke
description: Handles student registration and enrollment

# Service will need these privileges
permissions:
  - create:student
  - update:student
```

### Service Implementation

Create `src/server/go/enrollmentservice.go`:

```go
package main

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/components/data"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/elements"
	"laatoo.io/sdk/server/errors"
	"laatoo.io/sdk/server/log"
	"time"
	"github.com/twinj/uuid"
)

type EnrollmentService struct {
	core.Service
	dataManager elements.DataManager
	studentData data.DataComponent
}

func (s *EnrollmentService) Initialize(ctx core.ServerContext, conf config.Config) error {
	// Get DataManager from server context
	s.dataManager = ctx.GetServerElement(core.ServerElementDataManager).(elements.DataManager)
	
	// Define service parameters
	err := s.AddStringParam(ctx, "Name", "Student's full name")
	if err != nil {
		return errors.WrapError(ctx, err)
	}
	
	err = s.AddStringParam(ctx, "Email", "Student's email address")
	if err != nil {
		return errors.WrapError(ctx, err)
	}
	
	err = s.AddStringParam(ctx, "Grade", "Student's grade level")
	if err != nil {
		return errors.WrapError(ctx, err)
	}
	
	return nil
}

func (s *EnrollmentService) Start(ctx core.ServerContext) error {
	// Get Student data component
	studentSvc, err := ctx.GetService("Student.GetById")
	if err != nil {
		return errors.WrapError(ctx, err)
	}
	
	s.studentData = studentSvc.(data.DataComponent)
	return nil
}

func (s *EnrollmentService) Invoke(ctx core.RequestContext) error {
	// Get parameters
	name, ok := ctx.GetStringParam("Name")
	if !ok || name == "" {
		ctx.SetResponse(core.NewErrorResponse("Name is required"))
		return nil
	}
	
	email, ok := ctx.GetStringParam("Email")
	if !ok || email == "" {
		ctx.SetResponse(core.NewErrorResponse("Email is required"))
		return nil
	}
	
	grade, _ := ctx.GetStringParam("Grade")
	
	// Validate email doesn't already exist
	existing, err := s.studentData.Query(ctx, map[string]interface{}{
		"Email": email,
	}, "", 0, 1)
	
	if err != nil {
		log.Error(ctx, "Error checking existing student", "error", err)
		return err
	}
	
	if len(existing) > 0 {
		ctx.SetResponse(core.NewErrorResponse("Email already registered"))
		return nil
	}
	
	// Create student record
	studentId := uuid.NewV1().String()
	student := map[string]interface{}{
		"StudentId":      studentId,
		"Name":           name,
		"Email":          email,
		"Grade":          grade,
		"EnrollmentDate": time.Now(),
		"Status":         "Active",
	}
	
	created, err := s.studentData.Create(ctx, student)
	if err != nil {
		log.Error(ctx, "Error creating student", "error", err)
		return err
	}
	
	log.Info(ctx, "Student enrolled successfully", "studentId", studentId, "email", email)
	
	// Return success response
	ctx.SetResponse(&core.Response{
		Status: core.StatusOK,
		Data: map[string]interface{}{
			"message":   "Student enrolled successfully",
			"studentId": studentId,
			"student":   created,
		},
	})
	
	return nil
}
```

## Step 3: Create Exam Submission Service

Create `config/services/ExamSubmissionService.yml`:

```yaml
servicemethod: studentmgmt.ExamSubmissionService.Invoke
description: Handles exam attempt submission
```

Create `src/server/go/examsubmissionservice.go`:

```go
package main

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/components/data"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/elements"
	"laatoo.io/sdk/server/errors"
	"laatoo.io/sdk/server/log"
	"time"
	"github.com/twinj/uuid"
)

type ExamSubmissionService struct {
	core.Service
	dataManager    elements.DataManager
	examData       data.DataComponent
	attemptData    data.DataComponent
	resultData     data.DataComponent
	taskManager    elements.TaskManager
}

func (s *ExamSubmissionService) Initialize(ctx core.ServerContext, conf config.Config) error {
	s.dataManager = ctx.GetServerElement(core.ServerElementDataManager).(elements.DataManager)
	s.taskManager = ctx.GetServerElement(core.ServerElementTaskManager).(elements.TaskManager)
	
	// Define parameters
	err := s.AddStringParam(ctx, "StudentId", "Student ID")
	if err != nil {
		return errors.WrapError(ctx, err)
	}
	
	err = s.AddStringParam(ctx, "ExamId", "Exam ID")
	if err != nil {
		return errors.WrapError(ctx, err)
	}
	
	err = s.AddCustomObjectParam(ctx, "Answers", "Exam answers", "json", false, true, false)
	if err != nil {
		return errors.WrapError(ctx, err)
	}
	
	return nil
}

func (s *ExamSubmissionService) Start(ctx core.ServerContext) error {
	// Get data components
	examSvc, _ := ctx.GetService("Exam.GetById")
	s.examData = examSvc.(data.DataComponent)
	
	attemptSvc, _ := ctx.GetService("ExamAttempt.GetById")
	s.attemptData = attemptSvc.(data.DataComponent)
	
	resultSvc, _ := ctx.GetService("Result.GetById")
	s.resultData = resultSvc.(data.DataComponent)
	
	return nil
}

func (s *ExamSubmissionService) Invoke(ctx core.RequestContext) error {
	// Get parameters
	studentId, _ := ctx.GetStringParam("StudentId")
	examId, _ := ctx.GetStringParam("ExamId")
	answers, _ := ctx.GetParamValue("Answers")
	
	// Validate exam exists
	exam, err := s.examData.GetById(ctx, examId, "")
	if err != nil {
		ctx.SetResponse(core.NewErrorResponse("Exam not found"))
		return nil
	}
	
	examMap := exam.(map[string]interface{})
	
	// Check exam is published
	if examMap["Status"] != "Published" {
		ctx.SetResponse(core.NewErrorResponse("Exam is not available"))
		return nil
	}
	
	// Check if within exam window
	now := time.Now()
	startDate := examMap["StartDate"].(time.Time)
	endDate := examMap["EndDate"].(time.Time)
	
	if now.Before(startDate) || now.After(endDate) {
		ctx.SetResponse(core.NewErrorResponse("Exam is not currently available"))
		return nil
	}
	
	// Check attempts allowed
	existingAttempts, _ := s.attemptData.Query(ctx, map[string]interface{}{
		"StudentId": studentId,
		"ExamId":    examId,
	}, "", 0, 100)
	
	attemptsAllowed := int(examMap["AttemptsAllowed"].(int64))
	if len(existingAttempts) >= attemptsAllowed {
		ctx.SetResponse(core.NewErrorResponse("Maximum attempts exceeded"))
		return nil
	}
	
	// Create exam attempt
	attemptId := uuid.NewV1().String()
	attempt := map[string]interface{}{
		"AttemptId":     attemptId,
		"StudentId":     studentId,
		"ExamId":        examId,
		"AttemptNumber": len(existingAttempts) + 1,
		"StartTime":     time.Now().Add(-time.Duration(examMap["Duration"].(int64)) * time.Minute),
		"SubmitTime":    time.Now(),
		"Status":        "Submitted",
		"Answers":       answers,
	}
	
	created, err := s.attemptData.Create(ctx, attempt)
	if err != nil {
		log.Error(ctx, "Error creating exam attempt", "error", err)
		return err
	}
	
	log.Info(ctx, "Exam submitted successfully", "attemptId", attemptId)
	
	// Queue task for result processing
	if s.taskManager != nil {
		err = ctx.PushTask("exam-processing", map[string]interface{}{
			"attemptId": attemptId,
			"examId":    examId,
			"studentId": studentId,
		})
		
		if err != nil {
			log.Error(ctx, "Error queuing result processing task", "error", err)
		}
	}
	
	// Return response
	ctx.SetResponse(&core.Response{
		Status: core.StatusOK,
		Data: map[string]interface{}{
			"message":   "Exam submitted successfully",
			"attemptId": attemptId,
			"attempt":   created,
		},
	})
	
	return nil
}
```

## Step 4: Create Service Factory

Create `src/server/go/factory.go`:

```go
package main

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/elements"
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
	
	var service core.Service
	
	switch serviceName {
	case "EnrollmentService":
		svc := &EnrollmentService{}
		service = svc
		
	case "ExamSubmissionService":
		svc := &ExamSubmissionService{}
		service = svc
		
	default:
		return nil, errors.NotFound(ctx, "Service", serviceName)
	}
	
	return service, nil
}

// Export factory for plugin system
var Factory core.ServiceFactory = &StudentMgmtFactory{}
```

## Step 5: Create Plugin Manifest

Create `src/server/go/manifest.go` to expose your factory to the server:

```go
package main

import (
	"laatoo.io/sdk/server/core"
)

func Manifest(provider core.MetaDataProvider) []core.PluginComponent {
	return []core.PluginComponent{
		core.PluginComponent{Object: StudentMgmtFactory{}},
		core.PluginComponent{Object: EnrollmentService{}},
		core.PluginComponent{Object: ExamSubmissionService{}},
	}
}
```

## Step 6: Initialize Go Module

```bash
cd src/server/go

# Initialize module
go mod init studentmgmt

# Add SDK dependency
go get laatoo.io/sdk@latest

# Add UUID library
go get github.com/twinj/uuid

# Download dependencies
go mod tidy
```

Create `go.mod` if needed:
```go
module studentmgmt

go 1.19

require (
	laatoo.io/sdk v1.0.0
	github.com/twinj/uuid v1.0.0
)
```

## Step 7: Create API Channels

Services need channels to be accessible via HTTP.

Create `config/channels/enrollment-api.yml`:

```yaml
parent: root
method: POST
service: EnrollmentService
path: /api/enroll
description: Student enrollment endpoint
```

Create `config/channels/exam-submission-api.yml`:

```yaml
parent: root
method: POST
service: ExamSubmissionService
path: /api/exam/submit
description: Exam submission endpoint
authenticate: true  # Require authentication
```

## Step 8: Build the Plugin

```bash
cd ../../..  # Back to plugin root
laatoo plugin build student mgmt
```

Verify build output:
```
[INFO] Building Go plugin...
[INFO] Compiled: bin/objects/studentmgmt-local.so
[INFO] Copied services to: bin/config/services/
[INFO] Copied channels to: bin/config/channels/
[INFO] Build complete
```

## Step 9: Install and Test

```bash
# Install plugin
cp -r bin ../../applications/education/config/modules/studentmgmt

# Restart solution
cd ../../../..
laatoo solution run student-mgmt-system
```

### Test Enrollment Service

```bash
# Enroll a student
curl -X POST http://localhost:8080/api/enroll \
  -H "Content-Type: application/json" \
  -d '{
    "Name": "Bob Smith",
    "Email": "bob@example.com",
    "Grade": "11"
  }'
```

Expected response:
```json
{
  "message": "Student enrolled successfully",
  "studentId": "...",
  "student": {
    "StudentId": "...",
    "Name": "Bob Smith",
    "Email": "bob@example.com",
    "Status": "Active"
  }
}
```

### Test Exam Submission

```bash
# First, create an exam (using auto-generated endpoint)
curl -X POST http://localhost:8080/api/Exam \
  -H "Content-Type: application/json" \
  -d '{
    "ExamId": "MATH101",
    "Title": "Math Midterm",
    "Duration": 60,
    "MaxScore": 100,
    "PassingScore": 60,
    "AttemptsAllowed": 2,
    "StartDate": "2024-01-01T00:00:00Z",
    "EndDate": "2024-12-31T23:59:59Z",
    "Status": "Published"
  }'

# Submit exam
curl -X POST http://localhost:8080/api/exam/submit \
  -H "Content-Type: application/json" \
  -d '{
    "StudentId": "...",
    "ExamId": "MATH101",
    "Answers": {
      "q1": "A",
      "q2": "B",
      "q3": "C"
    }
  }'
```

## Understanding the Code

### Service Lifecycle

```go
// 1. Initialize - Called once during server startup
func (s *Service) Initialize(ctx core.ServerContext, conf config.Config) error {
    // Get dependencies
    // Define parameters
    // Setup configuration
}

// 2. Start - Called when service starts
func (s *Service) Start(ctx core.ServerContext) error {
    // Get other services
    // Final setup
}

// 3. Invoke - Called for each request
func (s *Service) Invoke(ctx core.RequestContext) error {
    // Handle request
    // Set response
}
```

### Getting Parameters

```go
// String parameter
name, ok := ctx.GetStringParam("Name")

// Integer parameter
age, ok := ctx.GetIntParam("Age")

// Custom object parameter
data, ok := ctx.GetParamValue("Answers")
```

### Accessing Data

```go
// Get by ID
student, err := s.studentData.GetById(ctx, studentId, "")

// Query with filter
results, err := s.studentData.Query(ctx, map[string]interface{}{
    "Email": email,
}, "", 0, 10)

// Create
created, err := s.studentData.Create(ctx, studentMap)

// Update
updated, err := s.studentData.Update(ctx, studentId, updates, "")
```

### Setting Response

```go
// Success response
ctx.SetResponse(&core.Response{
    Status: core.StatusOK,
    Data:   result,
})

// Error response
ctx.SetResponse(core.NewErrorResponse("Error message"))

// Status codes
core.StatusOK           // 200
core.StatusCreated      // 201
core.StatusBadRequest   // 400
core.StatusUnauthorized // 401
core.StatusNotFound     // 404
```

## Best Practices

### 1. Validate Inputs

```go
if name == "" {
    ctx.SetResponse(core.NewErrorResponse("Name is required"))
    return nil
}
```

### 2. Check Business Rules

```go
// Check attempts allowed
if len(existingAttempts) >= attemptsAllowed {
    ctx.SetResponse(core.NewErrorResponse("Maximum attempts exceeded"))
    return nil
}
```

### 3. Log Important Events

```go
log.Info(ctx, "Student enrolled", "studentId", studentId)
log.Error(ctx, "Failed to create student", "error", err)
```

### 4. Use Transactions (when available)

```go
// For multi-step operations
tx, err := s.dataManager.BeginTransaction(ctx)
defer tx.Rollback()

// ... perform operations ...

tx.Commit()
```

### 5. Queue Background Tasks

```go
// Don't block the request
ctx.PushTask("result-processing", data)
```

## Summary

You've created:
- ✅ EnrollmentService - Student registration with validation
- ✅ ExamSubmissionService - Exam attempt submission with business rules
- ✅ Service factory for plugin system
- ✅ Plugin manifest to expose components
- ✅ API channels for HTTP access
- ✅ Tested services via REST API

**Next**: [Chapter 4: Workflows](04-workflows.md) - Automate result processing with workflows
