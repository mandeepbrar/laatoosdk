# Creating Datasets

Datasets are declarative server-side data query definitions that provide structured access to entity collections. They enable you to define filtered, parameterized views of your data that can be consumed by UI components.

## What are Datasets?

Datasets serve as reusable data access layers that:
- Define filtered and transformed views of entity data
- Support parameterized queries with runtime filters
- Enable data projection through DTOs (Data Transfer Objects)
- Provide consistent data access patterns for UI components
- Support multi-tenancy and access control
- Integrate seamlessly with client-side Views

## Use Cases

**Common scenarios for datasets:**
- **List Views**: Populate tables and lists with filtered entity data
- **Dropdown Data**: Provide options for select/dropdown widgets
- **Search Results**: Return filtered and sorted entity collections
- **Reports**: Generate structured datasets for reporting

## Dataset Structure

Datasets are defined as YAML files in the `src/server/datasets/` directory of your plugin.

### File Location

```
dev/plugins/studentmgmt/
└── src/
    └── server/
        └── datasets/
            ├── studentresults.yml
            ├── courses.yml
            └── enrollments.yml
```

### Basic Definition

```yaml
Entity: <module>.<EntityName>
Params:
  <paramName>: <type>
Dao: <module>.<DtoName>
Properties:
  <fieldName>: <type>
Filters:
  - Lhs: <entity.property>
    LhsType: Property
    Rhs: <paramName>
    RhsType: Param
Sort: "<field> <asc|desc>"
```

## Core Properties

### Entity

**Required**. Specifies the source entity to query.

```yaml
Entity: studentmgmt.StudentResult
```

The entity name must be fully qualified: `<modulename>.<EntityName>`

### Params

**Optional**. Defines runtime parameters for filtering.

```yaml
Params:
  StudentId: string
  CourseId: string
  Semester: string
```

**Supported types**: `string`, `int`, `float`, `bool`, `date`

### Dao

**Optional**. Specifies a DTO for data projection.

```yaml
Dao: studentmgmt.ResultSummaryDTO
```

When specified, only the DTO properties are returned, providing:
- Field selection (return only needed fields)
- Performance optimization (reduced payload)
- Security (hide sensitive fields)

### Properties

**Required when using DAO**. Defines fields to include in the response.

```yaml
Properties:
  StudentName: string
  CourseName: string
  Grade: string
  Score: float
```

### Filters

**Optional**. Defines query conditions.

```yaml
Filters:
  - Lhs: Student.Id
    LhsType: Property
    Rhs: StudentId
    RhsType: Param
  - Lhs: Course.Id
    LhsType: Property
    Rhs: CourseId
    RhsType: Param
```

**Filter properties:**
- `Lhs`: Left-hand side (entity property path)
- `LhsType`: `Property` or `Constant`
- `Rhs`: Right-hand side (param name or value)
- `RhsType`: `Param`, `Constant`, or `Property`
- `Op`: Operator (`==`, `!=`, `>`, `<`, `>=`, `<=`, `in`, `contains`)

### Sort

**Optional**. Specifies default sorting.

```yaml
Sort: "StudentName asc, Score desc"
```

### Cache

**Optional**. Enables caching for improved performance.

```yaml
Cache: true
```

Use caching for slowly-changing reference data like courses or semesters.

### Expand

**Optional**. Auto-populates reference fields.

```yaml
Expand:
  Student: false      # Single reference
  Course: false       # Single reference
```

When enabled, referenced entities are automatically loaded and included in the response.

## Student Management Examples

### Example 1: Student Results Dataset

Display all results for a specific student.

**File**: `src/server/datasets/studentresults.yml`

```yaml
Entity: studentmgmt.StudentResult
Params:
  StudentId: string
Dao: studentmgmt.ResultSummaryDTO
Properties:
  CourseName: string
  CourseCode: string
  Grade: string
  Score: float
  Semester: string
  CreditsEarned: int
Filters:
  - Lhs: Student.Id
    LhsType: Property
    Rhs: StudentId
    RhsType: Param
Sort: "Semester desc, CourseName asc"
Expand:
  Course: false
```

**Features:**
- Filters results by StudentId parameter
- Projects data through ResultSummaryDTO
- Expands Course reference for full course details
- Sorts by semester (newest first) then course name

