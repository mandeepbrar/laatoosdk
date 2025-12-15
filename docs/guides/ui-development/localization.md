# Localization and Internationalization

Add multi-language support to your Laatoo applications with the properties system.

## Quick Start

1. Create property files in `properties/` directory
2. Define translations in YAML format
3. Reference in UI and server code

## Property Files

### Default Properties

**File**: `dev/plugins/studentmgmt/properties/default.yml`

```yaml
# Student Management Properties
studentname: Student Name
email: Email Address
grade: Grade Level
dateofbirth: Date of Birth
enrollstudent: Enroll Student
submit: Submit
cancel: Cancel

# Messages
enrollmentsuccess: Student enrolled successfully!
enrollmentfailed: Failed to enroll student
validationfailed: Please fix the form errors

# Button labels
save: Save
delete: Delete
edit: Edit

# Table headers
students: Students
exams: Exams
results: Results
```

### Language-Specific Properties

**File**: `properties/es.yml` (Spanish)

```yaml
studentname: Nombre del Estudiante
email: Correo Electrónico
grade: Grado
dateofbirth: Fecha de Nacimiento
enrollstudent: Inscribir Estudiante
submit: Enviar
cancel: Cancelar

enrollmentsuccess: ¡Estudiante inscrito exitosamente!
enrollmentfailed: Error al inscribir estudiante
```

**File**: `properties/fr.yml` (French)

```yaml
studentname: Nom de l'Étudiant
email: Adresse E-mail
grade: Niveau
dateofbirth: Date de Naissance
enrollstudent: Inscrire un Étudiant
```

## Using in UI

### In Forms

```xml
<!-- Form field labels use properties -->
<Form form="student_enrollment" module="formikforms">
  <Field 
    type="string" 
    name="Name" 
    label="{{jsreplace 'props.studentmgmt.studentname'}}"
    module="formikforms"
  />
  
  <Field 
    type="email" 
    name="Email" 
    label="{{jsreplace 'props.studentmgmt.email'}}"
    module="formikforms"
  />
  
  <Button type="submit">
    {{jsreplace "props.studentmgmt.submit"}}
  </Button>
</Form>
```

### In JavaScript Methods

```javascript
// Show localized messages
Window.showMessage({ 
  Default: module.properties.enrollmentsuccess 
});

Window.showError({ 
  Default: module.properties.enrollmentfailed 
});

// Access properties
const submitLabel = module.properties.submit;
const cancelLabel = module.properties.cancel;
```

### In Pages/Blocks

```xml
<Block className="student-header">
  <h2 module="html">{{jsreplace "props.studentmgmt.students"}}</h2>
</Block>
```

## Using in Server

Access properties from server-side code:

```go
func (s *StudentService) Invoke(ctx interfaces.ServerContext) error {
    // Get property from module
    errorMsg := ctx.GetProperty("enrollmentfailed")
    
    if err != nil {
        return ctx.SetError(errorMsg)
    }
    
    return nil
}
```

## Module Properties Access

Properties are organized by module/plugin name:

```javascript
// In UI Initialize method
function Initialize(appName, ins, mod, settings, def, req) {
    module = this;
    
    // Access module properties
    module.properties = Application.Properties[ins];
    
    // Use properties
    const studentLabel = module.properties.studentname;
    const emailLabel = module.properties.email;
}
```

## Format Patterns

### Simple Text

```yaml
welcome: Welcome to Student Management
```

```xml
<span module="html">{{jsreplace "props.studentmgmt.welcome"}}</span>
```

### Placeholders (Common Pattern)

```yaml
studentcount: "Total Students: {count}"
examschedule: "Exam on {date} at {time}"
```

```javascript
// Replace placeholders
let message = module.properties.studentcount;
message = message.replace('{count}', studentList.length);
```

### Multi-line Text

```yaml
welcomemessage: |
  Welcome to the Student Management System.
  Please select an option from the menu.
```

## Property Naming Conventions

1. **Use lowercase**: `studentname` not `StudentName`
2. **Use descriptive names**: `enrollmentsuccess` not `msg1`
3. **Group related properties**: All student-related properties together
4. **Consistent suffixes**:
   - `*success` for success messages
   - `*failed` or `*error` for errors
   - `*confirm` for confirmations

## Language Selection

