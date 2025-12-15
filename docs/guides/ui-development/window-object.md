# Window Object

The `Window` object (with capital W) is a global utility object providing essential functions for Laat

oo UI plugins.

## Available Methods

### Window.executeService

Execute a backend service:

```javascript
Window.executeService(serviceName, payload, urlParams, successCallback, errorCallback)
```

**Parameters**:
- `serviceName` (string): Service to call (e.g., 'Student.Query')
- `payload` (object): Request body data
- `urlParams` (object): URL parameters
- `successCallback` (function): Success handler, receives response
- `errorCallback` (function): Error handler, receives error

**Example**:
```javascript
Window.executeService(
  'Student.GetById',
  { id: '123' },
  null,
  function(response) {
    console.log('Student:', response.data);
    Window.showMessage({ Default: 'Loaded successfully' });
  },
  function(error) {
    console.error('Error:', error);
    Window.showError({ Default: 'Failed to load' });
  }
);
```

### Window.showMessage

Display a success/info message:

```javascript
Window.showMessage(messageObject)
```

**Parameters**:
- `messageObject` (object): Message with `Default` property

**Example**:
```javascript
Window.showMessage({ Default: 'Student enrolled successfully!' });
```

### Window.showError

Display an error message:

```javascript
Window.showError(errorObject, response)
```

**Parameters**:
- `errorObject` (object): Error with `Default` property
- `response` (optional): Response object

**Example**:
```javascript
Window.showError({ Default: 'Failed to enroll student' });
```

### Window.executeAction

Execute a registered action:

```javascript
Window.executeAction(action, payload, actionparams)
```

**Parameters**:
- `action` (string|object): Action name or action definition
- `payload` (object): Data to pass to action
- `actionparams` (object): Action parameters

**Example**:
```javascript
Window.executeAction('loadStudents', { filter: 'active' });

// With params
Window.executeAction('enrollStudent', 
  { name: 'John', email: 'john@example.com' },
  { successaction: 'navigateToDashboard' }
);
```

### Window.redirectPage

Navigate to a page:

```javascript
Window.redirectPage(pagePath, params)
```

**Parameters**:
- `pagePath` (string): Page route (e.g., '/students')
- `params` (object): URL parameters

**Example**:
```javascript
Window.redirectPage('/students');
Window.redirectPage('/student/:id', { id: '123' });
```

### Window.redirect

Navigate to a URL:

```javascript
Window.redirect(url, params, newWindow)
```

**Parameters**:
- `url` (string): Full URL
- `params` (object): URL parameters
- `newWindow` (boolean): Open in new window

**Example**:
```javascript
Window.redirect('https://example.com');
Window.redirect('/external/page', null, true); // Open in new window
```

### Window.showDialog

Display a dialog/modal:

```javascript
Window.showDialog(title, component, onClose, actions, contentStyle, titleStyle)
```

**Parameters**:
- `title` (string): Dialog title
- `component` (React component): Content component
- `onClose` (function): Close callback
- `actions` (array): Button actions
- `contentStyle` (object): Content CSS
- `titleStyle` (object): Title CSS

**Example**:
```javascript
Window.showDialog(
  'Confirm Enrollment',
  <ConfirmComponent data={studentData} />,
  () => console.log('Dialog closed'),
  [
    { label: 'Cancel', onClick: () => Window.closeDialog() },
    { label: 'Confirm', onClick: () => confirmEnrollment() }
  ]
);
```

### Window.closeDialog

Close the active dialog:

```javascript
Window.closeDialog()
```

### Window.showInteraction / Window.closeInteraction

Show/close interactions (drawers, modals):

```javascript
Window.showInteraction(interactiontype, title, component, onClose, actions, contentStyle, titleStyle)
Window.closeInteraction(interactiontype)
```

**Example**:
```javascript
Window.showInteraction('Drawer', 'Student Details', <StudentDetailsBlock />, null);
Window.closeInteraction('Drawer');
```

### Window.getActionCallback

Get a callback function for an action:

```javascript
const callback = Window.getActionCallback(action, params, actionparams)
```

**Example**:
```javascript
const handleClick = Window.getActionCallback('loadStudents', { filter: 'active' });
// Use in component: onClick={handleClick}
```

### Window.handleError

Global error handler:

```javascript
Window.handleError(errorObject, response)
```

**Example**:
```javascript
try {
  // some code
} catch (error) {
  Window.handleError(error);
}
```

