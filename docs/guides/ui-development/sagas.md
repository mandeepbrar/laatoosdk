# Sagas - Async Operations

Sagas handle asynchronous operations and side effects in Laatoo UI plugins using redux-saga. They process actions, call services, and dispatch state updates.

## Purpose

Sagas enable:
- Async service calls
- Complex action workflows
- Side effect management
- Action chaining
- Error handling

## Basic Saga Structure

```javascript
import { takeEvery, call, put, fork } from 'redux-saga/effects';
import Actions from './actions';

function* mySagaWorker(action) {
  try {
    // Perform async work
    const result = yield call(someAsyncFunction, action.payload);
    
    // Dispatch success action
    yield put({ type: Actions.SUCCESS, payload: result });
  } catch (error) {
    // Handle errors
    yield put({ type: Actions.FAILURE, payload: error });
  }
}

function* mySagaWatcher() {
  yield takeEvery(Actions.MY_ACTION, mySagaWorker);
}

function* rootSaga() {
  yield fork(mySagaWatcher);
}

// Register the saga
Application.Register('Sagas', 'mySaga', rootSaga);
```

## Registering Sagas

Use `Application.Register()` in your Initialize method or saga file:

```javascript
// In Initialize method
function Initialize(appName, ins, mod, settings, def, req) {
  module = this;
  
  // Define and register saga
  function* loadDataSaga() {
    yield takeEvery('LOAD_STUDENTS', function*(action) {
      try {
        const response = yield call(loadStudents, action.payload);
        yield put({ type: 'STUDENTS_LOADED', payload: response.data });
      } catch (error) {
        yield put({ type: 'LOAD_ERROR', payload: error });
      }
    });
  }
  
  Application.Register('Sagas', 'studentDataSaga', loadDataSaga);
}
```

## Common Saga Patterns

### 1. Service Call Saga

```javascript
import { DataSource, RequestBuilder } from 'jsui';

function* fetchStudentsSaga(action) {
  try {
    const req = RequestBuilder.DefaultRequest(null, action.payload);
    const response = yield call(DataSource.ExecuteService, 'Student.Query', req);
    yield put({ type: 'STUDENTS_LOADED', payload: response.data });
  } catch (error) {
    Window.showError({ Default: 'Failed to load students' });
  }
}
```

### 2. Action Execution Saga

```javascript
function* executeActionSaga(action) {
  const laatooAction = _reg('Actions', action.meta.actionname);
  
  if (laatooAction) {
    const payload = action.payload;
    const actionparams = action.meta.actionparams;
    
    switch(laatooaction.actiontype) {
      case 'executeservice':
        const req = RequestBuilder.URLParamsRequest(actionparams.serviceparams, payload);
        const response = yield call(DataSource.ExecuteService, laatooaction.servicename, req);
        yield put({ type: 'SERVICE_SUCCESS', payload: response });
        break;
        
      case 'method':
        const method = _reg('Methods', laatooAction.method);
        if (method) {
          method(payload, actionparams);
        }
        break;
    }
  }
}
```

### 3. Chained Actions Saga

```javascript
function* processExamSaga(action) {
  try {
    // Step 1: Submit exam
    const submitResponse = yield call(DataSource.ExecuteService, 'Exam.Submit', {
      data: action.payload
    });
    
    // Step 2: Trigger grading
    const gradeResponse = yield call(DataSource.ExecuteService, 'Exam.Grade', {
      examAttemptId: submitResponse.data.Id
    });
    
    // Step 3: Update UI
    yield put({ type: 'EXAM_PROCESSED', payload: gradeResponse.data });
    
    // Step 4: Show success
    Window.showMessage({ Default: 'Exam graded successfully!' });
    
  } catch (error) {
    yield put({ type: 'EXAM_ERROR', payload: error });
    Window.showError(error);
  }
}
```

## Saga Effects

### takeEvery

Listens for every occurrence of an action:

```javascript
yield takeEvery('LOAD_DATA', loadDataWorker);
```

### takeLatest

Cancels previous instances and only runs the latest:

```javascript
import { takeLatest } from 'redux-saga/effects';

yield takeLatest('SEARCH_STUDENTS', searchWorker);
```

### call

Calls a function and waits for the result:

```javascript
const result = yield call(DataSource.ExecuteService, 'MyService', request);
```

### put

Dispatches an action:

```javascript
yield put({ type: 'DATA_LOADED', payload: data });
```

### fork

Spawns a non-blocking saga:

```javascript
yield fork(backgroundSaga);
```

## File Location

Sagas can be defined in:
- `src/ui/js/index.js` (in Initialize method)
- `src/ui/js/sagas/*.js` (separate saga files)

## Example: Student Management Saga

```javascript
import { takeEvery, call, put } from 'redux-saga/effects';
import { DataSource, RequestBuilder } from 'jsui';

function* enrollStudentSaga(action) {
  try {
    // Show loading
    yield put({ type: 'ENROLLING_STUDENT' });
    
    // Call enrollment service
    const req = RequestBuilder.DefaultRequest(null, action.payload);
    const response = yield call(DataSource.ExecuteService, 'EnrollmentService', req);
    
    // Update state
    yield put({ type: 'STUDENT_ENROLLED', payload: response.data });
    
    // Show success
    Window.showMessage({ Default: 'Student enrolled successfully!' });
    
    // Navigate
    Window.redirectPage('/students');
    
  } catch (error) {
    yield put({ type: 'ENROLLMENT_ERROR', payload: error });
    Window.showError({ Default: 'Enrollment failed' });
  }
}

function* studentSagaWatcher() {
  yield takeEvery('ENROLL_STUDENT', enrollStudentSaga);
}

function* studentRootSaga() {
  yield fork(studentSagaWatcher);
}

// Register
Application.Register('Sagas', 'studentEnrollmentSaga', studentRootSaga);
```

## Best Practices

1. **Error Handling**: Always wrap saga logic in try-catch
2. **Loading States**: Dispatch loading actions before async calls
3. **User Feedback**: Use Window.showMessage/showError for UX
4. **Action Names**: Use descriptive, uppercase action constants
5. **Naming**: Give sagas descriptive names when registering
6. **Single Responsibility**: Each saga should handle one workflow
7. **Testing**: Sagas are easy to test with redux-saga-test-plan

## See Also

- [Reducers](reducers.md) - State management
- [Actions](actions.md) - Action definitions
- [Initialize Method](initialize.md) - Plugin initialization
- [Window Object](window-object.md) - Global utilities
