# Database Setup - MongoDB

This guide explains how to configure MongoDB as the database for the Student Management System.

## Prerequisites

- MongoDB installed and running locally or access to a MongoDB instance
- Laatoo CLI installed
- Student Management solution created

## Overview

Laatoo uses the `mongodatabase` plugin to connect to MongoDB. The configuration involves:
1. Creating a module configuration file
2. Setting up environment variables for connection
3. Configuring entities to use the database
4. Optionally setting up master data

## Step 1: Configure MongoDB Module

Create the MongoDB module configuration file:

```bash
# Create the modules directory if it doesn't exist
mkdir -p config/modules

# Create the mongodatabase configuration
```

**File**: `config/modules/mongodatabase.yml`

```yaml
plugin: mongodatabase
settings:
  factoryname: database
  mongodatabase: studentmanagement
```

**Configuration Explanation**:
- `plugin`: Must be `mongodatabase` to use the MongoDB plugin
- `factoryname`: Name of the database factory (used by other modules to reference this connection)
- `mongodatabase`: Name of the MongoDB database to create/use

## Step 2: Set MongoDB Connection String

The connection string is provided via environment variables for security.

### Development Environment

Create or update `.env` file in your solution root:

```bash
# .env
mongoconnectionstring=mongodb://localhost:27017
```

### Production Environment

For production, use secrets manager:

```bash
# Add connection string to secrets
laatoo security addsecret mongoconnectionstring "mongodb://your-production-host:27017"
```

Or set as environment variable in your deployment:

```bash
export mongoconnectionstring="mongodb://your-production-host:27017"
```

### Docker Compose Setup

If using Docker Compose for development:

```yaml
# docker-compose.yml
version: '3.8'

services:
  mongodb:
    image: mongo:6.0
    container_name: studentmgmt-mongodb
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
    environment:
      MONGO_INITDB_DATABASE: studentmanagement

  studentmgmt:
    build: .
    ports:
      - "8080:8080"
    environment:
      - mongoconnectionstring=mongodb://mongodb:27017
    depends_on:
      - mongodb

volumes:
  mongodb_data:
```

## Step 3: Configure Entities for Database

Each entity in your plugins needs to specify which database factory to use.

### Update Entity Module Configurations

For each entity (Student, Exam, ExamAttempt, Result), create or update the module configuration:

**File**: `dev/plugins/studentmgmt/config/server/modules/Student.yml`

```yaml
plugin: entity
settings:
  object: studentmgmt.Student
  dataservicefactory: database  # References the 'factoryname' from mongodatabase.yml
```

Repeat for all entities:
- `Exam.yml`
- `ExamAttempt.yml`
- `Result.yml`

## Step 4: Master Data Setup (Optional)

To pre-populate your database with initial data:

### Configure Master Data Manager

**File**: `config/modules/masterdatamanager.yml`

```yaml
plugin: masterdatamanager
settings:
  dataservicefactory: database  # Must match mongodatabase factoryname
```

### Create Master Data Files

Create a `data` directory in your plugin and add JSON files for initial data:

```bash
mkdir -p dev/plugins/studentmgmt/data
```

**File**: `dev/plugins/studentmgmt/data/Student.json`

```json
[
  {
    "Id": "student-001",
    "Name": "Alice Johnson",
    "Email": "alice.johnson@example.com",
    "Grade": 10,
    "DateOfBirth": "2008-05-15T00:00:00Z",
    "Status": "active"
  },
  {
    "Id": "student-002",
    "Name": "Bob Smith",
    "Email": "bob.smith@example.com",
    "Grade": 11,
    "DateOfBirth": "2007-08-22T00:00:00Z",
    "Status": "active"
  }
]
```

**File**: `dev/plugins/studentmgmt/data/Exam.json`

