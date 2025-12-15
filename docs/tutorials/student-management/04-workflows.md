# Chapter 4: Workflows

Automate exam result processing using Laatoo's workflow engine.

## Overview

Workflows orchestrate complex, multi-step business processes. In this chapter, we'll create:
- **ExamProcessingWorkflow**: Automatically grade submitted exams and generate results
- **NotificationWorkflow**: Send result notifications to students

Workflows in Laatoo:
- Are defined in YAML
- Can call activities (services)
- Support parallel execution
- Can wait for external events
- Handle errors gracefully

## Step 1: Create Workflow Directory

```bash
cd dev/plugins/studentmgmt
mkdir -p config/workflows
```

## Step 2: Define Exam Processing Workflow

Create `config/workflows/ExamProcessingWorkflow.yml`:

```yaml
# Workflow for processing exam submissions
name: ExamProcessingWorkflow
description: Grades exam and generates result

# Task queue for this workflow
taskqueues:
  - exam-processing

# Workflow variables
variables:
  attemptId: ""
  examId: ""
  studentId: ""
  score: 0
  maxScore: 0
  passed: false

# Main workflow definition
workflow:
  sequence:
    elements:
      # Step 1: Load exam attempt
      - activity:
          name: LoadAttempt
          service: ExamAttempt.GetById
          arguments:
            id: "$attemptId"
          result: attempt
          
      # Step 2: Load exam details
      - activity:
          name: LoadExam
          service: Exam.GetById
          arguments:
            id: "$examId"
          result: exam
          
      # Step 3: Grade the exam
      - activity:
          name: GradeExam
          service: GradingService
          arguments:
            attempt: "$attempt"
            exam: "$exam"
          result: gradingResult
          
      # Step 4: Calculate pass/fail
      - activity:
          name: CalculateResult
          handler: inline
          script: |
            score = gradingResult.score
            maxScore = exam.MaxScore
            passingScore = exam.PassingScore
            passed = score >= passingScore
            percentage = (score / maxScore) * 100
            
            # Determine grade
            if percentage >= 90:
              grade = "A"
            elif percentage >= 80:
              grade = "B"
            elif percentage >= 70:
              grade = "C"
            elif percentage >= 60:
              grade = "D"
            else:
              grade = "F"
          
      # Step 5: Create result record
      - activity:
          name: CreateResult
          service: Result.Create
          arguments:
            data:
              AttemptId: "$attemptId"
              StudentId: "$studentId"
              ExamId: "$examId"
              Score: "$score"
              MaxScore: "$maxScore"
              Percentage: "$percentage"
              Grade: "$grade"
              Passed: "$passed"
              ProcessedDate: "$now"
          result: resultRecord
          
      # Step 6: Update attempt status
      - activity:
          name: UpdateAttempt
          service: ExamAttempt.Update
          arguments:
            id: "$attemptId"
            data:
              Status: "Graded"
              
      # Step 7: Send notification (parallel)
      - parallel:
          branches:
            - activity:
                name: SendEmail
                service: NotificationService
                arguments:
                  type: "email"
                  studentId: "$studentId"
                  subject: "Exam Result Available"
                  templateId": "exam-result"
                  data:
                    examTitle: "$exam.Title"
                    score: "$score"
                    grade: "$grade"
                    
            - activity:
                name: LogResult
                service: AuditService
                arguments:
                  event: "exam_graded"
                  data:
                    attemptId: "$attemptId"
                    score: "$score"
                    grade: "$grade"

# Error handling
errorHandling:
  retry:
    maxAttempts: 3
    initialInterval: 1s
    backoffCoefficient: 2.0
  
  onError:
    - activity:
        name: LogError
        service: ErrorLogService
        arguments:
          workflowId: "$workflowId"
          error: "$error"
```

## Step 3: Create Grading Service

Create `config/services/GradingService.yml`:

```yaml
servicemethod: studentmgmt.GradingService.Invoke
description: Grades exam attempts
```

Create `src/server/go/gradingservice.go`:

```go
package main

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/errors"
	"laatoo.io/sdk/server/log"
)

type GradingService struct {
	core.Service
}

func (s *GradingService) Initialize(ctx core.ServerContext, conf config.Config) error {
	// Define parameters
	err := s.AddCustomObjectParam(ctx, "attempt", "Exam attempt", "ExamAttempt", false, true, false)
	if err != nil {
		return errors.WrapError(ctx, err)
	}
	
	err = s.AddCustomObjectParam(ctx, "exam", "Exam configuration", "Exam", false, true, false)
	if err != nil {
		return errors.WrapError(ctx, err)
	}
	
	return nil
}

func (s *GradingService) Invoke(ctx core.RequestContext) error {
	// Get parameters
	attempt, ok := ctx.GetParamValue("attempt")
	if !ok {
		ctx.SetResponse(core.NewErrorResponse("Attempt is required"))
		return nil
	}
	
	exam, ok := ctx.GetParamValue("exam")
	if !ok {
		ctx.SetResponse(core.NewErrorResponse("Exam is required"))
		return nil
	}
	
	attemptMap := attempt.(map[string]interface{})
	examMap := exam.(map[string]interface{})
	
	// Get answers from attempt
	answers := attemptMap["Answers"].(map[string]interface{})
	
	// Get correct answers from exam (would be stored in exam definition)
	// For this example, we'll calculate a simple score
	maxScore := int(examMap["MaxScore"].(int64))
	
	// Grade the exam
	// In a real system, you'd compare answers to correct answers
	// Here we'll use a simple mock scoring
	score := s.calculateScore(answers, maxScore)
	
	log.Info(ctx, "Exam graded", "score", score, "maxScore", maxScore)
	
	// Return grading result
	ctx.SetResponse(&core.Response{
		Status: core.StatusOK,
		Data: map[string]interface{}{
			"score":     score,
			"maxScore":  maxScore,
			"graded":    true,
			"timestamp": ctx.Now(),
		},
	})
	
	return nil
}

func (s *GradingService) calculateScore(answers map[string]interface{}, maxScore int) int {
	// Mock implementation
	// In reality, you'd compare with correct answers
	numQuestions := len(answers)
	if numQuestions == 0 {
		return 0
	}
	
	// For demo: give points per answer
	pointsPerQuestion := maxScore / numQuestions
	return pointsPerQuestion * numQuestions
}
```

## Step 4: Configure Workflow Task Queue

Edit `applications/education/isolations/school-a/config/isolation.yml`:

```yaml
name: school-a
description: School A isolation

database:
  type: postgres
  host: localhost
  port: 5432
  database: school_a
  user: laatoo
  password: laatoo

# Configure task queues for workflows
taskqueues:
  exam-processing:
    workers: 5
    maxConcurrent: 10
    pollingInterval: 1s
```

## Step 5: Update Factory

Add GradingService to `src/server/go/factory.go`:

```go
func (f *StudentMgmtFactory) CreateService(
	ctx core.ServerContext,
	serviceName string,
	conf config.Config,
) (elements.Service, error) {
	
	var service elements.Service
	
	switch serviceName {
	case "EnrollmentService":
		svc := &EnrollmentService{}
		svc.SetName(serviceName)
		service = svc
		
	case "ExamSubmissionService":
		svc := &ExamSubmissionService{}
		svc.SetName(serviceName)
		service = svc
		
	case "GradingService":  // Add this case
		svc := &GradingService{}
		svc.SetName(serviceName)
		service = svc
		
	default:
		return nil, errors.NotFound(ctx, "Service", serviceName)
	}
	
	return service, nil
}
```

## Step 6: Build and Deploy

```bash
# Build plugin
laatoo plugin build studentmgmt

# Copy to application
cp -r bin ../../applications/education/config/modules/studentmgmt

# Restart solution
cd ../../../..
laatoo solution run student-mgmt-system
```

## Step 7: Test the Workflow

The workflow is automatically triggered when ExamSubmissionService pushes a task:

```go
// From ExamSubmissionService.Invoke()
err = ctx.PushTask("exam-processing", map[string]interface{}{
    "attemptId": attemptId,
    "examId":    examId,
    "studentId": studentId,
})
```

### Complete End-to-End Test

```bash
# 1. Create a student
curl -X POST http://localhost:8080/api/enroll \
  -H "Content-Type: application/json" \
  -d '{
    "Name": "Charlie Brown",
    "Email": "charlie@example.com",
    "Grade": "10"
  }'
# Save the studentId from response

# 2. Create an exam
curl -X POST http://localhost:8080/api/Exam \
  -H "Content-Type: application/json" \
  -d '{
    "ExamId": "SCIENCE101",
    "Title": "Science Quiz",
    "Duration": 30,
    "MaxScore": 100,
    "PassingScore": 60,
    "AttemptsAllowed": 2,
    "StartDate": "2024-01-01T00:00:00Z",
    "EndDate": "2024-12-31T23:59:59Z",
    "Status": "Published"
  }'

# 3. Submit exam (triggers workflow)
curl -X POST http://localhost:8080/api/exam/submit \
  -H "Content-Type: application/json" \
  -d '{
    "StudentId": "<studentId-from-step-1>",
    "ExamId": "SCIENCE101",
    "Answers": {
      "q1": "A",
      "q2": "B",
      "q3": "C",
      "q4": "A",
      "q5": "D"
    }
  }'
# Save the attemptId from response

# 4. Wait a few seconds for workflow to complete

# 5. Check the result
curl http://localhost:8080/api/Result?filter=AttemptId eq '<attemptId>Query'
```

