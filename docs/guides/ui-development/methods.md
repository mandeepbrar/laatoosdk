# Methods Registry

Methods are custom JavaScript functions registered in the Methods registry to provide reusable UI logic and event handlers.

## Purpose

Methods enable:
- Custom form validation
- Complex UI interactions
- Reusable business logic
- Event handlers for actions
- Data transformations

## Creating Methods

Methods are JavaScript files in `src/ui/registry/Methods/`:

```
src/ui/registry/
└── Methods/
    ├── validateForm.js
    ├── calculateTotal.js
    └── submitExam.js
```

## Basic Method Structure

```javascript
// src/ui/registry/Methods/validateForm.js
function validateForm(payload, actionparams, event, form) {
  // payload: data passed to the method
  // actionparams: parameters from action definition
  // event: DOM event (if triggered by UI event)
  // form: Formik form object (if called from form)
  
  if (!payload.Email || !payload.Email.includes('@')) {
    Window.showError({ Default: 'Invalid email address' });
    return false;
  }
  
  return true;
}
```

## Method Signature

Methods receive four parameters:

```javascript
function methodName(payload, actionparams, event, form) {
  // payload: data from action/form
  // actionparams: action parameters
  // event: browser event object
  // form: Formik form instance (has setFieldValue, setFieldError, etc.)
}
```

## Accessing Methods

Methods are automatically registered and can be called:

### From Actions

```yaml
# src/ui/registry/Actions/validateStudent.yml
actiontype: method
method: validateForm
```

### From Sagas

```javascript
const method = _reg('Methods', 'validateForm');
if (method) {
  method(payload, actionparams);
}
```

### From Other JavaScript

```javascript
const calculateTotal = _reg('Methods', 'calculateTotal');
const total = calculateTotal({ items: cartItems });
```

## Common Method Patterns

### 1. Form Validation

```javascript
// src/ui/registry/Methods/validateStudentForm.js
function validateStudentForm(payload, actionparams, event, form) {
  const errors = {};
  
  if (!payload.Name || payload.Name.trim() === '') {
    errors.Name = 'Name is required';
  }
  
  if (!payload.Email) {
    errors.Email = 'Email is required';
  } else if (!/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i.test(payload.Email)) {
    errors.Email = 'Invalid email address';
  }
  
  if (!payload.DateOfBirth) {
    errors.DateOfBirth = 'Date of birth is required';
  }
  
  if (Object.keys(errors).length > 0) {
    // Set errors on form
    if (form) {
      Object.keys(errors).forEach(field => {
        form.setFieldError(field, errors[field]);
      });
    }
    Window.showError({ Default: 'Please fix form errors' });
    return false;
  }
  
  return true;
}
```

### 2. Service Call

```javascript
// src/ui/registry/Methods/loadStudentData.js
function loadStudentData(payload, actionparams, event, form) {
  const studentId = payload.studentId || actionparams.studentId;
  
  Window.executeService(
    'Student.GetById',
    { id: studentId },
    null,
    function(response) {
      // Success callback
      if (form) {
        // Populate form with data
        form.setValues(response.data);
      }
      Window.showMessage({ Default: 'Student data loaded' });
    },
    function(error) {
      // Error callback
      Window.showError({ Default: 'Failed to load student' });
    }
  );
}
```

### 3. Data Transformation

```javascript
// src/ui/registry/Methods/formatExamData.js
function formatExamData(payload, actionparams, event, form) {
  const formatted = {
    ...payload,
    ExamId: (payload.ExamId || '').toUpperCase(),
    StartDate: new Date(payload.StartDate).toISOString(),
    EndDate: new Date(payload.EndDate).toISOString(),
    MaxScore: parseInt(payload.MaxScore, 10),
    PassingScore: parseInt(payload.PassingScore, 10)
  };
  
  return formatted;
}
```

### 4. Dynamic Field Updates

```javascript
// src/ui/registry/Methods/updateGradeOptions.js
function updateGradeOptions(payload, actionparams, event, form) {
  const examType = payload.ExamType;
  
  if (form) {
    if (examType === 'Midterm') {
      form.setFieldValue('MaxScore', 50);
      form.setFieldValue('PassingScore', 30);
    } else if (examType === 'Final') {
      form.setFieldValue('MaxScore', 100);
      form.setFieldValue('PassingScore', 60);
    }
  }
}
```

### 5. Dispatch Redux Action

```javascript
// src/ui/registry/Methods/updateStudentFilter.js
function updateStudentFilter(payload, actionparams, event, form) {
  const { createAction } = require('jsui');
  
  Application.store.dispatch(createAction('SET_STUDENT_FILTER', {
    grade: payload.grade,
    status: payload.status
  }));
  
  // Trigger data reload
  Window.executeAction('loadStudents');
}
```

## Using Window Object in Methods

Methods have access to the global `Window` object:

```javascript
function submitExam(payload, actionparams, event, form) {
  // Execute service
  Window.executeService('Exam.Submit', payload, null,
    function(response) {
      Window.showMessage({ Default: 'Exam submitted successfully!' });
      Window.redirectPage('/exams');
    },
    function(error) {
      Window.showError({ Default: 'Failed to submit exam' });
    }
  );
}
```

## Example: Student Management Methods

```javascript
// src/ui/registry/Methods/enrollStudent.js
function enrollStudent(payload, actionparams, event, form) {
  // Validate
  if (!payload.Name || !payload.Email) {
    Window.showError({ Default: 'Name and email are required' });
    return;
  }
  
  // Show loading (optional)
  if (form) {
    form.setSubmitting(true);
  }
  
  // Call enrollment service
  Window.executeService(
    'EnrollmentService',
    {
      Name: payload.Name,
      Email: payload.Email,
      Grade: payload.Grade,
      DateOfBirth: payload.DateOfBirth
    },
    null,
    function(response) {
      // Success
      if (form) {
        form.setSubmitting(false);
        form.resetForm();
      }
      
      Window.showMessage({ Default: 'Student enrolled successfully!' });
      
      // Dispatch to Redux
      Application.store.dispatch({
        type: 'STUDENT_ENROLLED',
        payload: response.data
      });
      
      // Navigate
      setTimeout(() => {
        Window.redirectPage('/students');
      }, 1500);
    },
    function(error) {
      // Error
      if (form) {
        form.setSubmitting(false);
      }
      Window.showError({ Default: 'Failed to enroll student' });
    }
  );
}
```

## Best Practices

1. **Single Responsibility**: Each method should do one thing well
2. **Error Handling**: Always handle errors gracefully
3. **User Feedback**: Use Window.showMessage/showError for UX
4. **Return Values**: Return meaningful values (true/false for validation)
5. **Form Integration**: Use the form parameter to interact with Formik
6. **Naming**: Use descriptive names (verb + noun pattern)
7. **Documentation**: Add comments for complex logic

## File Naming

- Use camelCase for filenames: `validateForm.js`
- Use descriptive names: `calculateGrade.js` not `calc.js`
- Match function name to filename

## See Also

- [Actions](actions.md) - Calling methods from actions
- [Forms](forms.md) - Form integration
- [Window Object](window-object.md) - Global utilities
- [Sagas](sagas.md) - Calling methods from sagas
