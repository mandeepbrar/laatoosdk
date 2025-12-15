# Security Configuration

This guide covers how to configure authentication and authorization in your Laatoo application, including user management, roles, and permissions.

## Overview

Laatoo provides a comprehensive security framework with:

1. **Authentication**: JWT-based token authentication
2. **Authorization**: Two-level access control (channel and service permissions)
3. **Role-Based Access**: Assign roles to users with specific permissions
4. **Multi-Tenancy**: Tenant-based data isolation

## Authentication Setup

### User Entity

Define a User entity to store user credentials and profile information.

**File**: `build/entities/User.yml`

```yaml
name: User
titleField: Username
multitenant: true
cacheable: true
presave: true
fields:
  Username:
    type: string
  Password:
    type: string
    widget:
      props:
        password: true
  Email:
    type: string
  Name:
    type: string
  Status:
    type: int
  Roles:
    type: storableref
    entity: role.Role
    list: true
```

**Key features:**
- `presave: true` - Enables PreSave hook for password hashing
- `Roles` - List of role references for user permissions

### Role Entity

Define roles that group permissions.

**File**: `build/entities/Role.yml`

```yaml
name: Role
titleField: Name
fields:
  Name:
    type: string
  Description:
    type: string
  Permissions:
    type: string
    list: true
```

### Permission Entity

Define individual permissions.

**File**: `build/entities/Permission.yml`

```yaml
name: Permission
fields:
  Name:
    type: string
  Description:
    type: string
```

## Two-Level Security Model

Laatoo implements two layers of access control:

### Level 1: Channel Permissions

Controls access to API endpoints based on roles and permissions.

**Channel Configuration** (`config/server/channels/grades.yml`):

```yaml
id: grades_api
path: /api/grades/:studentId
service: grades_service
permission: grades.read
method: GET
```

**How it works:**
1. User makes request to `/api/grades/123`
2. Security handler checks if user has `grades.read` permission
3. If authorized, request proceeds to service
4. If not, request rejected with 403 Forbidden

### Level 2: Service Access Control

Controls which roles can access specific internal services using Casbin policies.

**Directory Structure:**

```
applications/studentapp/
└── config/
    └── security/
        ├── serviceaccess/
        │   ├── Admin.csv
        │   ├── Teacher.csv
        │   └── Student.csv
        └── rolepermissions/
            ├── Admin.csv
            ├── Teacher.csv
            └── Student.csv
```

**Service Access Format** (CSV):

```csv
p,<role>,<tenant>,<service>
```

**Role Inheritance:**

```csv
g,<child_role>,<parent_role>
```

## Student Management System Example

### Defining Roles

**Admin Role** - Full system access

**File**: `config/security/serviceaccess/Admin.csv`

```csv
p,Admin,*,*
```

**Teacher Role** - Manage courses and grades

**File**: `config/security/serviceaccess/Teacher.csv`

```csv
g,Teacher,Student
p,Teacher,*,dataadapter.*.Course
p,Teacher,*,dataadapter.*.StudentResult
p,Teacher,*,courses.*
p,Teacher,*,grades.*
```

**Student Role** - View own data only

**File**: `config/security/serviceaccess/Student.csv`

```csv
p,Student,*,profile.read
p,Student,*,profile.update.own
p,Student,*,grades.read.own
p,Student,*,courses.list
p,Student,*,enrollment.create
```

### Role Permissions

**Admin Permissions** (`config/security/rolepermissions/Admin.csv`):

```csv
p,Admin,*,*
```

**Teacher Permissions** (`config/security/rolepermissions/Teacher.csv`):

```csv
g,Teacher,Student
p,Teacher,*,courses.create
p,Teacher,*,courses.update
p,Teacher,*,courses.delete
p,Teacher,*,grades.create
p,Teacher,*,grades.update
p,Teacher,*,students.read
```

**Student Permissions** (`config/security/rolepermissions/Student.csv`):

```csv
p,Student,*,profile.read
p,Student,*,profile.update.own
p,Student,*,courses.read
p,Student,*,grades.read.own
p,Student,*,enrollment.create
p,Student,*,enrollment.read.own
```

## Plugin Permissions