```json
[
  {
    "Id": "exam-math-001",
    "Name": "Mathematics Midterm",
    "Subject": "Mathematics",
    "MaxScore": 100,
    "PassingScore": 60,
    "StartDate": "2024-06-01T09:00:00Z",
    "EndDate": "2024-06-01T11:00:00Z"
  }
]
```

## Step 5: Verify Configuration

### Check Configuration Files

Ensure you have:
- `config/modules/mongodatabase.yml` - Database connection
- `config/modules/masterdatamanager.yml` - (Optional) Master data
- Entity modules with `dataservicefactory: database`

### Test Connection

Start the solution and check logs:

```bash
# Build and run
laatoo build
laatoo run

# Check logs for MongoDB connection
# Look for: "Connected to MongoDB: studentmanagement"
```

### Verify Database

Connect to MongoDB and verify:

```bash
# Using MongoDB shell
mongosh

use studentmanagement
show collections

# Should see collections for your entities:
# - Student
# - Exam
# - ExamAttempt  
# - Result

# Check sample data
db.Student.find().pretty()
```

## Connection String Formats

### Local Development
```
mongodb://localhost:27017
```

### With Authentication
```
mongodb://username:password@localhost:27017
```

### MongoDB Atlas (Cloud)
```
mongodb+srv://username:password@cluster.mongodb.net/studentmanagement
```

### Replica Set
```
mongodb://host1:27017,host2:27017,host3:27017/?replicaSet=myReplicaSet
```

## Environment Variables Reference

| Variable | Description | Example |
|----------|-------------|---------|
| `mongoconnectionstring` | MongoDB connection URL | `mongodb://localhost:27017` |
| `mongoconnectionstringkey` | Secret key for connection string | `mongoconnectionstringlocal` |

## Switching Between Databases

The `factoryname` abstraction allows easy switching between databases:

### Switch to SQL Database

1. Create `config/modules/sqldatabase.yml`:
```yaml
plugin: sqldatabase
settings:
  factoryname: database  # Same name!
  sqlvendor: postgres
```

2. Set SQL connection string:
```bash
export sqlconnectionstring="postgres://user:pass@localhost:5432/studentmanagement"
```

3. Remove `mongodatabase.yml`

All entity modules continue to work because they reference `database` factory name.

## Troubleshooting

### Connection Refused
```
Error: connect ECONNREFUSED 127.0.0.1:27017
```
**Solution**: Ensure MongoDB is running:
```bash
# Start MongoDB
sudo systemctl start mongod

# Or with Docker
docker run -d -p 27017:27017 --name mongodb mongo:6.0
```

### Authentication Failed
```
Error: Authentication failed
```
**Solution**: Check username/password in connection string or use correct auth database:
```
mongodb://user:pass@localhost:27017/studentmanagement?authSource=admin
```

### Database Not Created
```
Collection not found
```
**Solution**: MongoDB creates databases/collections on first write. Insert a document to trigger creation, or ensure entities are being saved during solution startup.

### Multiple Database Instances

If running multiple Laatoo solutions, use different database names:

```yaml
# Solution 1: config/modules/mongodatabase.yml
settings:
  mongodatabase: studentmanagement

# Solution 2: config/modules/mongodatabase.yml
settings:
  mongodatabase: anothersolution
```

## Best Practices

1. **Never Commit Connection Strings**: Use `.env` files (add to `.gitignore`)
2. **Use Secrets in Production**: Use `laatoo security addsecret` for sensitive data
3. **Consistent Factory Names**: Always use `database` as factoryname for portability
4. **Index Your Entities**: Add indexes in entity definitions for performance
5. **Backup Regularly**: Set up automated MongoDB backups

## Next Steps

- Define indexes in your entity configurations
- Set up database backups
- Configure security rules
- Set up monitoring and logging

## See Also

- [Entity Configuration](../server-development/entities.md)
- [Data Modeling](../server-development/data-modeling.md)
- [Master Data Management](../server-development/master-data.md)
- [Production Deployment](09-deployment.md)
