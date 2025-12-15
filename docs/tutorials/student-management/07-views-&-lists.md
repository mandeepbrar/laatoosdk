# Chapter 7: Views & Lists

Display and browse students, exams, and results using views.

## Overview

Views in Laatoo provide:
- Data display with automatic pagination
- Filtering and sorting
- Custom item rendering
- Actions on items
- Integration with datasets

We'll create:
- **Student List View**: Browse all students
- **Exam List View**: View available exams
- **Results List View**: Display exam results

## Step 1: Create Student Dataset

Create `src/server/datasets/student_list.yml`:

```yaml
Entity: studentmgmt.Student
Params:
  Status: string
  Grade: string
Filters:
  - Lhs: Status
    LhsType: Property
    Rhs: Status
    RhsType: Param
  - Lhs: Grade
    LhsType: Property
    Rhs: Grade
    RhsType: Param
Sort: "Name asc"
```

## Step 2: Create Student List View

Create `src/ui/registry/Views/student-list.yml`:

```yaml
# Student list view configuration
dataset: student_list
overlay: false
pagination: true
pageSize: 20

# Header configuration
header:
  title: "Students"
  actions:
    - label: "Enroll New Student"
      action: navigateToEnroll
      icon: "fa fa-plus"
      className: "btn-primary"

# Filter configuration
filters:
  - name: Status
    field: Status
    type: select
    options:
      - value: "Active"
        label: "Active"
      - value: "Suspended"
        label: "Suspended"
      - value: "Graduated"
        label: "Graduated"
        
  - name: Grade
    field: Grade
    type: select
    options:
      - value: "9"
      - value: "10"
      - value: "11"
      - value: "12"

# Search configuration
search:
  enabled: true
  placeholder: "Search by name or email..."
  fields: [Name, Email]

# Item configuration
item:
  type: custom
  blockId: student-card

# Sorting
sort:
  default: Name
  options:
    - field: Name
      label: "Name"
    - field: EnrollmentDate
      label: "Enrollment Date"
    - field: Grade
      label: "Grade"
```

## Step 3: Create Student Card Block

Create `src/ui/registry/Blocks/student-card.xml`:

```xml
<Block className="student-card">
  <Block className="student-header">
    <h3 module="html">{Name}</h3>
    <Block className="student-status status-{Status}">{Status}</Block>
  </Block>
  
  <Block className="student-details">
    <Block className="detail-row">
      <span className="label" module="html">Email:</span>
      <span className="value" module="html">{Email}</span>
    </Block>
    
    <Block className="detail-row">
      <span className="label" module="html">Grade:</span>
      <span className="value" module="html">{Grade}</span>
    </Block>
    
    <Block className="detail-row">
      <span className="label" module="html">Enrolled:</span>
      <span className="value" module="html">{EnrollmentDate|date}</span>
    </Block>
  </Block>
  
  <Block className="student-actions">
    <Action name="viewStudent" data="{StudentId}">
      <Content>View Profile</Content>
    </Action>
    <Action name="editStudent" data="{StudentId}">
      <Content>Edit</Content>
    </Action>
  </Block>
</Block>
```

## Step 4: Create Exam List View

Create `src/ui/registry/Views/exam-list.yml`:

```yaml
serviceName: Exam.Query
overlay: false
pagination: true
pageSize: 12

header:
  title: "Exams"
  actions:
    - label: "Create New Exam"
      action: navigateToCreateExam
      icon: "fa fa-plus"
      className: "btn-primary"
      roles: [Teacher, Admin]

filters:
  - name: Status
    field: Status
    type: select
    options:
      - value: "Published"
        label: "Published"
      - value: "Draft"
        label: "Draft"
      - value: "Closed"
        label: "Closed"

search:
  enabled: true
  placeholder: "Search exams..."
  fields: [Title, Description]

item:
  type: custom
  blockId: exam-card

sort:
  default: StartDate
  direction: desc
  options:
    - field: Title
      label: "Title"
    - field: StartDate
      label: "Start Date"
    - field: Duration
      label: "Duration"
```

## Step 5: Create Exam Card Block

Create `src/ui/registry/Blocks/exam-card.xml`:

```xml
<Block className="exam-card">
  <Block className="exam-header">
    <h3 module="html">{Title}</h3>
    <Block className="exam-status status-{Status}">{Status}</Block>
  </Block>
  
  <Block className="exam-description">
    <p module="html">{Description}</p>
  </Block>
  
  <Block className="exam-metadata">
    <Block className="metadata-item">
      <i className="fa fa-clock" module="html"></i>
      <span module="html">{Duration} minutes</span>
    </Block>
    
    <Block className="metadata-item">
      <i className="fa fa-star" module="html"></i>
      <span module="html">{MaxScore} points</span>
    </Block>
    
    <Block className="metadata-item">
      <i className="fa fa-calendar" module="html"></i>
      <span module="html">{StartDate|date} - {EndDate|date}</span>
    </Block>
    
    <Block className="metadata-item">
      <i className="fa fa-redo" module="html"></i>
      <span module="html">{AttemptsAllowed} attempts</span>
    </Block>
  </Block>
  
  <Block className="exam-actions">
    <Action name="takeExam" data="{ExamId}" roles="[Student]">
      <Content>Take Exam</Content>
    </Action>
    <Action name="viewExamDetails" data="{ExamId}">
      <Content>View Details</Content>
    </Action>
    <Action name="editExam" data="{ExamId}" roles="[Teacher,Admin]">
      <Content>Edit</Content>
    </Action>
  </Block>
</Block>
```

## Step 6: Create Results List View

Create `src/ui/registry/Views/results-list.yml`:

```yaml
serviceName: Result.Query
overlay: false
pagination: true
pageSize: 25

header:
  title: "Exam Results"

filters:
  - name: Passed
    field: Passed
    type: select
    options:
      - value: "true"
        label: "Passed"
      - value: "false"
        label: "Failed"
        
  - name: Grade
    field: Grade
    type: select
    options:
      - value: "A"
      - value: "B"
      - value: "C"
      - value: "D"
      - value: "F"

item:
  type: table
  columns:
    - field: StudentId
      label: "Student ID"
      width: "15%"
      
    - field: ExamId
      label: "Exam"
      width: "20%"
      
    - field: Score
      label: "Score"
      width: "10%"
      render: "{Score}/{MaxScore}"
      
    - field: Percentage
      label: "Percentage"
      width: "10%"
      render: "{Percentage}%"
      
    - field: Grade
      label: "Grade"
      width: "10%"
      className: "grade-{Grade}"
      
    - field: Passed
      label: "Status"
      width: "10%"
      render: "{Passed|yesno:Passed,Failed}"
      
    - field: ProcessedDate
      label: "Date"
      width: "15%"
      render: "{ProcessedDate|date}"
      
    - field: actions
      label: "Actions"
      width: "10%"
      type: actions
      actions:
        - label: "View"
          action: viewResult
          data: "{ResultId}"

sort:
  default: ProcessedDate
  direction: desc
```

## Step 7: Create View Pages

Create `src/ui/registry/Pages/students.yml`:

```yaml
route: "/students"
authenticate: true
roles: [Teacher, Admin]
component:
  type: layout
  layout: 2col
  leftcol:
    type: menu
    id: main-menu
  rightcol:
    type: view
    id: student-list
    className: content-area
```

Create `src/ui/registry/Pages/exams.yml`:

```yaml
route: "/exams"
authenticate: true
component:
  type: layout
  layout: 2col
  leftcol:
    type: menu
    id: main-menu
  rightcol:
    type: view
    id: exam-list
    className: content-area
```

Create `src/ui/registry/Pages/results.yml`:

```yaml
route: "/results"
authenticate: true
roles: [Teacher, Admin]
component:
  type: layout
  layout: 2col
  leftcol:
    type: menu
    id: main-menu
  rightcol:
    type: view
    id: results-list
    className: content-area
```

## Step 8: Add View Styles

Update `src/ui/styles/main.scss`:

```scss
// Content area
.content-area {
  padding: 20px;
  background: #f5f5f5;
  min-height: 100vh;
}

// Student Card
.student-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 15px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }
}

.student-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  
  h3 {
    margin: 0;
    color: #333;
  }
}

.student-status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  
  &.status-Active {
    background: #d4edda;
    color: #155724;
  }
  
  &.status-Suspended {
    background: #f8d7da;
    color: #721c24;
  }
  
  &.status-Graduated {
    background: #d1ecf1;
    color: #0c5460;
  }
}

.student-details {
  .detail-row {
    display: flex;
    margin-bottom: 8px;
    
    .label {
      font-weight: 600;
      color: #666;
      min-width: 100px;
    }
    
    .value {
      color: #333;
    }
  }
}

.student-actions {
  display: flex;
  gap: 10px;
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #eee;
  
  button {
    padding: 6px 16px;
    border: 1px solid #ddd;
    background: white;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.2s;
    
    &:hover {
      background: #667eea;
      color: white;
      border-color: #667eea;
    }
  }
}

// Exam Card
.exam-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 15px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.exam-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  
  h3 {
    margin: 0;
    color: #333;
  }
}

.exam-description {
  color: #666;
  margin-bottom: 15px;
  
  p {
    margin: 0;
  }
}

.exam-metadata {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
  margin-bottom: 15px;
  
  .metadata-item {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #666;
    font-size: 14px;
    
    i {
      color: #667eea;
    }
  }
}

// Results Table
.results-table {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  
  table {
    width: 100%;
    border-collapse: collapse;
    
    thead {
      background: #f8f9fa;
      
      th {
        padding: 12px;
        text-align: left;
        font-weight: 600;
        color: #333;
        border-bottom: 2px solid #dee2e6;
      }
    }
    
    tbody {
      tr {
        border-bottom: 1px solid #eee;
        transition: background 0.2s;
        
        &:hover {
          background: #f8f9fa;
        }
        
        td {
          padding: 12px;
          color: #666;
        }
      }
    }
  }
}

// Grade badges
.grade-A { color: #28a745; font-weight: bold; }
.grade-B { color: #17a2b8; font-weight: bold; }
.grade-C { color: #ffc107; font-weight: bold; }
.grade-D { color: #fd7e14; font-weight: bold; }
.grade-F { color: #dc3545; font-weight: bold; }

// Pagination
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  margin-top: 20px;
  
  button {
    padding: 8px 12px;
    border: 1px solid #ddd;
    background: white;
    border-radius: 4px;
    cursor: pointer;
    
    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
    
    &.active {
      background: #667eea;
      color: white;
      border-color: #667eea;
    }
  }
}
```

## Understanding Views

### View Types

**1. Custom Item View**:
```yaml
item:
  type: custom
  blockId: student-card  # Uses custom block
```

**2. Table View**:
```yaml
item:
  type: table
  columns:
    - field: Name
      label: "Student Name"
      width: "30%"
```

**3. Grid View**:
```yaml
item:
  type: grid
  columns: 3  # Number of columns
  blockId: item-card
```

### Filters

```yaml
filters:
  # Select filter
  - name: Status
    field: Status
    type: select
    options:
      - value: "Active"
        label: "Active Students"
        
  # Date range filter
  - name: DateRange
    field: EnrollmentDate
    type: daterange
    
  # Number range filter
  - name: ScoreRange
    field: Score
    type: numberrange
    min: 0
    max: 100
```

### Search

```yaml
search:
  enabled: true
  placeholder: "Search..."
  fields: [Name, Email, Phone]  # Fields to search in
  minLength: 3                   # Minimum characters to trigger
```

### Sorting

```yaml
sort:
  default: Name        # Default sort field
  direction: asc       # asc or desc
  options:
    - field: Name
      label: "Name (A-Z)"
    - field: EnrollmentDate
      label: "Newest First"
```

### Pagination

```yaml
pagination: true
pageSize: 20
pageSizeOptions: [10, 20, 50, 100]
```

## Build and Test

```bash
# Rebuild
laatoo plugin build studentmgmt-ui --getbuildpackages

# Deploy
cp -r bin ../../applications/education/config/modules/studentmgmt-ui

# Restart
cd ../../../..
laatoo solution run student-mgmt-system
```

Test the views:
- `/students` - Student list
- `/exams` - Exam list
- `/results` - Results table

## Summary

You've created:
- ✅ Student list view with custom cards
- ✅ Exam list view with metadata
- ✅ Results table view
- ✅ Filters, search, and sorting
- ✅ Item actions and navigation
- ✅ Responsive styles

**Next**: [Chapter 8: Dashboards](08-dashboards.md) - Create analytics dashboards