Expose permissions from your plugin configuration.

**Plugin Configuration** (`config/config.yml`):

```yaml
version: 1.0.0
description: Student Management Plugin
permissions:
  students.create: "Create Student"
  students.read: "Read Student"
  students.update: "Update Student"
  students.delete: "Delete Student"
  grades.create: "Create Grade"
  grades.read: "Read Grade"
  grades.update: "Update Grade"
  courses.create: "Create Course"
  courses.read: "Read Course"
  courses.update: "Update Course"
  enrollment.create: "Create Enrollment"
  enrollment.read: "Read Enrollment"
```

## Channel Configuration with Permissions

**Grades Channel** (`config/server/channels/studentgrades.yml`):

```yaml
id: student_grades
path: /api/student/:studentId/grades
service: student_grades_service
permission: grades.read.own
method: GET
```

**Course Management Channel** (`config/server/channels/coursemgmt.yml`):

```yaml
id: course_management
path: /api/courses
service: course_service
permission: courses.update
method: POST
```

## Checking Permissions in Code

You can programmatically check permissions in your services.

```go
func (s *GradesService) Invoke(ctx core.RequestContext) error {
    secHandler := ctx.GetSecurityHandler()
    
    // Check if user has permission to update grades
    if !secHandler.HasPermission(ctx, "grades.update") {
        return errors.Forbidden(ctx, "Insufficient permissions")
    }
    
    // Check if user has a specific role
    if !secHandler.HasRole(ctx, "Teacher") {
        return errors.Forbidden(ctx, "Only teachers can update grades")
    }
    
    // Proceed with grade update
    return nil
}
```

## Multi-Tenant Security

For multi-tenant applications, scope permissions by tenant.

**Tenant-Specific Access** (`config/security/serviceaccess/Student.csv`):

```csv
p,Student,{{tenant_id}},grades.read.own
p,Student,{{tenant_id}},profile.read
```

**Global Admin** (`config/security/serviceaccess/SuperAdmin.csv`):

```csv
p,SuperAdmin,*,*
```

## Best Practices

### User Management

1. **Hash Passwords**: Always hash passwords in PreSave hook
2. **Email Verification**: Verify email addresses before activation
3. **Strong Passwords**: Enforce password complexity requirements
4. **Account Lockout**: Implement lockout after failed login attempts

### Permission Design

1. **Principle of Least Privilege**: Grant minimum required permissions
2. **Hierarchical Permissions**: Use dot notation (`resource.action`)
3. **Granular Controls**: Separate read/write/delete permissions
4. **Document Permissions**: Maintain clear permission documentation

### Role Configuration

1. **Clear Role Names**: Use descriptive role names
2. **Role Inheritance**: Use inheritance to avoid duplication
3. **Avoid Wildcards**: Minimize use of `*` except for admin roles
4. **Regular Audits**: Review and audit role configurations

### Security Checklist

- [ ] User entity with encrypted passwords
- [ ] Role and permission entities defined
- [ ] Service access policies configured
- [ ] Role permissions mapped
- [ ] Channel permissions assigned
- [ ] Plugin permissions exposed
- [ ] Multi-tenancy configured (if applicable)
- [ ] Security documentation updated

## Common Patterns

### Self-Service Pattern

Allow users to manage their own data:

```csv
# Student can only read own data
p,Student,*,profile.read.own
p,Student,*,grades.read.own
```

Implement in service:

```go
// Verify user is accessing their own data
if studentId !== ctx.GetUserId() {
    return errors.Forbidden(ctx, "Can only access own data")
}
```

### Hierarchical Roles

Use role inheritance:

```csv
# Teacher inherits all Student permissions
g,Teacher,Student

# Admin inherits all Teacher permissions
g,Admin,Teacher
```

### Resource-Based Permissions

Scope permissions by resource:

```csv
# Read all courses
p,Student,*,courses.read

# Update only own courses
p,Instructor,*,courses.update.own

# Update any course
p,Admin,*,courses.update
```

## See Also

- [Services Guide](services.md) - Implement permission checks in services
- [Entities Guide](entities.md) - Define User, Role, and Permission entities
- [Channels Guide](channels.md) - Configure channel-level permissions
