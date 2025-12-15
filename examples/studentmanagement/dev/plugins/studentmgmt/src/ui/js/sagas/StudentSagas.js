// Student enrollment saga
import { takeEvery, call, put, fork } from 'redux-saga/effects';
import { DataSource, RequestBuilder } from 'jsui';

function* enrollStudentSaga(action) {
    try {
        // Dispatch loading action
        yield put({ type: 'ENROLLING_STUDENT' });

        // Call createStudent service (defined in Services.yml)
        const req = RequestBuilder.DefaultRequest(null, action.payload);
        const response = yield call(
            DataSource.ExecuteService,
            'createStudent',
            req
        );

        // Update state with new student
        yield put({ type: 'STUDENT_ENROLLED', payload: response.data });

        // Show success message
        Window.showMessage({ Default: 'Student enrolled successfully!' });

        // Navigate to students list
        setTimeout(() => {
            Window.redirectPage('/students');
        }, 1500);

    } catch (error) {
        yield put({ type: 'ENROLLMENT_ERROR', payload: error });
        Window.showError({ Default: 'Failed to enroll student' });
        console.error('Enrollment error:', error);
    }
}

function* loadStudentsSaga(action) {
    try {
        yield put({ type: 'LOADING_STUDENTS' });

        // Call queryStudents service (defined in Services.yml)
        const req = RequestBuilder.DefaultRequest(null, action.payload || {});
        const response = yield call(
            DataSource.ExecuteService,
            'queryStudents',
            req
        );

        yield put({ type: 'STUDENTS_LOADED', payload: response.data });

    } catch (error) {
        yield put({ type: 'LOAD_ERROR', payload: error });
        Window.showError({ Default: 'Failed to load students' });
    }
}

function* studentSagaWatcher() {
    yield takeEvery('ENROLL_STUDENT', enrollStudentSaga);
    yield takeEvery('LOAD_STUDENTS', loadStudentsSaga);
}

function* studentRootSaga() {
    yield fork(studentSagaWatcher);
}

Application.Register('Sagas', 'studentMgmtSaga', studentRootSaga);
