// validateEnrollment.js - Validate student enrollment form
function validateEnrollment(payload, actionparams, event, form) {
    const errors = {};

    // Name validation
    if (!payload.Name || payload.Name.trim() === '') {
        errors.Name = 'Student name is required';
    }

    // Email validation
    if (!payload.Email) {
        errors.Email = 'Email is required';
    } else if (!/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i.test(payload.Email)) {
        errors.Email = 'Invalid email format';
    }

    // Date of Birth validation
    if (!payload.DateOfBirth) {
        errors.DateOfBirth = 'Date of birth is required';
    } else {
        const dob = new Date(payload.DateOfBirth);
        const age = (new Date() - dob) / (1000 * 60 * 60 * 24 * 365);
        if (age < 5 || age > 100) {
            errors.DateOfBirth = 'Please enter a valid date of birth';
        }
    }

    // Grade validation
    if (!payload.Grade) {
        errors.Grade = 'Grade is required';
    }

    if (Object.keys(errors).length > 0) {
        // Show errors on form fields
        if (form) {
            Object.keys(errors).forEach(field => {
                form.setFieldError(field, errors[field]);
            });
        }

        Window.showError({ Default: 'Please fix the form errors' });
        return false;
    }

    Window.showMessage({ Default: 'Validation passed' });
    return true;
}
