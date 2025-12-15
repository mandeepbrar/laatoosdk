# Running Laatoo Solutions

This guide explains how to run the Student Management System and other Laatoo solutions in various modes.

## Prerequisites

- Laatoo CLI installed
- MongoDB running (local or Docker)
- Node.js 18+ (for local UI development)
- Docker and Docker Compose (for containerized deployment)

## Quick Start

```bash
# Navigate to your solution
cd examples/studentmanagement

# Copy environment template
cp .env.example .env

# Start MongoDB (if not running)
docker run -d -p 27017:27017 --name mongodb mongo:6.0

# Run the solution
laatoo solution run --runmode local
```

The solution will be available at `http://localhost:8080`

## Run Modes

Laatoo supports three primary run modes for different development and deployment scenarios.

### Local Mode (Development)

**Best for**: Day-to-day development with fast iteration

```bash
laatoo solution run --runmode local
```

**Features**:
- Runs directly on your machine
- Fast startup (~5-10 seconds)
- Hot reload for code changes
- Direct database connection
- Easy debugging

**Environment File**: `environment_local.yml`

**Example for Student Management**:
```bash
# Start with local MongoDB
laatoo solution run --runmode local

# With debugger attached
laatoo solution run --runmode local --withlaatoodebugger
```

### Docker Debug Mode (Testing)

**Best for**: Testing in production-like environment

```bash
laatoo solution run --runmode dockerdebug
```

**Features**:
- Containerized application
- Simulates production setup
- Service isolation
- Network simulation

**Environment File**: `environment_dockerdebug.yml`

**Example**:
```bash
# Run in Docker containers
laatoo solution run --runmode dockerdebug

# With log routing
laatoo solution run --runmode dockerdebug --routealllogs

# Clean up after testing
laatoo solution run --runmode dockerdebug --shutdowncontainers
```

### Production Mode

**Best for**: Deploying to production

```bash
laatoo solution run --runmode prod
```

**Environment File**: `environment_prod.yml` (deployment-specific)

## Environment Configuration

### Environment Files

Laatoo uses YAML files to configure runtime behavior:

- `environment_local.yml` - Local development
- `environment_dockerdebug.yml` - Docker testing
- `environment_prod.yml` - Production (optional)

### Student Management: environment_local.yml

```yaml
mongoconnectionstring: mongodb://localhost:27017
production: false
sslenabled: false
registrationallowed: true
port: 8080
```

### Environment Variables with LAATOO_ Prefix

Any environment variable starting with `LAATOO_` is automatically available to your solution.

**Pattern**:
```bash
LAATOO_<variablename>=<value>
```

**Example**:
```bash
# Override MongoDB connection
export LAATOO_mongoconnectionstring=mongodb://localhost:27018

# Set production mode
export LAATOO_production=true

# Custom variables
export LAATOO_maxuploadsize=10485760

# Run the solution
laatoo solution run --runmode local
```

### Common Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `LAATOO_mongoconnectionstring` | MongoDB connection URL | `mongodb://localhost:27017` |
| `LAATOO_production` | Production mode flag | `true` or `false` |
| `LAATOO_port` | HTTP server port | `8080` |
| `LAATOO_redishost` | Redis server | `localhost:6379` |
| `LAATOO_sslenabled` | Enable SSL/TLS | `true` |
| `LAATOO_seedpass` | Decrypt secrets | `your-seed-password` |

### Priority Order

Configuration is loaded in this order (later overrides earlier):

1. Default `config.yml`
2. Environment-specific file (`environment_local.yml`)
3. **`LAATOO_` environment variables** (highest priority)

## Docker Setup

### Using Docker Compose

Create `docker-compose.yml` for Student Management:

```yaml
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
    container_name: studentmgmt-app
    ports:
      - "8080:8080"
    environment:
      - LAATOO_mongoconnectionstring=mongodb://mongodb:27017
      - LAATOO_production=false
    depends_on:
      - mongodb
    volumes:
      - ./dev:/app/dev
      - ./config:/app/config

volumes:
  mongodb_data:
```

**Start Services**:
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Remove volumes (clean database)
docker-compose down -v
```

## Development Workflow

### 1. First Time Setup

```bash
# Clone and setup
git clone <repo>
cd studentmanagement

# Copy environment file
cp .env.example .env

# Install dependencies (if needed)
npm install

# Start MongoDB
docker run -d -p 27017:27017 --name mongodb mongo:6.0

# Run solution
laatoo solution run --runmode local
```

### 2. Daily Development

```bash
# Start development server
laatoo solution run --runmode local

# Make changes to code
# Server auto-reloads on file changes

# Test endpoints
curl http://localhost:8080/api/students
```

### 3. Testing Before Deployment

```bash
# Test in Docker
laatoo solution run --runmode dockerdebug

# Run tests
npm test

# Check logs
docker-compose logs -f

