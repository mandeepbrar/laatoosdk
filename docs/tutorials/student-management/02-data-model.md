# Chapter 2: Data Model

Define the entities for your Student Management System.

## Overview

In this chapter, you'll create the data model for the Student Management System. Laatoo uses YAML-based entity definitions that automatically generate:
- Database schema
- CRUD operations
- REST API endpoints
- Data validation

## Entity Design

Our system needs four core entities:

```
Student → registers for → Exam
Student → takes → ExamAttempt → generates → Result
```

## Step 1: Create Entity Definitions Directory

```bash
cd dev/plugins/studentmgmt
mkdir -p config/entities
```

## Step 2: Define Student Entity

Create `config/entities/Student.yml`:

```yaml
# Student Entity Definition
name: Student
description: Student enrolled in the system

attributes:
  - name: StudentId
    type: string
    primary: true
    description: Unique student identifier
    
  - name: Name
    type: string
    required: true
    description: Full name of the student
    
  - name: Email
    type: string
    required: true
    unique: true
    description: Student email address
    
  - name: EnrollmentDate
    type: datetime
    required: true
    description: Date of enrollment
    
  - name: Status
    type: string
    default: "Active"
    description: Student status (Active, Suspended, Graduated)
    
  - name: Grade
    type: string
    description: Current grade level
    
  - name: DateOfBirth
    type: date
    description: Student's date of birth
    
  - name: Phone
    type: string
    description: Contact phone number

# Define indexes for common queries
indexes:
  - name: idx_email
    columns: [Email]
    unique: true
  - name: idx_status
    columns: [Status]
```

## Step 3: Define Exam Entity

Create `config/entities/Exam.yml`:

```yaml
# Exam Entity Definition
name: Exam
description: Examination configuration

attributes:
  - name: ExamId
    type: string
    primary: true
    description: Unique exam identifier
    
  - name: Title
    type: string
    required: true
    description: Exam title
    
  - name: Description
    type: text
    description: Detailed exam description
    
  - name: Duration
    type: int
    required: true
    description: Exam duration in minutes
    
  - name: MaxScore
    type: int
    required: true
    description: Maximum achievable score
    
  - name: PassingScore
    type: int
    required: true
    description: Minimum score to pass
    
  - name: AttemptsAllowed
    type: int
    default: 1
    description: Number of attempts allowed per student
    
  - name: StartDate
    type: datetime
    required: true
    description: Exam availability start date
    
  - name: EndDate
    type: datetime
    required: true
    description: Exam availability end date
    
  - name: Status
    type: string
    default: "Draft"
    description: Exam status (Draft, Published, Closed)
    
  - name: CreatedBy
    type: string
    description: User who created the exam
    
  - name: CreatedDate
    type: datetime
    description: Creation timestamp

indexes:
  - name: idx_status_dates
    columns: [Status, StartDate, EndDate]
```

## Step 4: Define ExamAttempt Entity

Create `config/entities/ExamAttempt.yml`:

```yaml
# ExamAttempt Entity Definition
name: ExamAttempt
description: Student's exam attempt record

attributes:
  - name: AttemptId
    type: string
    primary: true
    description: Unique attempt identifier
    
  - name: StudentId
    type: string
    required: true
    reference: Student.StudentId
    description: Reference to Student
    
  - name: ExamId
    type: string
    required: true
    reference: Exam.ExamId
    description: Reference to Exam
    
  - name: AttemptNumber
    type: int
    required: true
    description: Attempt number (1, 2, 3...)
    
  - name: StartTime
    type: datetime
    required: true
    description: When student started the exam
    
  - name: SubmitTime
    type: datetime
    description: When student submitted the exam
    
  - name: Status
    type: string
    default: "InProgress"
    description: Attempt status (InProgress, Submitted, Graded)
    
  - name: Answers
    type: json
    description: Student's answers (JSON format)
    
  - name: TimeSpent
    type: int
    description: Time spent in minutes

indexes:
  - name: idx_student_exam
    columns: [StudentId, ExamId]
  - name: idx_status
    columns: [Status]

# Ensure unique attempt numbers per student per exam
constraints:
  - name: unique_attempt
    type: unique
    columns: [StudentId, ExamId, AttemptNumber]
```

## Step 5: Define Result Entity

Create `config/entities/Result.yml`:

```yaml
# Result Entity Definition
name: Result
description: Exam result after grading

attributes:
  - name: ResultId
    type: string
    primary: true
    description: Unique result identifier
    
  - name: AttemptId
    type: string
    required: true
    reference: ExamAttempt.AttemptId
    unique: true
    description: Reference to ExamAttempt (one-to-one)
    
  - name: StudentId
    type: string
    required: true
    reference: Student.StudentId
    description: Reference to Student (denormalized for queries)
    
  - name: ExamId
    type: string
    required: true
    reference: Exam.ExamId
    description: Reference to Exam (denormalized for queries)
    
  - name: Score
    type: int
    required: true
    description: Achieved score
    
  - name: MaxScore
    type: int
    required: true
    description: Maximum possible score
    
  - name: Percentage
    type: float
    description: Score percentage (calculated)
    
  - name: Grade
    type: string
    description: Letter grade (A, B, C, D, F)
    
  - name: Passed
    type: boolean
    required: true
    description: Whether student passed
    
  - name: ProcessedDate
    type: datetime
    required: true
    description: When result was generated
    
  - name: Feedback
    type: text
    description: Teacher feedback
    
  - name: ProcessedBy
    type: string
    description: Who processed/graded the exam

indexes:
  - name: idx_student
    columns: [StudentId]
  - name: idx_exam
    columns: [ExamId]
  - name: idx_passed
    columns: [Passed]
```

## Understanding Entity Attributes

### Attribute Types

Laatoo supports these data types:

| Type | Description | Example |
|------|-------------|---------|
| `string` | Variable-length text | "John Doe" |
| `text` | Long text/paragraphs | Descriptions |
| `int` | Integer number | 100 |
| `float` | Decimal number | 95.5 |
| `boolean` | True/false | true |
| `date` | Date only | 2024-01-15 |
| `datetime` | Date and time | 2024-01-15T10:30:00Z |
| `json` | JSON object | {"answer": "A"} |

### Attribute Properties

- **`primary: true`**: Marks as primary key
- **`required: true`**: Cannot be null
- **`unique: true`**: Must be unique across records
- **`default: value`**: Default value if not provided
- **`reference: Entity.Field`**: Foreign key relationship
- **`description`**: Documentation for the field

### Indexes and Constraints

```yaml
# Single column index
indexes:
  - name: idx_email
    columns: [Email]
    unique: true

# Composite index
indexes:
  - name: idx_student_exam
    columns: [StudentId, ExamId]

# Unique constraint
constraints:
  - name: unique_attempt
    type: unique
    columns: [StudentId, ExamId, AttemptNumber]
```

## Step 6: Update Plugin Configuration

Edit `config/config.yml` to include entities:

```yaml
name: studentmgmt
version: 1.0.0
description: Student Management System Plugin

entities:
  - Student
  - Exam
  - ExamAttempt
  - Result

# Data adapter (will be configured in application)
dataadapter: default
```

## Step 7: Build the Plugin

```bash
cd dev/plugins/studentmgmt
laatoo plugin build studentmgmt
```

The build process will:
1. Read entity definitions
2. Generate database schema
3. Create CRUD services automatically
4. Package everything into `bin/`

## What Gets Generated

For each entity, Laatoo automatically creates:

### Database Schema
```sql
CREATE TABLE Student (
    StudentId VARCHAR(255) PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Email VARCHAR(255) NOT NULL UNIQUE,
    EnrollmentDate TIMESTAMP NOT NULL,
    Status VARCHAR(255) DEFAULT 'Active',
    Grade VARCHAR(255),
    DateOfBirth DATE,
    Phone VARCHAR(255)
);

CREATE INDEX idx_email ON Student(Email);
CREATE INDEX idx_status ON Student(Status);
```

### CRUD Services
- `Student.Create` - Create new student
- `Student.GetById` - Get student by ID
- `Student.Update` - Update student
- `Student.Delete` - Delete student
- `Student.Query` - Query with filters

### REST API Endpoints
- `POST /api/Student` - Create
- `GET /api/Student/{id}` - Read
- `PUT /api/Student/{id}` - Update
- `DELETE /api/Student/{id}` - Delete
- `GET /api/Student?filter=...` - Query

## Step 8: Install and Configure

