# Entity Data Models

Complete guide to defining entities in Laatoo.

## Overview

Entities are YAML-defined data models that automatically generate:
- Database schema
- CRUD operations
- REST API endpoints
- Data validation

## Basic Entity Definition

```yaml
name: Product
description: Product in catalog
objectcn: Product
objectpluralcn: Products
parent: Catalog
instance: product
dataservicefactory: sql
dbtype: sql

attributes:
  - name: ProductId
    type: string
    primary: true
    
  - name: Name
    type: string
    required: true
    
  - name: Price
    type: float
    required: true
    min: 0
    
  - name: Quantity
    type: int
    default: 0
```

## Entity Metadata Properties

| Property | Type | Description |
|----------|------|-------------|
| `name` | string | Name of the entity object. |
| `description` | string | Description of the entity. |
| `objectcn` | string | Common name for the object (human readable). |
| `objectpluralcn` | string | Plural common name for the object. |
| `parent` | string | Parent channel to be used for entity APIs. |
| `instance` | string | Name of the instance (e.g., `myuser` instead of `user`) to be used for the object. Defaults to object name. |
| `dataservicefactory` | string | Data service factory to be used for entity storing (e.g., `sql`, `nosql`). |
| `dbtype` | string | Type of database (`sql`/`nosql`). |
| `list_page_route` | string | Route for the list page on SPA UI. Default: `/list_instance`. |
| `new_form_route` | string | Route for the new entity page on application UI. Default: `/new_form_instance`. |
| `update_form_route` | string | Route for the update entity page on application UI. Default: `/update_form_instance/:entityId`. |
| `view_page_route` | string | Route for the view page on application UI. Default: `/instance/:entityId`. |
| `create_dataservice` | bool | Create data service for object. |
| `create_getservice` | bool | Create get data service. |
| `create_saveservice` | bool | Create save data service. |
| `create_updateservice` | bool | Create update data service. |
| `create_deleteservice` | bool | Create delete data service. |


## Data Types

| Type | SQL Type | Description |
|------|----------|-------------|
| string | VARCHAR(255) | Short text |
| text | TEXT | Long text |
| int | INTEGER | Whole numbers |
| float | DOUBLE PRECISION | Decimals |
| boolean | BOOLEAN | True/false |
| date | DATE | Date only |
| datetime | TIMESTAMP | Date and time |
| json | JSON/JSONB | JSON data |
| uuid | UUID | UUID values |

## Attribute Properties

```yaml
attributes:
  - name: Email
    type: string
    required: true      # Cannot be null
    unique: true        # Must be unique
    minlength: 5        # Minimum length
    maxlength: 255      # Maximum length
    pattern: "^[^@]+@[^@]+$"  # Regex pattern
    default: ""         # Default value
```

## Relationships

### One-to-Many

```yaml
# Order entity
attributes:
  - name: CustomerId
    type: string
    reference: Customer.CustomerId  # Foreign key
    required: true
```

### Many-to-Many

Create junction entity:

```yaml
name: StudentCourse
attributes:
  - name: StudentId
    type: string
    reference: Student.StudentId
    
  - name: CourseId
    type: string
    reference: Course.CourseId
    
constraints:
  - name: pk_studentcourse
    type: primary
    columns: [StudentId, CourseId]
```

## Indexes

```yaml
indexes:
  # Single column
  - name: idx_email
    columns: [Email]
    unique: true
    
  # Composite index
  - name: idx_name_status
    columns: [Name, Status]
    
  # Partial index
  - name: idx_active_products
    columns: [Status]
    where: "Status = 'Active'"
```

## Constraints

```yaml
constraints:
  # Unique constraint
  - name: unique_email
    type: unique
    columns: [Email]
    
  # Check constraint
  - name: check_positive_price
    type: check
    expression: "Price >= 0"
```

## Best Practices

1. **Use Descriptive Names**: `CustomerEmail` not `ce`
2. **Add Indexes for Queries**: Index fields used in WHERE clauses
3. **Denormalize When Needed**: Include frequently-queried foreign fields
4. **Use Appropriate Types**: `datetime` for timestamps, `boolean` for flags
5. **Add References**: Define foreign keys explicitly