Set user language via:

**Environment Variable**:
```bash
export LAATOO_defaultlanguage=es
```

**User Preferences**:
```javascript
// Store user language preference
Application.setUserLanguage('es');
```

## Complete Example: Student Form

**Properties** (`properties/default.yml`):
```yaml
# Form labels
form_student_title: Student Enrollment
form_name: Full Name
form_name_placeholder: Enter student's full name
form_email: Email Address
form_email_placeholder: student@example.com
form_grade: Grade Level
form_dob: Date of Birth

# Buttons
btn_enroll: Enroll Student
btn_cancel: Cancel

# Messages
msg_enroll_success: Student enrolled successfully!
msg_enroll_error: Failed to enroll student
msg_validation_error: Please correct the form errors

# Validation
val_name_required: Name is required
val_email_required: Email is required
val_email_invalid: Invalid email address
```

**Form UI** (`src/ui/registry/Forms/studentform.xml`):
```xml
<Form form="student_enrollment" module="formikforms">
  <h2 module="html">{{jsreplace "props.studentmgmt.form_student_title"}}</h2>
  
  <Field 
    type="string" 
    name="Name"
    label="{{jsreplace 'props.studentmgmt.form_name'}}"
    placeholder="{{jsreplace 'props.studentmgmt.form_name_placeholder'}}"
    required="true"
    module="formikforms"
  />
  
  <Field 
    type="email" 
    name="Email"
    label="{{jsreplace 'props.studentmgmt.form_email'}}"
    placeholder="{{jsreplace 'props.studentmgmt.form_email_placeholder'}}"
    required="true"
    module="formikforms"
  />
  
  <Field 
    type="number" 
    name="Grade"
    label="{{jsreplace 'props.studentmgmt.form_grade'}}"
    module="formikforms"
  />
  
  <Button type="submit">
    {{jsreplace "props.studentmgmt.btn_enroll"}}
  </Button>
  
  <Button type="button" action="cancelEnrollment">
    {{jsreplace "props.studentmgmt.btn_cancel"}}
  </Button>
</Form>
```

**Method** (`src/ui/registry/Methods/enrollStudent.js`):
```javascript
function enrollStudent(payload, actionparams, event, form) {
    if (form) {
        form.setSubmitting(true);
    }
    
    Window.executeService('createStudent', payload, null,
        function(response) {
            if (form) {
                form.setSubmitting(false);
                form.resetForm();
            }
            
            // Use localized success message
            Window.showMessage({ 
                Default: module.properties.msg_enroll_success 
            });
            
            Window.redirectPage('/students');
        },
        function(error) {
            if (form) {
                form.setSubmitting(false);
            }
            
            // Use localized error message
            Window.showError({ 
                Default: module.properties.msg_enroll_error 
            });
        }
    );
}
```

## Best Practices

1. **Default Properties Required**: Always have `default.yml`
2. **Consistent Keys**: Use same keys across all language files
3. **Test All Languages**: Verify translations display correctly
4. **Professional Translations**: Use native speakers for translations
5. **Context Matters**: "Close" (verb) vs "Close" (adjective) need different keys
6. **Avoid Hardcoding**: Never hardcode display text
7. **Document Property Usage**: Comment property purpose in YAML
8. **Version Control**: Track property file changes

## Common Patterns

### Error Messages

```yaml
# Validation errors
error_required: This field is required
error_invalid_email: Please enter a valid email
error_min_length: Minimum {min} characters required
error_max_length: Maximum {max} characters allowed

# System errors  
error_network: Network error occurred
error_server: Server error, please try again
error_unauthorized: You don't have permission
```

### Confirmation Dialogs

```yaml
confirm_delete: Are you sure you want to delete this student?
confirm_grade: Submit final grade for this exam?
confirm_cancel: Discard unsaved changes?
```

## Troubleshooting

### Properties Not Loading

```javascript
// Check if properties loaded
console.log('Properties:', module.properties);

// Verify module name
console.log('Module:', ins);
```

### Wrong Language Displayed

```javascript
// Check current language
console.log('Language:', Application.getCurrentLanguage());

// Force language
Application.setLanguage('es');
```

## See Also

- [Forms Guide](forms.md)
- [UI Development](creating-ui-plugins.md)
- [Blocks and Components](blocks.md)
