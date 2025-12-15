# UI Services

UI Services are HTTP endpoint definitions that provide shortcuts for calling backend services from the UI. They are defined in the UI plugin's registry.

## Purpose

UI Services enable:
- Simplified HTTP calls to backend APIs
- Centralized endpoint configuration
- Easy service discovery
- Type-safe service definitions
- Consistent URL patterns

## File Locations

Services can be defined in two ways:

### Option 1: Services.yml File

`src/ui/registry/Services.yml`:

```yaml
loadStudent:
  method: GET
  url: /studentmgmt.Student/:id

queryStudents:
  method: GET
  url: /studentmgmt.Student

createStudent:
  method: POST
  url: /studentmgmt.Student

updateStudent:
  method: PUT
  url: /studentmgmt.Student/:id

deleteStudent:
  method: DELETE
  url: /studentmgmt.Student/:id
```

### Option 2: Services Directory

`src/ui/registry/Services/*.yml`:

```
Services/
├── loadStudent.yml
├── queryStudents.yml
├── createStudent.yml
└── deleteStudent.yml
```

Each file contains:
```yaml
method: GET
url: /studentmgmt.Student/:id
```

## Service Definition

### Basic Structure

```yaml
serviceName:
  method: GET|POST|PUT|DELETE|PATCH
  url: /module.Entity/path
```

**Properties**:
- `method`: HTTP method (GET, POST, PUT, DELETE, PATCH)
- `url`: API endpoint path

### URL Parameters

Use `:paramName` syntax for path parameters:

```yaml
getStudentExams:
  method: GET
  url: /studentmgmt.Exam/student/:studentId

getExamAttempt:
  method: GET
  url: /studentmgmt.ExamAttempt/:attemptId
```

## Using Services

### From Window.executeService

```javascript
// Calls the 'loadStudent' service
Window.executeService(
  'loadStudent',
  { id: '123' },  // URL params
  null,           // Query params
  function(response) {
    console.log('Student:', response.data);
  },
  function(error) {
    console.error('Error:', error);
  }
);
```

### From Actions

```yaml
# src/ui/registry/Actions/loadStudent.yml
actiontype: executeservice
servicename: loadStudent
actionparams:
  serviceparams:
    id: '{{ctx.studentId}}'
```

### From Sagas

```javascript
import { DataSource, RequestBuilder } from 'jsui';

function* loadStudentSaga(action) {
  const req = RequestBuilder.URLParamsRequest({ id: action.payload.id });
  const response = yield call(DataSource.ExecuteService, 'loadStudent', req);
  yield put({ type: 'STUDENT_LOADED', payload: response.data });
}
```

### From Methods

```javascript
function loadStudentData(payload, actionparams) {
  Window.executeService(
    'loadStudent',
    { id: payload.studentId },
    null,
    function(response) {
      // Handle response
    }
  );
}
```

## Common Service Patterns

### 1. CRUD Operations

```yaml
# Services.yml
getEntity:
  method: GET
  url: /module.Entity/:id

queryEntities:
  method: GET
  url: /module.Entity

createEntity:
  method: POST
  url: /module.Entity

updateEntity:
  method: PUT
  url: /module.Entity/:id

deleteEntity:
  method: DELETE
  url: /module.Entity/:id
```

### 2. Custom Actions

```yaml
gradeExam:
  method: POST
  url: /studentmgmt.Exam/:examId/grade

publishExam:
  method: POST
  url: /studentmgmt.Exam/:examId/publish

submitExamAttempt:
  method: POST
  url: /studentmgmt.ExamAttempt/:attemptId/submit
```

### 3. Nested Resources

```yaml
getStudentEnrollments:
  method: GET
  url: /studentmgmt.Student/:studentId/enrollments

getExamResults:
  method: GET
  url: /studentmgmt.Exam/:examId/results
```

### 4. Query Operations

```yaml
searchStudents:
  method: GET
  url: /studentmgmt.Student/search

filterExams:
  method: GET
  url: /studentmgmt.Exam/filter
```

## Example: Student Management Services