# Clean up
laatoo solution run --runmode dockerdebug --shutdowncontainers
```

## Debugging

### Local Development Debugging

**VS Code Configuration** (`.vscode/launch.json`):

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Laatoo",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}",
      "args": ["solution", "run", "--runmode", "local", "--withlaatoodebugger"],
      "envFile": "${workspaceFolder}/.env"
    }
  ]
}
```

**Run with Debugger**:
```bash
laatoo solution run --runmode local --withlaatoodebugger
```

Then attach your IDE debugger to the process.

### Docker Debugging

```bash
# Start with debug logging
laatoo solution run --runmode dockerdebug --routealllogs

# View specific service logs
docker-compose logs -f studentmgmt

# Execute commands in container
docker exec -it studentmgmt-app sh

# Check MongoDB connection
docker exec -it studentmgmt-mongodb mongosh
```

### Common Debug Scenarios

**Check if server is running**:
```bash
curl http://localhost:8080/health
```

**View environment variables**:
```bash
laatoo solution run --runmode local --withlaatoodebugger
# Check logs for configuration loaded
```

**Test database connection**:
```bash
# MongoDB shell
mongosh mongodb://localhost:27017/studentmanagement

# List collections
show collections

# Query students
db.Student.find().pretty()
```

## Common Issues

### Issue: Port Already in Use

**Error**: `Address already in use ::8080`

**Solution**:
```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>

# Or use different port
export LAATOO_port=8081
laatoo solution run --runmode local
```

### Issue: MongoDB Connection Failed

**Error**: `Failed to connect to MongoDB`

**Solution**:
```bash
# Check if MongoDB is running
docker ps | grep mongo

# Start MongoDB
docker run -d -p 27017:27017 --name mongodb mongo:6.0

# Verify connection string
echo $LAATOO_mongoconnectionstring

# Test connection
mongosh mongodb://localhost:27017
```

### Issue: Hot Reload Not Working

**Solution**:
```bash
# Check devhome flag
laatoo solution run --runmode local --devhome /path/to/project

# Verify file permissions
ls -la dev/plugins/
```

### Issue: Containers Won't Start

**Solution**:
```bash
# Check Docker daemon
docker ps

# View error logs
docker-compose logs

# Restart with clean state
docker-compose down -v
laatoo solution run --runmode dockerdebug --removecontainers
```

### Issue: Database Not Populated

**Error**: No data in collections

**Solution**:
```bash
# Check master data files
ls -la dev/plugins/studentmgmt/data/

# Verify masterdatamanager configuration
cat config/modules/masterdatamanager.yml

# Restart to reload data
laatoo solution run --runmode local
```

## Command Reference

### Run Commands

```bash
# Basic local run
laatoo solution run

# Specific run mode
laatoo solution run --runmode local|dockerdebug|prod

# With debugger
laatoo solution run --runmode local --withlaatoodebugger

# Route all logs
laatoo solution run --runmode dockerdebug --routealllogs

# Clean up containers
laatoo solution run --runmode dockerdebug --removecontainers

# Shutdown containers
laatoo solution run --shutdowncontainers

# Synchronous startup
laatoo solution run --sync

# Override environment
laatoo solution run --overrideenv
```

### Docker Commands

```bash
# View containers
docker ps

# View logs
docker-compose logs -f

# View specific service
docker-compose logs -f studentmgmt

# Stop services
docker-compose down

# Remove volumes
docker-compose down -v

# Rebuild images
docker-compose build --no-cache

# Execute in container
docker exec -it studentmgmt-app sh
```

## Best Practices

1. **Use Local Mode for Development**: Fastest iteration cycle
2. **Test in Docker Before Deployment**: Catch environment issues early
3. **Use Environment Files**: Don't hardcode configuration
4. **Use `LAATOO_` Variables for Secrets**: Never commit sensitive data
5. **Route Logs in Docker**: Use `--routealllogs` for debugging
6. **Clean Up Resources**: Use `--removecontainers` after testing
7. **Version Control Environment Templates**: Commit `.env.example`, not `.env`
8. **Document Required Variables**: List in README

## Student Management Specific

### Default Configuration

- **Port**: 8080
- **MongoDB**: localhost:27017
- **Database**: studentmanagement
- **Sample Data**: 3 students, 2 exams

### Accessing the Application

```bash
# API Endpoints
curl http://localhost:8080/studentmgmt.Student
curl http://localhost:8080/studentmgmt.Exam

# UI (if configured)
http://localhost:8080/students
```

### Verifying Setup

```bash
# 1. Check server is running
curl http://localhost:8080/health

# 2. Check database connection
mongosh mongodb://localhost:27017/studentmanagement

# 3. Verify sample data loaded
db.Student.find().pretty()

# 4. Test API
curl http://localhost:8080/studentmgmt.Student
```

## Next Steps

- [Database Setup](../tutorials/student-management/10-database-setup.md)
- [Security Configuration](security-detailed.md)
- [Testing Guide](server-development/testing.md)
- [Deployment](../tutorials/student-management/09-deployment.md)

## See Also

- [Creating Solutions](creating-solutions.md)
- [Environment Configuration](../reference/configuration.md)
- [CLI Commands](../reference/cli-commands.md)