## Example: Complete Entity

```yaml
name: Order
description: Customer order

attributes:
  - name: OrderId
    type: uuid
    primary: true
    
  - name: CustomerId
    type: string
    reference: Customer.CustomerId
    required: true
    
  - name: OrderDate
    type: datetime
    required: true
    default: NOW()
    
  - name: Status
    type: string
    required: true
    default: "Pending"
    
  - name: TotalAmount
    type: float
    required: true
    min: 0
    
  - name: ShippingAddress
    type: text
    
indexes:
  - name: idx_customer_date
    columns: [CustomerId, OrderDate]
  - name: idx_status
    columns: [Status]
    
constraints:
  - name: check_total
    type: check
    expression: "TotalAmount >= 0"
```

## Generated CRUD Operations

For each entity, Laatoo creates:

- `<Entity>.Create` - Create new record
- `<Entity>.GetById` - Get by ID
- `<Entity>.Update` - Update record
- `<Entity>.Delete` - Delete record
- `<Entity>.Query` - Query with filters

## Understanding Relationships: StorableRef vs Storable

### When to Use StorableRef

Use `storableref` for **references to independent entities** (like foreign keys):

```yaml
fields:
  Department:
    type: storableref
    entity: departments.Department
```

**Characteristics:**
- Referenced entity exists independently
- Multiple entities can reference the same instance
- Changes to referenced entity reflect everywhere
- Stored as entity ID reference

**Student Management example:**
- Student → Department (many students, one department)
- StudentResult → Course (many results, one course)
- Course → Teacher (many courses, one teacher)

### When to Use Storable

Use `storable` for **embedded/composed data** (part-of relationship):

```yaml
fields:
  Address:
    type: storable
    entity: common.Address
```

**Characteristics:**
- Data is part of parent entity
- No independent existence
- Changes affect only parent entity
- Stored inline with parent

**Student Management example:**
- Student → Address (address is part of student data)
- Enrollment → PaymentInfo (payment details specific to enrollment)

### Student Management Examples

**Student with embedded Address and Department reference:**

```yaml
name: Student
fields:
  StudentId:
    type: string
  Name:
    type: string
  
  # Embedded address (storable)
  Address:
    type: storable
    entity: common.Address
  
  # Department reference (storableref)
  Department:
    type: storableref
    entity: departments.Department
```

**StudentResult linking Student and Course:**

```yaml
name: StudentResult
fields:
  Student:
    type: storableref
    entity: students.Student
  Course:
    type: storableref
    entity: courses.Course
  Grade:
    type: string
  Score:
    type: float
```

## Lifecycle Hooks

### PreSave Hook

Execute logic before saving to database.

**Enable:**
```yaml
name: Student
presave: true
```

**Implement:**
```go
func (s *Student) PreSave(ctx ctx.Context) error {
    // Normalize email
    s.Email = strings.ToLower(s.Email)
    
    // Generate student ID if not provided
    if s.StudentId == "" {
        s.StudentId = generateStudentId()
    }
    
    return nil
}
```

**Use cases:**
- Data validation
- Field normalization
- Auto-generating IDs
- Password hashing

### PostSave Hook

Execute logic after successful save.

**Enable:**
```yaml
name: StudentResult
postsave: true
```

**Implement:**
```go
func (r *StudentResult) PostSave(ctx ctx.Context) error {
    // Send notification about new grade
    return sendGradeNotification(ctx, r.Student.Id, r.Course, r.Grade)
}
```

**Use cases:**
- Sending notifications
- Updating related entities
- Logging audit trail
- Triggering workflows

### PostLoad Hook

Execute logic after loading from database.

**Enable:**
```yaml
name: Student
postload: true
```

**Implement:**
```go
func (s *Student) PostLoad(ctx ctx.Context) error {
    // Calculate derived fields
    s.Age = calculateAge(s.DateOfBirth)
    return nil
}
```

**Use cases:**
- Calculating derived fields
- Decrypting sensitive data
- Loading additional context

## See Also

- [Services Guide](services.md) - Create custom services
- [Tutorial: Data Model](../../tutorials/student-management/02-data-model.md)
