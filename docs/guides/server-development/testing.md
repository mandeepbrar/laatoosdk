# Testing Guide

Learn how to test Laatoo services and plugins with unit tests and integration tests.

## Quick Start

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestStudentService ./src/server/services
```

## Unit Testing Services

### Basic Service Test

```go
// src/server/services/student_test.go
package services

import (
    "testing"
    "lat.io/sdk/server/interfaces"
)

func TestStudentCreate(t *testing.T) {
    // Create mock context
    ctx := &MockServerContext{
        params: map[string]interface{}{
            "Name": "Test Student",
            "Email": "test@example.com",
            "Grade": 10,
        },
    }
    
   // Create service with mock DataComponent
    service := &StudentService{
        dataComponent: &MockDataComponent{},
    }
    
    // Test
    err := service.Invoke(ctx)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    
    // Verify response
    if ctx.response == nil {
        t.Error("Expected response to be set")
    }
}
```

### Mock ServerContext

```go
// test/mocks/servercontext.go
type MockServerContext struct {
    params   map[string]interface{}
    response interface{}
    error    string
}

func (m *MockServerContext) GetParamValue(key string) interface{} {
    return m.params[key]
}

func (m *MockServerContext) SetResponse(data interface{}) {
    m.response = data
}

func (m *MockServerContext) SetError(msg string) error {
    m.error = msg
    return errors.New(msg)
}
```

### Mock DataComponent

```go
// test/mocks/datacomponent.go
type MockDataComponent struct {
    createdEntities []map[string]interface{}
}

func (m *MockDataComponent) CreateEntity(entityName string, data map[string]interface{}) (interface{}, error) {
    m.createdEntities = append(m.createdEntities, data)
    data["Id"] = "mock-id-123"
    return data, nil
}

func (m *MockDataComponent) GetEntity(entityName string, id string) (interface{}, error) {
    return map[string]interface{}{
        "Id": id,
        "Name": "Mock Student",
    }, nil
}
```

## Testing Examples

### Example 1: Test Validation

```go
func TestStudentValidation(t *testing.T) {
    tests := []struct {
        name        string
        params      map[string]interface{}
        expectError bool
    }{
        {
            name: "Valid student",
            params: map[string]interface{}{
                "Name": "John Doe",
                "Email": "john@example.com",
                "Grade": 10,
            },
            expectError: false,
        },
        {
            name: "Missing name",
            params: map[string]interface{}{
                "Email": "john@example.com",
                "Grade": 10,
            },
            expectError: true,
        },
        {
            name: "Missing email",
            params: map[string]interface{}{
                "Name": "John Doe",
                "Grade": 10,
            },
            expectError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := &MockServerContext{params: tt.params}
            service := &StudentService{dataComponent: &MockDataComponent{}}
            
            err := service.Invoke(ctx)
            
            if tt.expectError && err == nil {
                t.Error("Expected error but got none")
            }
            if !tt.expectError && err != nil {
                t.Errorf("Expected no error but got: %v", err)
            }
        })
    }
}
```

### Example 2: Test Grading Logic

```go
func TestGradeCalculation(t *testing.T) {
    tests := []struct {
        score      int
        maxScore   int
        passing    int
        expected   float64
        shouldPass bool
    }{
        {90, 100, 60, 90.0, true},
        {50, 100, 60, 50.0, false},
        {75, 100, 70, 75.0, true},
    }
    
    for _, tt := range tests {
        ctx := &MockServerContext{
            params: map[string]interface{}{
                "Score":        tt.score,
                "MaxScore":     tt.maxScore,
                "PassingScore": tt.passing,
            },
        }
        
        service := &GradingService{}
        err := service.Invoke(ctx)
        
        if err != nil {
            t.Errorf("Unexpected error: %v", err)
        }
        
        result := ctx.response.(map[string]interface{})
        percentage := result["Percentage"].(float64)
        passed := result["Passed"].(bool)
        
        if percentage != tt.expected {
            t.Errorf("Expected percentage %v, got %v", tt.expected, percentage)
        }
        
        if passed != tt.shouldPass {
            t.Errorf("Expected passed %v, got %v", tt.shouldPass, passed)
        }
    }
}
```

### Example 3: Test Data Creation

```go
func TestCreateStudent(t *testing.T) {
    mock := &MockDataComponent{}
    service := &StudentService{dataComponent: mock}
    
    ctx := &MockServerContext{
        params: map[string]interface{}{
            "Name": "Alice",
            "Email": "alice@example.com",
            "Grade": 11,
        },
    }
    
    err := service.Invoke(ctx)
    if err != nil {
        t.Fatalf("Failed to create student: %v", err)
    }
    
    // Verify entity was created
    if len(mock.createdEntities) != 1 {
        t.Errorf("Expected 1 entity created, got %d", len(mock.createdEntities))
    }
    
    created := mock.createdEntities[0]
    if created["Name"] != "Alice" {
        t.Errorf("Expected name 'Alice', got %v", created["Name"])
    }
}
```

## Integration Tests

### Database Integration

```go
func TestStudentIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Setup test database
    db := setupTestDatabase(t)
    defer cleanupTestDatabase(t, db)
    
    // Create real DataComponent with test DB
    dataComponent := createTestDataComponent(db)
    
    service := &StudentService{dataComponent: dataComponent}
    ctx := &ServerContextImpl{
        params: map[string]interface{}{
            "Name": "Integration Test Student",
            "Email": "integration@example.com",
            "Grade": 10,
        },
    }
    
    // Execute
    err := service.Invoke(ctx)
    if err != nil {
        t.Fatalf("Service failed: %v", err)
    }
    
    // Verify in database
    var count int
    db.QueryRow("SELECT COUNT(*) FROM students WHERE email = ?", 
        "integration@example.com").Scan(&count)
    
    if count != 1 {
        t.Errorf("Expected 1 student in database, got %d", count)
    }
}
```

## Test Helpers

### Test Data Factory

```go
// test/helpers/factory.go
type TestDataFactory struct{}

