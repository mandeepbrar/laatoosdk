# Chapter 6: Forms

Create forms for student enrollment and exam creation.

## Overview

Forms in Laatoo are declaratively defined in XML and support:
- Various field types (text, email, date, number, select, etc.)
- Built-in validation
- Actions (presubmit, submit, success)
- Dynamic field visibility
- File uploads

We'll create:
- **Student Enrollment Form**: Register new students
- **Exam Creation Form**: Create and configure exams

> **Note**: Forms in Laatoo are not standalone. They must be included in a Page to be accessible via a URL. We will create the pages for these forms in Step 7.

## Step 1: Create Student Enrollment Form

Create `src/ui/registry/Forms/student-enrollment.xml`:

```xml
<Form form="student_enrollment" module="formikforms" action="enrollStudent">
  <h2 module="html">Student Enrollment</h2>
  <p module="html">Enter student information to register</p>
  
  <!-- Name Field -->
  <Field 
    type="string" 
    name="Name" 
    label="Full Name"
    className="w100" 
    module="formikforms"
    required="true"
    placeholder="Enter student's full name"/>
    
  <!-- Email Field -->
  <Field 
    type="email" 
    name="Email" 
    label="Email Address"
    className="w100" 
    module="formikforms"
    required="true"
    placeholder="student@example.com"
    validation="email"/>
    
  <!-- Grade Field -->
  <Field 
    type="select" 
    name="Grade" 
    label="Grade Level"
    className="w50" 
    module="formikforms"
    required="true">
    <option value="">Select Grade</option>
    <option value="9">Grade 9</option>
    <option value="10">Grade 10</option>
    <option value="11">Grade 11</option>
    <option value="12">Grade 12</option>
  </Field>
  
  <!-- Date of Birth -->
  <Field 
    type="date" 
    name="DateOfBirth" 
    label="Date of Birth"
    className="w50" 
    module="formikforms"/>
    
  <!-- Phone Number -->
  <Field 
    type="string" 
    name="Phone" 
    label="Phone Number"
    className="w50" 
    module="formikforms"
    placeholder="(555) 123-4567"
    pattern="^[\d\(\)\-\s]+$"/>
    
  <!-- Submit Button -->
  <Button type="submit" className="btn-primary">
    Enroll Student
  </Button>
</Form>
```

## Step 2: Create Form Actions

Forms need actions to handle submission. We'll create a `save` action that calls the backend service.

Create `src/ui/registry/Actions/save.yml`:

```yaml
actiontype: executeservice
servicename: studentmgmt.Student.Create
successaction: navigateToDashboard
params:
  data: {{jsreplace "ctx.formData"}}
```

Create `src/ui/registry/Actions/navigateToDashboard.yml`:

```yaml
actiontype: navigate
page: /dashboard
```

## Step 3: Create Form Pages

Create `src/ui/registry/Forms/exam-creation.xml`:

```xml
<Form form="exam_creation" module="formikforms" action="createExam">
  <h2 module="html">Create Exam</h2>
  
  <!-- Basic Information -->
  <Block className="form-section">
    <h3 module="html">Basic Information</h3>
    
    <Field 
      type="string" 
      name="ExamId" 
      label="Exam ID"
      className="w50" 
      module="formikforms"
      required="true"
      pattern="^[A-Z0-9_-]+$"
      placeholder="MATH_FINAL_2024"/>
      
    <Field 
      type="string" 
      name="Title" 
      label="Exam Title"
      className="w100" 
      module="formikforms"
      required="true"
      placeholder="Mathematics Final Exam"/>
      
    <Field 
      type="textarea" 
      name="Description" 
      label="Description"
      className="w100" 
      module="formikforms"
      rows="4"
      placeholder="Enter exam description and instructions"/>
  </Block>
  
  <!-- Scoring Configuration -->
  <Block className="form-section">
    <h3 module="html">Scoring</h3>
    
    <Field 
      type="number" 
      name="MaxScore" 
      label="Maximum Score"
      className="w33" 
      module="formikforms"
      required="true"
      min="1"
      placeholder="100"/>
      
    <Field 
      type="number" 
      name="PassingScore" 
      label="Passing Score"
      className="w33" 
      module="formikforms"
      required="true"
      min="1"
      placeholder="60"/>
      
    <Field 
      type="number" 
      name="Duration" 
      label="Duration (minutes)"
      className="w33" 
      module="formikforms"
      required="true"
      min="1"
      placeholder="120"/>
  </Block>
  
  <!-- Scheduling -->
  <Block className="form-section">
    <h3 module="html">Schedule</h3>
    
    <Field 
      type="datetime" 
      name="StartDate" 
      label="StartDate"
      className="w50" 
      module="formikforms"
      required="true"/>
      
    <Field 
      type="datetime" 
      name="EndDate" 
      label="End Date"
      className="w50" 
      module="formikforms"
      required="true"/>
      
    <Field 
      type="number" 
      name="AttemptsAllowed" 
      label="Attempts Allowed"
      className="w50" 
      module="formikforms"
      required="true"
      min="1"
      max="10"
      default="1"/>
  </Block>
  
  <!-- Status -->
  <Field 
    type="select" 
    name="Status" 
    label="Status"
    className="w50" 
    module="formikforms"
    required="true"
    default="Draft">
    <option value="Draft">Draft</option>
    <option value="Published">Published</option>
    <option value="Closed">Closed</option>
  </Field>
  
  <!-- Submit Button -->
  <Button type="submit" className="btn-primary">
    Create Exam
  </Button>
</Form>
```

## Step 4: Create Form Pages

Create `src/ui/registry/Pages/enroll-student.yml`:

```yaml
route: "/enroll"
authenticate: true
roles: [Teacher, Admin]
component:
  type: layout
  layout: 2col
  leftcol:
    type: menu
    id: main-menu
  rightcol:
    type: form
    id: student-enrollment
    className: form-container
```

Create `src/ui/registry/Pages/create-exam.yml`:

```yaml
route: "/create-exam"
authenticate: true
roles: [Teacher, Admin]
component:
  type: layout
  layout: 2col
  leftcol:
    type: menu
    id: main-menu
  rightcol:
    type: form
    id: exam-creation
    className: form-container
```

## Step 5: Add Form Styles

Update `src/ui/styles/main.scss`:

```scss
// Form styles
.form-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  
  h2 {
    color: #333;
    margin-bottom: 10px;
  }
  
  p {
    color: #666;
    margin-bottom: 30px;
  }
}

.form-section {
  margin-bottom: 30px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 4px;
  
  h3 {
    color: #333;
    margin-bottom: 15px;
    font-size: 18px;
  }
}

// Field widths
.w100 { width: 100%; }
.w50 { width: 48%; display: inline-block; margin-right: 2%; }
.w33 { width: 32%; display: inline-block; margin-right: 1.3%; }

// Button styles
.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 12px 30px;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: transform 0.2s;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
  }
  
  &:active {
    transform: translateY(0);
  }
}

// Field styles
input[type="text"],
input[type="email"],
input[type="number"],
input[type="date"],
input[type="datetime-local"],
select,
textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.2s;
  
  &:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }
  
  &.error {
    border-color: #e74c3c;
  }
}

.field-error {
  color: #e74c3c;
  font-size: 12px;
  margin-top: 5px;
}
```

## Understanding Forms

### Field Types

```xml
<!-- Text Input -->
<Field type="string" name="Name" />

<!-- Email with Validation -->
<Field type="email" name="Email" validation="email" />

<!-- Number with Min/Max -->
<Field type="number" name="Age" min="0" max="120" />

<!-- Date Picker -->
<Field type="date" name="BirthDate" />

<!-- DateTime Picker -->
<Field type="datetime" name="ExamStart" />

<!-- Select Dropdown -->
<Field type="select" name="Grade">
  <option value="9">Grade 9</option>
  <option value="10">Grade 10</option>
</Field>

<!-- Textarea -->
<Field type="textarea" name="Description" rows="5" />

<!-- Checkbox -->
<Field type="boolean" name="AgreeToTerms" />

<!-- File Upload -->
<Field type="file" name="ProfilePicture" accept="image/*" />
```

### Validation

```xml
<!-- Required Field -->
<Field name="Email" required="true" />

<!-- Pattern Matching -->
<Field name="Phone" pattern="^[\d\(\)\-\s]+$" />

<!-- Min/Max Length -->
<Field name="Password" minlength="8" maxlength="50" />

<!-- Email Validation -->
<Field type="email" validation="email" />
```

### Form Actions

```yaml
# Call a method
actiontype: callmethod
actionparams:
  method: submitForm
  
# Call a service directly
actiontype: callservice
actionparams:
  service: MyService
  params:
    field1: "$formData.Field1"
    
# Open a page
actiontype: navigate
actionparams:
  path: /success
  
# Custom action
actiontype: customaction
actionparams:
  handler: myCustomHandler
```

### Form Lifecycle Hooks

```yaml
actionparams:
  method: submitForm
  presubmit: validateForm      # Called before submission
  success: onSuccess           # Called on success
  error: onError               # Called on error
```

```javascript
// In Initialize.js
return {
  validateForm: (formData) => {
    // Return false to prevent submission
    if (!formData.Email.includes('@')) {
      return false;
    }
    return true;
  },
  
  onSuccess: (response) => {
    console.log('Form submitted successfully:', response);
  },
  
  onError: (error) => {
    console.error('Form submission failed:', error);
  }
};
```

## Best Practices

### 1. Group Related Fields

```xml
<Block className="form-section">
  <h3 module="html">Personal Information</h3>
  <Field name="FirstName" />
  <Field name="LastName" />
</Block>

<Block className="form-section">
  <h3 module="html">Contact Information</h3>
  <Field name="Email" />
  <Field name="Phone" />
</Block>
```

### 2. Provide Clear Labels and Placeholders

```xml
<Field 
  name="Email" 
  label="Email Address"
  placeholder="john@example.com"
  help="We'll never share your email"/>
```

### 3. Use Appropriate Input Types

```xml
<!-- Good: Uses specific types -->
<Field type="email" name="Email" />
<Field type="date" name="BirthDate" />
<Field type="number" name="Age" />

<!-- Avoid: Generic string for everything -->
<Field type="string" name="Email" />
<Field type="string" name="BirthDate" />
```


### 4. Use Form-Level Validation

Forms support `required` attributes and validation patterns. Server-side validation is also performed.

```xml
<Field type="email" name="Email" required="true" validation="email" />
```

## Build and Test

```bash
# Rebuild UI plugin
laatoo plugin build studentmgmt-ui --getbuildpackages

# Copy to application
cp -r bin ../../applications/education/config/modules/studentmgmt-ui

# Restart
cd ../../../..
laatoo solution run student-mgmt-system
```

Test the forms:
- Navigate to `http://localhost:8080/enroll`
- Navigate to `http://localhost:8080/create-exam`

## Summary

You've created:
-✅ Student enrollment form with validation
- ✅ Exam creation form with grouped fields
- ✅ Form actions and handlers
- ✅ Form pages with navigation
- ✅ Styled form components

**Next**: [Chapter 7: Views & Lists](07-views-lists.md) - Display and browse data