```yaml
# src/ui/registry/Services.yml

# Student Services
loadStudent:
  method: GET
  url: /studentmgmt.Student/:id

queryStudents:
  method: GET
  url: /studentmgmt.Student

createStudent:
  method: POST
  url: /studentmgmt.Student

updateStudent:
  method: PUT
  url: /studentmgmt.Student/:id

deleteStudent:
  method: DELETE
  url: /studentmgmt.Student/:id

# Exam Services
loadExam:
  method: GET
  url: /studentmgmt.Exam/:id

queryExams:
  method: GET
  url: /studentmgmt.Exam

createExam:
  method: POST
  url: /studentmgmt.Exam

publishExam:
  method: POST
  url: /studentmgmt.Exam/:id/publish

# Exam Attempt Services
submitExamAttempt:
  method: POST
  url: /studentmgmt.ExamAttempt/submit

gradeExamAttempt:
  method: POST
  url: /studentmgmt.ExamAttempt/:id/grade

getStudentAttempts:
  method: GET
  url: /studentmgmt.ExamAttempt/student/:studentId

# Result Services
queryResults:
  method: GET
  url: /studentmgmt.Result

getStudentResults:
  method: GET
  url: /studentmgmt.Result/student/:studentId
```

## Using Services in Complete Workflow

```javascript
// src/ui/registry/Methods/submitExam.js
function submitExam(payload, actionparams, event, form) {
  if (form) {
    form.setSubmitting(true);
  }
  
  // Call submitExamAttempt service
  Window.executeService(
    'submitExamAttempt',
    {
      examId: payload.examId,
      studentId: payload.studentId,
      answers: payload.answers
    },
    null,
    function(response) {
      const attemptId = response.data.Id;
      
      // Call gradeExamAttempt service
      Window.executeService(
        'gradeExamAttempt',
        { id: attemptId },
        null,
        function(gradeResponse) {
          if (form) {
            form.setSubmitting(false);
          }
          
          Window.showMessage({ 
            Default: `Exam graded! Score: ${gradeResponse.data.Percentage}%` 
          });
          
          Window.redirectPage('/results');
        },
        function(error) {
          if (form) {
            form.setSubmitting(false);
          }
          Window.showError({ Default: 'Failed to grade exam' });
        }
      );
    },
    function(error) {
      if (form) {
        form.setSubmitting(false);
      }
      Window.showError({ Default: 'Failed to submit exam' });
    }
  );
}
```

## RequestBuilder Patterns

### URL Parameters

```javascript
const req = RequestBuilder.URLParamsRequest({ id: '123' });
// Replaces :id in URL with '123'
```

### Query Parameters

```javascript
const req = RequestBuilder.URLParamsRequest(
  { id: '123' },
  { filter: 'active', limit: 10 }
);
// URL: /entity/123?filter=active&limit=10
```

### POST/PUT Body

```javascript
const req = RequestBuilder.DefaultRequest(null, {
  Name: 'John',
  Email: 'john@example.com'
});
```

## Best Practices

1. **Naming**: Use descriptive, action-based names (loadStudent, createExam)
2. **Consistency**: Follow REST conventions (GET for read, POST for create, etc.)
3. **Organization**: Group related services together
4. **URL Patterns**: Use consistent module.Entity/action patterns
5. **Parameters**: Use path params for IDs, query params for filters
6. **Documentation**: Comment complex services

## vs Direct DataSource Calls

### Using Service Definition (Recommended)

```javascript
Window.executeService('loadStudent', { id: '123' }, null, success, error);
```

### Direct DataSource Call

```javascript
DataSource.ExecuteService(
  'studentmgmt.Student.GetById',
  RequestBuilder.DefaultRequest(null, { id: '123' })
).then(success).catch(error);
```

**Advantages of Service Definitions**:
- Centralized configuration
- Easier refactoring
- Auto-completion support
- Consistent error handling
- URL pattern reuse

## See Also

- [Window Object](window-object.md) - executeService method
- [Sagas](sagas.md) - Using services in sagas
- [Methods](methods.md) - Using services in methods
- [Actions](actions.md) - executeservice action type