func (f *TestDataFactory) CreateStudent(overrides ...map[string]interface{}) map[string]interface{} {
    student := map[string]interface{}{
        "Name": "Test Student",
        "Email": fmt.Sprintf("test%d@example.com", time.Now().Unix()),
        "Grade": 10,
        "Status": "active",
    }
    
    if len(overrides) > 0 {
        for k, v := range overrides[0] {
            student[k] = v
        }
    }
    
    return student
}

func (f *TestDataFactory) CreateExam() map[string]interface{} {
    return map[string]interface{}{
        "Name": "Test Exam",
        "MaxScore": 100,
        "PassingScore": 60,
    }
}
```

### Test Assertions

```go
// test/helpers/assertions.go
func AssertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }
}

func AssertEqual(t *testing.T, expected, actual interface{}) {
    t.Helper()
    if expected != actual {
        t.Errorf("Expected %v, got %v", expected, actual)
    }
}

func AssertNotNil(t *testing.T, value interface{}) {
    t.Helper()
    if value == nil {
        t.Error("Expected non-nil value")
    }
}
```

## Running Tests

### Basic Commands

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./src/server/services

# Run specific test
go test -run TestStudentCreate ./src/server/services

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Skip Integration Tests

```bash
# Skip slow tests
go test -short ./...
```

In test code:
```go
func TestSlowOperation(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping slow test")
    }
    // Long-running test...
}
```

## Best Practices

1. **Test Naming**: Use descriptive names (`TestStudentCreate`, `TestGradeCalculation`)
2. **Table-Driven Tests**: Use test tables for multiple scenarios
3. **Test Isolation**: Each test should be independent
4. **Mock Externals**: Mock database, APIs, external services
5. **Fast Unit Tests**: Keep unit tests fast (< 100ms)
6. **Cleanup**: Always cleanup resources in `defer`
7. **Test Coverage**: Aim for > 80% coverage
8. **Edge Cases**: Test boundary conditions

## CI/CD Integration

### GitHub Actions

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run Tests
        run: go test -v -cover ./...
      
      - name: Generate Coverage
        run: go test -coverprofile=coverage.out ./...
      
      - name: Upload Coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

## See Also

- [Services Guide](services.md)
- [Request Context](request-context.md)
- [Workflows](workflows.md)
