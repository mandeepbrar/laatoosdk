
// enrollStudent.js - Submit student enrollment
function enrollStudent(payload, actionparams, event, form) {
    // Show loading state
    if (form) {
        form.setSubmitting(true);
    }

    // Call createStudent service (defined in Services.yml)
    Window.executeService(
        'createStudent',
        {
            Name: payload.Name,
            Email: payload.Email,
            Grade: payload.Grade,
            DateOfBirth: payload.DateOfBirth,
            Phone: payload.Phone || ''
        },
        null,
        function (response) {
            // Success callback
            if (form) {
                form.setSubmitting(false);
                form.resetForm();
            }

            Window.showMessage({ Default: 'Student enrolled successfully!' });

            // Dispatch to Redux store
            Application.store.dispatch({
                type: 'STUDENT_ENROLLED',
                payload: response.data
            });

            // Navigate to students list after delay
            setTimeout(() => {
                Window.redirectPage('/students');
            }, 1500);
        },
        function (error) {
            // Error callback
            if (form) {
                form.setSubmitting(false);
            }
            Window.showError({ Default: 'Failed to enroll student. Please try again.' });
            console.error('Enrollment error:', error);
        }
    );
}