Expected workflow execution in logs:
```
[INFO] Task received: exam-processing
[INFO] Starting workflow: ExamProcessingWorkflow
[INFO] Activity: LoadAttempt - Started
[INFO] Activity: LoadAttempt - Completed
[INFO] Activity: LoadExam - Started
[INFO] Activity: LoadExam - Completed
[INFO] Activity: GradeExam - Started
[INFO] Exam graded - score: 80, maxScore: 100
[INFO] Activity: GradeExam - Completed
[INFO] Activity: CalculateResult - Started
[INFO] Activity: CalculateResult - Completed
[INFO] Activity: CreateResult - Started
[INFO] Activity: CreateResult - Completed
[INFO] Activity: UpdateAttempt - Started
[INFO] Activity: UpdateAttempt - Completed
[INFO] Workflow completed successfully
```

## Understanding Workflows

### Workflow Structure

```yaml
workflow:
  sequence:        # Sequential execution
    elements:
      - activity:  # Call a service
          name: ActivityName
          service: ServiceName
          arguments:
            param: value
          result: variableName
```

### Activities Types

**1. Service Activity**
```yaml
- activity:
    name: LoadStudent
    service: Student.GetById
    arguments:
      id: "$studentId"
    result: student
```

**2. Inline Script Activity**
```yaml
- activity:
    name: Calculate
    handler: inline
    script: |
      result = value1 + value2
```

**3. Workflow Activity** (Call another workflow)
```yaml
- activity:
    name: ProcessPayment
    workflow: PaymentWorkflow
    arguments:
      amount: "$total"
```

### Parallel Execution

```yaml
- parallel:
    branches:
      - activity:
          name: SendEmail
          service: EmailService
      - activity:
          name: SendSMS
          service: SMSService
```

### Conditional Logic

```yaml
- condition:
    if: "$score >= 60"
    then:
      - activity:
          name: GeneratePassCertificate
          service: CertificateService
    else:
      - activity:
          name: ScheduleRetake
          service: RetakeService
```

### Loops

```yaml
- loop:
    collection: "$students"
    item: "student"
    body:
      - activity:
          name: ProcessStudent
          service: ProcessingService
          arguments:
            studentData: "$student"
```

### Error Handling

```yaml
errorHandling:
  retry:
    maxAttempts: 3
    initialInterval: 1s
    backoffCoefficient: 2.0
    maximumInterval: 1m
    
  onError:
    - activity:
        name: NotifyAdmin
        service: AdminNotificationService
        arguments:
          error: "$error"
```

## Workflow Variables

### Predefined Variables

- `$workflowId` - Unique workflow instance ID
- `$now` - Current timestamp
- `$error` - Error information (in error handlers)
- `$user` - Current user context

### Custom Variables

```yaml
variables:
  studentId: ""
  examId: ""
  score: 0
  
workflow:
  sequence:
    elements:
      - activity:
          name: Calculate
          handler: inline
          script: |
            score = 85  # Sets the variable
```

### Accessing Variables

```yaml
# In arguments
arguments:
  id: "$studentId"
  
# In conditions
if: "$score >= 60"

# In scripts
script: |
  total = score + bonus
```

## Best Practices

### 1. Keep Activities Small

```yaml
# Good: Multiple small activities
- activity:
    name: ValidateInput
    service: ValidationService
- activity:
    name: ProcessData
    service: ProcessingService
    
# Avoid: One giant activity
``yaml

### 2. Use Meaningful Names

```yaml
# Good
- activity:
    name: CalculateStudentGrade
    service: GradingService
    
# Avoid
- activity:
    name: Activity1
    service: Service1
```

### 3. Handle Errors

```yaml
errorHandling:
  retry:
    maxAttempts: 3
  onError:
    - activity:
        name: LogError
        service: ErrorLogService
```

### 4. Use Parallel for Independent Tasks

```yaml
# These can run in parallel
- parallel:
    branches:
      - activity:
          name: SendEmail
      - activity:
          name: UpdateStats
      - activity:
          name: LogEvent
```

### 5. Document Complex Workflows

```yaml
name: ComplexWorkflow
description: |
  This workflow processes exam results in multiple steps:
  1. Loads and validates the attempt
  2. Grades the exam
  3. Calculates final score
  4. Generates and stores result
  5. Sends notifications
```

## Troubleshooting

**Issue**: Workflow not executing

**Solution**: Check task queue configuration and ensure queue name matches


**Issue**: Activity fails with "Service not found"

**Solution**: Ensure service is registered and plugin is loaded

**Issue**: Variables not accessible

**Solution**: Check variable names have `$` prefix when accessed

## Summary

You've created:
- ✅ ExamProcessingWorkflow - Automated exam grading
- ✅ GradingService - Scoring logic
- ✅ Task queue configuration
- ✅ Complete end-to-end testing

**Next**: [Chapter 5: UI Plugin](05-ui-plugin.md) - Create user interface for the system
