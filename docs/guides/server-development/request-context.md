# Server and Request Context

Understanding the difference between ServerContext and RequestContext is fundamental to building Laatoo services.

## Two Different Contexts

Laatoo uses two distinct context types for different purposes:

### ServerContext
Used during **service lifecycle** (initialization, startup, shutdown)

### RequestContext  
Used during **request handling** (processing individual user requests)

## ServerContext

**When**: Service initialization, startup, and configuration

**Purpose**: Access server infrastructure and setup phase

### Common Methods

```go
// Access server elements (managers)
secHandler := ctx.GetServerElement(core.ServerElementSecurityHandler)
dataManager := ctx.GetServerElement(core.ServerElementDataManager)

// Get other services
userService, err := ctx.GetService("UserDataService")

// Access configuration
userObject, ok := service.GetStringConfiguration(ctx, "userobject")

// Create objects via ObjectLoader
obj, err := ctx.CreateObject("User")
```

### Real Example: Initialize Method

From Student Management service:

```go
func (s *StudentService) Initialize(ctx core.ServerContext, conf config.Config) error {
    // Get DataManager from server
    dataManager := ctx.GetServerElement(core.ServerElementDataManager)
    if dataManager == nil {
        return errors.ThrowError(ctx, "Data manager not found")
    }
    
    // Store for later use
    s.dataComponent = dataManager.(data.DataComponent)
    
    // Get configuration
    maxStudents, ok := s.GetIntConfiguration(ctx, "maxstudents")
    if ok {
        s.maxStudents = maxStudents
    }
    
    return nil
}
```

### Real Example: Start Method

```go
func (s *StudentService) Start(ctx core.ServerContext) error {
    // Get grade calculation service
    gradeSvcName, _ := s.GetConfiguration(ctx, "gradeservice")
    
    // Retrieve the service using ServerContext
    gradeService, err := ctx.GetService(gradeSvcName.(string))
    if err != nil {
        return errors.RethrowError(ctx, "Grade service not found", err)
    }
    
    // Type assert and store
    s.gradeService = gradeService.(GradeService)
    
    log.Info(ctx, "Student service started successfully")
    return nil
}
```

## RequestContext

**When**: Handling individual requests in the `Invoke` method

**Purpose**: Process request data and generate responses

### Common Methods

```go
// Get request parameters
credentials, _ := ctx.GetParamValue("credentials")
userId, _ := ctx.GetParamValue("userId")

// Set response
ctx.SetResponse(core.SuccessResponse(data))
ctx.SetResponse(core.StatusUnauthorizedResponse)

// Get user information
user := ctx.GetUser()
tenant := ctx.GetTenant()

// Session operations
ctx.PushTask(queueName, data)
ctx.GetFromSession("key")
ctx.SetInSession("key", value)

// Call other services
response, err := ctx.Invoke("ServiceAlias", params)
```

### Real Example: Invoke Method

From Student Management grading service:

```go
func (s *GradingService) Invoke(ctx core.RequestContext) error {
    // Get parameters from request
    studentId := ctx.GetParamValue("studentId")
    examId := ctx.GetParamValue("examId")
    score := ctx.GetParamValue("Score")
    
    if studentId == nil || examId == nil {
        ctx.SetResponse(core.BadRequestResponse("Student ID and Exam ID required"))
        return nil
    }
    
    // Get exam details using stored DataComponent
    exam, err := s.dataComponent.GetEntity("studentmgmt.Exam", examId.(string))
    if err != nil {
        log.Error(ctx, "Failed to get exam", "err", err)
        ctx.SetResponse(core.InternalErrorResponse("Failed to fetch exam"))
        return nil
    }
    
    // Calculate grade
    maxScore := exam["MaxScore"].(int)
    percentage := (float64(score.(int)) / float64(maxScore)) * 100
    passed := percentage >= float64(exam["PassingScore"].(int))
    
    // Create result
    resultData := map[string]interface{}{
        "StudentId":  studentId,
        "ExamId":     examId,
        "Score":      score,
        "Percentage": percentage,
        "Passed":     passed,
    }
    
    // Save result
    result, err := s.dataComponent.CreateEntity("studentmgmt.Result", resultData)
    if err != nil {
        ctx.SetResponse(core.InternalErrorResponse("Failed to save result"))
        return nil
    }
    
    // Return success response
    ctx.SetResponse(core.SuccessResponse(result))
    return nil
}
```

## Complete Service Example

Here's a complete service showing both contexts:

