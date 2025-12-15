// calculateGrade.js - Calculate exam grade percentage
function calculateGrade(payload, actionparams, event, form) {
    const score = parseFloat(payload.Score) || 0;
    const maxScore = parseFloat(payload.MaxScore) || 100;

    if (maxScore === 0) {
        Window.showError({ Default: 'Max score cannot be zero' });
        return;
    }

    const percentage = (score / maxScore) * 100;
    const rounded = Math.round(percentage * 100) / 100;

    // Update form field if available
    if (form) {
        form.setFieldValue('Percentage', rounded);
        form.setFieldValue('Passed', rounded >= (payload.PassingScore || 60));
    }

    return rounded;
}
