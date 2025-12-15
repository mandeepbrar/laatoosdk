// StudentData reducer - Manages students state
import Actions from '../actions';

const initialState = {
    students: [],
    selectedStudent: null,
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
        case 'LOADING_STUDENTS':
            return {
                ...state,
                loading: true,
                error: null
            };

        case 'STUDENTS_LOADED':
            return {
                ...state,
                students: action.payload || [],
                loading: false
            };

        case 'LOAD_ERROR':
            return {
                ...state,
                loading: false,
                error: action.payload
            };

        case 'STUDENT_ENROLLED':
            return {
                ...state,
                students: [...state.students, action.payload]
            };

        case 'SELECT_STUDENT':
            return {
                ...state,
                selectedStudent: action.payload
            };

        case 'SET_STUDENT_FILTER':
            return {
                ...state,
                filters: {
                    ...state.filters,
                    ...action.payload
                }
            };

        case Actions.LOGOUT:
            return initialState;

        default:
            return state;
    }
};

Application.Register('Reducers', 'StudentData', StudentData);