```go
package main

import (
    "laatoo.io/sdk/config"
    "laatoo.io/sdk/server/core"
    "laatoo.io/sdk/server/components/data"
    "laatoo.io/sdk/server/errors"
)

type StudentService struct {
    core.Service
    dataComponent data.DataComponent
}

// Initialize - Uses ServerContext
func (s *StudentService) Initialize(ctx core.ServerContext, conf config.Config) error {
    // Access server elements
    dataManager := ctx.GetServerElement(core.ServerElementDataManager)
    
    // Store for later use
    s.dataComponent = dataManager.(data.DataComponent)
    
    return nil
}

// Start - Uses ServerContext
func (s *StudentService) Start(ctx core.ServerContext) error {
    // Optional: Get dependencies from other services
    // otherSvc, err := ctx.GetService("OtherServiceAlias")
    
    return nil
}

// Invoke - Uses RequestContext
func (s *StudentService) Invoke(ctx core.RequestContext) error {
    // Get request parameters
    name := ctx.GetParamValue("Name")
    email := ctx.GetParamValue("Email")
    grade := ctx.GetParamValue("Grade")
    
    // Validate
    if name == nil || email == nil {
        ctx.SetResponse(core.BadRequestResponse("Name and Email required"))
        return nil
    }
    
    // Create student data
    studentData := map[string]interface{}{
        "Name": name,
        "Email": email,
        "Grade": grade,
    }
    
    // Use DataComponent (stored during Initialize)
    result, err := s.dataComponent.CreateEntity("studentmgmt.Student", studentData)
    if err != nil {
        return errors.WrapError(ctx, err)
    }
    
    // Set response
    ctx.SetResponse(core.SuccessResponse(result))
    return nil
}

// Stop - Uses ServerContext
func (s *StudentService) Stop(ctx core.ServerContext) error {
    // Cleanup resources
    return nil
}
```

## Key Differences

| Aspect | ServerContext | RequestContext |
|--------|---------------|----------------|
| **When** | Initialize, Start, Stop | Invoke |
| **Purpose** | Setup & configuration | Handle requests |
| **Lifespan** | Service lifetime | Single request |
| **Access** | Server elements, services | Request data, user info |
| **Operations** | GetService, GetServerElement | GetParamValue, SetResponse |

## Common Patterns

### Pattern 1: Store Dependencies in Initialize

```go
type MyService struct {
    core.Service
    dataComponent data.DataComponent
    otherService  SomeService
}

func (s *MyService) Initialize(ctx core.ServerContext, conf config.Config) error {
    // Get and store dependencies
    s.dataComponent = ctx.GetServerElement(core.ServerElementDataManager).(data.DataComponent)
    
    svc, _ := ctx.GetService("OtherService")
    s.otherService = svc.(SomeService)
    
    return nil
}

func (s *MyService) Invoke(ctx core.RequestContext) error {
    // Use stored dependencies
    data, err := s.dataComponent.GetEntity("Student", "123")
    result := s.otherService.Process(data)
    
    ctx.SetResponse(core.SuccessResponse(result))
    return nil
}
```

### Pattern 2: Parameter Validation

```go
func (s *ExamService) Invoke(ctx core.RequestContext) error {
    // Get all required parameters
    examId := ctx.GetParamValue("examId")
    score := ctx.GetParamValue("Score")
    maxScore := ctx.GetParamValue("MaxScore")
    
    // Validate
    if examId == nil {
        ctx.SetResponse(core.BadRequestResponse("Exam ID required"))
        return nil
    }
    
    // Continue with logic...
}
```

### Pattern 3: Error Responses

```go
func (s *Service) Invoke(ctx core.RequestContext) error {
    data, err := s.fetchData()
    
    if err != nil {
        // Log and return error response
        log.Error(ctx, "Failed to fetch data", "err", err)
        ctx.SetResponse(core.InternalErrorResponse("Data fetch failed"))
        return nil
    }
    
    ctx.SetResponse(core.SuccessResponse(data))
    return nil
}
```

## Response Types

```go
// Success
ctx.SetResponse(core.SuccessResponse(data))

// Standard status responses
ctx.SetResponse(core.StatusUnauthorizedResponse)
ctx.SetResponse(core.StatusNotFoundResponse)

// Custom status with data
ctx.SetResponse(core.NewServiceResponse(core.StatusMoreInfo, data))

// Error responses
ctx.SetResponse(core.BadRequestResponse("Invalid input"))
ctx.SetResponse(core.InternalErrorResponse("Server error"))
```

## Common Mistakes

### ❌ Using RequestContext in Initialize

```go
// WRONG - Initialize receives ServerContext, not RequestContext
func (s *Service) Initialize(ctx core.RequestContext, conf config.Config) error {
    // This won't compile!
}
```

### ❌ Using ServerContext in Invoke

```go
// WRONG - Invoke receives RequestContext, not ServerContext
func (s *Service) Invoke(ctx core.ServerContext) error {
    // This won't compile!
}
```

### ❌ Trying to get request params in Initialize

```go
// WRONG - No request during initialization
func (s *Service) Initialize(ctx core.ServerContext, conf config.Config) error {
    userId := ctx.GetParamValue("userId") // Won't work!
}
```

### ✅ Correct Pattern

```go
// CORRECT - Use ServerContext for setup
func (s *Service) Initialize(ctx core.ServerContext, conf config.Config) error {
    s.dataComponent = ctx.GetServerElement(core.ServerElementDataManager)
    return nil
}

// CORRECT - Use RequestContext for requests
func (s *Service) Invoke(ctx core.RequestContext) error {
    userId := ctx.GetParamValue("userId")
    // Process request...
    return nil
}
```

## See Also

- [Services Guide](services.md)
- [Data Components](../ui-development/datasets.md)
- [Testing](testing.md)