### Example 2: Course Results Dataset

Display all student results for a specific course.

**File**: `src/server/datasets/courseresults.yml`

```yaml
Entity: studentmgmt.StudentResult
Params:
  CourseId: string
  Semester: string
Dao: studentmgmt.CourseResultsDTO
Properties:
  StudentName: string
  StudentId: string
  Grade: string
  Score: float
  Status: string
Filters:
  - Lhs: Course.Id
    LhsType: Property
    Rhs: CourseId
    RhsType: Param
  - Lhs: Semester
    LhsType: Property
    Rhs: Semester
    RhsType: Param
Sort: "StudentName asc"
Expand:
  Student: false
```

**Features:**
- Filters by both course and semester
- Returns student information for each result
- Expands Student reference
- Sorted alphabetically by student name

### Example 3: Active Courses Dataset

List all active courses for dropdown selection.

**File**: `src/server/datasets/activecourses.yml`

```yaml
Entity: studentmgmt.Course
Dao: studentmgmt.CourseListDTO
Properties:
  CourseId: string
  CourseName: string
  CourseCode: string
  Credits: int
  Department: string
Filters:
  - Lhs: Active
    LhsType: Property
    Rhs: "true"
    RhsType: Constant
Sort: "Department asc, CourseName asc"
Cache: true
```

**Features:**
- Constant filter for active courses only
- Cached for better performance
- Sorted by department and name
- Ideal for dropdown/select widgets

### Example 4: Grade Statistics Dataset

Advanced dataset with multiple parameters and operators.

**File**: `src/server/datasets/gradestatistics.yml`

```yaml
Entity: studentmgmt.StudentResult
Params:
  Department: string
  Semester: string
  MinScore: float
Dao: studentmgmt.GradeStatsDTO
Properties:
  CourseName: string
  Department: string
  AverageScore: float
  PassRate: float
  StudentCount: int
Filters:
  - Lhs: Course.Department
    LhsType: Property
    Rhs: Department
    RhsType: Param
  - Lhs: Semester
    LhsType: Property
    Rhs: Semester
    RhsType: Param
  - Lhs: Score
    LhsType: Property
    Rhs: MinScore
    RhsType: Param
    Op: ">="
```

**Features:**
- Multiple parameters for flexible filtering
- Uses comparison operator for minimum score
- Nested property filtering (Course.Department)

## Integration with UI Views

Datasets are consumed by UI Views through the `dataset` property.

### View Configuration

**File**: `src/ui/registry/Views/studentresults.yml`

```yaml
title: My Results
dataset: studentresults
filters:
  StudentId: [[ctx.user.id]]
item:
  type: block
  id: resultcard
```

**How it works:**
1. View references dataset by name
2. Parameters passed through `filters`
3. Dataset executes with provided parameters
4. Results rendered using item template

### Dynamic Parameters

Use context variables for dynamic filtering:

```yaml
dataset: courseresults
filters:
  CourseId: [[Application.Route.Params.courseId]]
  Semester: [[ctx.pageContext.selectedSemester]]
```

**Available context:**
- `ctx.user.*` - User information
- `ctx.pageContext.*` - Page-level context
- `Application.Route.Params.*` - URL parameters

## Creating DTOs

DTOs are defined as entities with selected fields.

**File**: `build/entities/ResultSummaryDTO.yml`

```yaml
name: ResultSummaryDTO
fields:
  CourseName:
    type: string
  CourseCode:
    type: string
  Grade:
    type: string
  Score:
    type: float
  Semester:
    type: string
```

**Best practices:**
- Suffix names with `DTO`
- Include only necessary fields
- Flatten nested relationships
- Exclude sensitive data

## Best Practices

1. **Use DTOs for Production**: Always define DTOs to control the data shape
2. **Enable Caching Wisely**: Cache reference data, not frequently-changing data
3. **Index Filter Fields**: Add database indexes on filtered properties
4. **Limit Expand Usage**: Only expand necessary references
5. **Descriptive Names**: Use clear dataset names matching their purpose
6. **Parameter Validation**: Define appropriate parameter types

## See Also

- [Entities Guide](../server-development/entities.md) - Define entities and DTOs
- [Views Guide](views.md) - Consume datasets in UI
- [Forms Guide](forms.md) - Use datasets in form selects
