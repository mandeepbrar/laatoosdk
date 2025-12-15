# Reducers - State Management

Reducers manage application state in Laatoo UI plugins using Redux. They handle state updates in response to actions.

## Purpose

Reducers enable:
- Centralized state management
- Predictable state updates
- State persistence across components
- Time-travel debugging (Redux DevTools)

## Basic Reducer Structure

```javascript
import Actions from './actions';

const initialState = {
  students: [],
  loading: false,
  error: null
};

const StudentReducer = (state, action) => {
  if (!state) {
    return initialState;
  }
  
  switch (action.type) {
    case 'LOAD_STUDENTS':
      return {
        ...state,
        loading: true,
        error: null
      };
      
    case 'STUDENTS_LOADED':
      return {
        ...state,
        students: action.payload,
        loading: false
      };
      
    case 'LOAD_ERROR':
      return {
        ...state,
        loading: false,
        error: action.payload
      };
      
    default:
      return state;
  }
};

// Register the reducer
Application.Register('Reducers', 'Students', StudentReducer);
```

## Registering Reducers

Use `Application.Register()` to register reducers:

```javascript
// In separate file: src/ui/js/reducers/Students.js
const Students = (state, action) => {
  // reducer logic
};

Application.Register('Reducers', 'Students', Students);
export default Students;
```

Or in Initialize method:

```javascript
function Initialize(appName, ins, mod, settings, def, req) {
  module = this;
  
  const MyReducer = (state = {}, action) => {
    switch (action.type) {
      case 'UPDATE_DATA':
        return { ...state, data: action.payload };
      default:
        return state;
    }
  };
  
  Application.Register('Reducers', 'MyData', MyReducer);
}
```

## Accessing State

Access Redux state in components via `Application.store`:

```javascript
// Get current state
const currentState = Application.store.getState();
const students = currentState.Students.students;

// Subscribe to state changes
Application.store.subscribe(() => {
  const state = Application.store.getState();
  console.log('State updated:', state);
});
```

## Dispatching Actions

Update state by dispatching actions:

```javascript
import { createAction } from 'jsui';

// Dispatch an action
Application.store.dispatch(createAction('LOAD_STUDENTS', { filter: 'active' }));

// Or using Window.executeAction
Window.executeAction('loadStudents', { filter: 'active' });
```

## Common Reducer Patterns

### 1. List Management

```javascript
const ExamList = (state = { exams: [], selected: null }, action) => {
  switch (action.type) {
    case 'ADD_EXAM':
      return {
        ...state,
        exams: [...state.exams, action.payload]
      };
      
    case 'UPDATE_EXAM':
      return {
        ...state,
        exams: state.exams.map(exam =>
          exam.Id === action.payload.Id ? action.payload : exam
        )
      };
      
    case 'REMOVE_EXAM':
      return {
        ...state,
        exams: state.exams.filter(exam => exam.Id !== action.payload)
      };
      
    case 'SELECT_EXAM':
      return {
        ...state,
        selected: action.payload
      };
      
    default:
      return state;
  }
};

Application.Register('Reducers', 'ExamList', ExamList);
```

### 2. UI State Management

```javascript
const UIState = (state = { showMenu: false, activeTab: 'students' }, action) => {
  switch (action.type) {
    case 'TOGGLE_MENU':
      return {
        ...state,
        showMenu: !state.showMenu
      };
      
    case 'SET_ACTIVE_TAB':
      return {
        ...state,
        activeTab: action.payload
      };
      
    default:
      return state;
  }
};

Application.Register('Reducers', 'UIState', UIState);
```

### 3. Messages/Notifications Reducer

```javascript
const Messages = (state = {}, action) => {
  switch (action.type) {
    case 'DISPLAY_ERROR':
      return {
        Message: action.payload.message,
        Type: 'Error',
        Time: new Date().getTime()
      };
      
    case 'SHOW_MESSAGE':
      return {
        Message: action.payload.message,
        Type: 'Message',
        Time: new Date().getTime()
      };
      
    case 'CLEAR_MESSAGE':
      return {};
      
    default:
      return state || {};
  }
};

Application.Register('Reducers', 'Messages', Messages);
```

## File Location

Reducers are typically defined in:
- `src/ui/js/reducers/*.js` (separate reducer files)
- `src/ui/js/index.js` (in Initialize method)

## Example: Student Management Reducers

```javascript
// src/ui/js/reducers/StudentData.js
const initialState = {
  students: [],
  enrollments: [],
  loading: false,
  error: null,
  filters: {
    grade: null,
    status: 'active'
  }
};

const StudentData = (state, action) => {
  if (!state) {
    return initialState;
  }
  
  switch (action.type) {
    case 'ENROLLING_STUDENT':
      return {
        ...state,
        loading: true,
        error: null
      };
      
    case 'STUDENT_ENROLLED':
      return {
        ...state,
        students: [...state.students, action.payload],
        loading: false
      };
      
    case 'ENROLLMENT_ERROR':
      return {
        ...state,
        loading: false,
        error: action.payload
      };
      
    case 'SET_STUDENT_FILTER':
      return {
        ...state,
        filters: {
          ...state.filters,
          ...action.payload
        }
      };
      
    case 'STUDENTS_LOADED':
      return {
        ...state,
        students: action.payload,
        loading: false
      };
      
    case 'LOGOUT':
      return initialState;
      
    default:
      return state;
  }
};

Application.Register('Reducers', 'StudentData', StudentData);
```

## Best Practices

1. **Immutability**: Always return new state objects, never mutate
2. **Default State**: Always handle undefined state by returning initialState
3. **Switch Statements**: Use switch for clarity with multiple actions
4. **Logout Handling**: Reset to initialState on LOGOUT action
5. **Naming**: Use descriptive reducer names when registering
6. **Separation**: Keep reducers focused on a single domain
7. **Type Constants**: Define action types as constants to avoid typos

## Combining with Sagas

Reducers work with sagas for complete state management:

```javascript
// Saga dispatches actions
function* loadStudentsSaga() {
  try {
    // Dispatch loading action (handled by reducer)
    yield put({ type: 'LOADING_STUDENTS' });
    
    const response = yield call(DataSource.ExecuteService, 'Student.Query');
    
    // Dispatch success action (handled by reducer)
    yield put({ type: 'STUDENTS_LOADED', payload: response.data });
  } catch (error) {
    // Dispatch error action (handled by reducer)
    yield put({ type: 'LOAD_ERROR', payload: error });
  }
}

// Reducer handles the actions
const Students = (state = { items: [], loading: false }, action) => {
  switch (action.type) {
    case 'LOADING_STUDENTS':
      return { ...state, loading: true };
    case 'STUDENTS_LOADED':
      return { items: action.payload, loading: false };
    case 'LOAD_ERROR':
      return { ...state, loading: false, error: action.payload };
    default:
      return state;
  }
};
```

## See Also

- [Sagas](sagas.md) - Async operations
- [Actions](actions.md) - Action definitions
- [Initialize Method](initialize.md) - Plugin initialization