## Common Usage Patterns

### 1. Load and Display Data

```javascript
function loadStudent(studentId) {
  Window.executeService(
    'Student.GetById',
    { id: studentId },
    null,
    function(response) {
      // Update UI with response.data
      displayStudent(response.data);
      Window.showMessage({ Default: 'Student loaded' });
    },
    function(error) {
      Window.showError({ Default: 'Failed to load student' });
    }
  );
}
```

### 2. Submit Form and Navigate

```javascript
function submitEnrollment(formData) {
  Window.executeService(
    'EnrollmentService',
    formData,
    null,
    function(response) {
      Window.showMessage({ Default: 'Enrollment successful!' });
      setTimeout(() => {
        Window.redirectPage('/students');
      }, 1500);
    },
    function(error) {
      Window.showError({ Default: 'Enrollment failed' });
    }
  );
}
```

### 3. Confirm Action with Dialog

```javascript
function deleteStudent(studentId) {
  Window.showDialog(
    'Confirm Delete',
    <p>Are you sure you want to delete this student?</p>,
    null,
    [
      {
        label: 'Cancel',
        onClick: () => Window.closeDialog()
      },
      {
        label: 'Delete',
        onClick: () => {
          Window.closeDialog();
          performDelete(studentId);
        }
      }
    ]
  );
}

function performDelete(studentId) {
  Window.executeService(
    'Student.Delete',
    { id: studentId },
    null,
    function(response) {
      Window.showMessage({ Default: 'Student deleted' });
      Window.executeAction('loadStudents'); // Refresh list
    },
    function(error) {
      Window.showError({ Default: 'Failed to delete student' });
    }
  );
}
```

### 4. Chain Actions

```javascript
function submitExam(examData) {
  Window.executeAction(
    'submitExam',
    examData,
    {
      successaction: 'gradeExam',
      successCallback: function(response) {
        Window.showMessage({ Default: 'Exam submitted and graded!' });
        Window.redirectPage('/results');
      }
    }
  );
}
```

### 5. Conditional Navigation

```javascript
function navigateToExam(examId, attemptsRemaining) {
 if (attemptsRemaining <= 0) {
    Window.showError({ Default: 'No attempts remaining for this exam' });
    return;
  }
  
  Window.redirectPage('/take-exam/' + examId);
}
```

## Integration with Methods

Methods can use Window object extensively:

```javascript
// src/ui/registry/Methods/enrollStudent.js
function enrollStudent(payload, actionparams, event, form) {
  if (form) {
    form.setSubmitting(true);
  }
  
  Window.executeService(
    'EnrollmentService',
    payload,
    null,
    function(response) {
      if (form) {
        form.setSubmitting(false);
        form.resetForm();
      }
      
      Window.showMessage({ Default: 'Student enrolled!' });
      
      // Dispatch to Redux
      Application.store.dispatch({
        type: 'STUDENT_ENROLLED',
        payload: response.data
      });
      
      // Navigate after delay
      setTimeout(() => {
        Window.redirectPage('/students');
      }, 1500);
    },
    function(error) {
      if (form) {
        form.setSubmitting(false);
      }
      Window.showError({ Default: 'Enrollment failed' });
    }
  );
}
```

## Integration with Sagas

Sagas can also use Window methods:

```javascript
import { call } from 'redux-saga/effects';

function* enrollStudentSaga(action) {
  try {
    const response = yield call(
      DataSource.ExecuteService,
      'EnrollmentService',
      RequestBuilder.DefaultRequest(null, action.payload)
    );
    
    Window.showMessage({ Default: 'Enrollment successful!' });
    Window.redirectPage('/students');
    
  } catch (error) {
    Window.showError({ Default: 'Enrollment failed' });
  }
}
```

## Best Practices

1. **Error Handling**: Always provide error callbacks
2. **User Feedback**: Use showMessage/showError for UX
3. **Async Awareness**: Remember service calls are async
4. **Navigation Timing**: Add delays before navigation to show messages
5. **Action Parameters**: Use actionparams for configuration
6. **Consistent Messaging**: Use clear, user-friendly messages
7. **Try-Catch**: Wrap Window calls in try-catch for safety

## See Also

- [Methods](methods.md) - Custom JavaScript functions
- [Sagas](sagas.md) - Async operations
- [Actions](actions.md) - Action definitions
- [Services](services.md) - Backend service calls