```bash
# Copy plugin to application
cp -r bin ../../applications/education/config/modules/studentmgmt

# Configure database adapter
# Edit applications/education/isolations/school-a/config/isolation.yml
```

Add to `isolation.yml`:

```yaml
name: school-a
description: School A isolation

# Database configuration
database:
  type: postgres
  host: localhost
  port: 5432
  database: school_a
  user: laatoo
  password: laatoo

# Data adapter configuration
dataadapters:
  default:
    type: sqldatabase
    database: school_a
```

## Step 9: Run Migrations

When you start the solution, Laatoo will:
1. Connect to the database
2. Compare entity definitions with existing schema
3. Generate and run migrations automatically

```bash
cd ../../../..
laatoo solution run student-mgmt-system
```

Check the logs for:
```
[INFO] Creating table: Student
[INFO] Creating table: Exam
[INFO] Creating table: ExamAttempt
[INFO] Creating table: Result
[INFO] Migrations completed successfully
```

## Testing the Data Model

### Using REST API

```bash
# Create a student
curl -X POST http://localhost:8080/api/Student \
  -H "Content-Type: application/json" \
  -d '{
    "StudentId": "STU001",
    "Name": "Alice Johnson",
    "Email": "alice@example.com",
    "EnrollmentDate": "2024-01-15T00:00:00Z",
    "Grade": "10"
  }'

# Get student by ID
curl http://localhost:8080/api/Student/STU001

# Query students
curl "http://localhost:8080/api/Student?filter=Grade eq '10'"

# Create an exam
curl -X POST http://localhost:8080/api/Exam \
  -H "Content-Type: application/json" \
  -d '{
    "ExamId": "EX001",
    "Title": "Mathematics Final",
    "Duration": 120,
    "MaxScore": 100,
    "PassingScore": 60,
    "StartDate": "2024-06-01T09:00:00Z",
    "EndDate": "2024-06-01T11:00:00Z",
    "Status": "Published"
  }'
```

## Data Relationships

Understanding the relationships:

```
Student (1) ←→ (N) ExamAttempt
Exam (1) ←→ (N) ExamAttempt
ExamAttempt (1) ←→ (1) Result
```

- One student can have many exam attempts
- One exam can have many attempts (from different students)
- Each attempt generates exactly one result

## Best Practices

### 1. Use Meaningful IDs
```yaml
# Good: Descriptive IDs
StudentId: "STU001"
ExamId: "MATH-FINAL-2024"

# Avoid: Generic UUIDs in user-facing fields
StudentId: "a1b2c3d4-e5f6..."
```

### 2. Add Constraints
```yaml
# Ensure data integrity
constraints:
  - name: unique_student_exam_attempt
    type: unique
    columns: [StudentId, ExamId, AttemptNumber]
```

### 3. Use Appropriate Types
```yaml
# Use specific types
Duration: int          # minutes as integer
Percentage: float      # 95.5
Passed: boolean        # true/false
```

### 4. Add Indexes for Common Queries
```yaml
# Index fields used in WHERE clauses
indexes:
  - name: idx_status_dates
    columns: [Status, StartDate, EndDate]
```

### 5. Denormalize When Needed
```yaml
# In Result entity, include both AttemptId and StudentId
# for faster queries even though StudentId can be derived
StudentId:
  type: string
  reference: Student.StudentId
  description: Denormalized for query performance
```

## Troubleshooting

**Issue**: Entity not found during build

**Solution**: Ensure entity is listed in `config/config.yml`:
```yaml
entities:
  - Student
  - Exam
```

---

**Issue**: Database migration fails

**Solution**: Check database connection in `isolation.yml` and ensure database exists:
```bash
# Create database if it doesn't exist
createdb school_a
```

---

**Issue**: Foreign key constraint violation

**Solution**: Ensure referenced entities exist before creating relationships:
```bash
# Create Student first
POST /api/Student {...}

# Then create ExamAttempt
POST /api/ExamAttempt {
  "StudentId": "STU001",  # Must exist
  ...
}
```

## Summary

You've successfully:
- ✅ Defined 4 entities: Student, Exam, ExamAttempt, Result
- ✅ Set up relationships between entities
- ✅ Configured indexes and constraints
- ✅ Built the plugin
- ✅ Tested CRUD operations

**Next**: [Chapter 3: Server Plugin](03-server-plugin.md) - Create custom services for business logic
