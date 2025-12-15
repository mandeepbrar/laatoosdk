# Chapter 1: Setup Solution

Create the foundation for your Student Management System.

## Create Solution Structure

\`\`\`bash
# Create solution
laatoo solution create student-mgmt-system
cd student-mgmt-system

# Create application
mkdir -p applications/education
mkdir -p applications/education/config
mkdir -p applications/education/isolations/school-a/config

# Create server configuration
mkdir -p config/modules
mkdir -p config/engines
```

## Configure Solution

Create `config/config.yml`:
```yaml
name: student-mgmt-system
logging:
  level: debug
```

Create `config/engines/http.yml`:
```yaml
name: root
enginetype: http
address: ":8080"
path: "/"
```

Create `config/modules/user.yml`:
```yaml
plugin: user
```

## Configure Application

Create `applications/education/config/config.yml`:
```yaml
name: education
description: Educational management application
\`\`\`

## Configure Isolation

Create `applications/education/isolations/school-a/config/config.yml`:
\`\`\`yaml
name: school-a
description: School A tenant

database:
  type: postgres
  host: localhost
  port: 5432
  database: school_a
  user: laatoo
  password: laatoo
\`\`\`

## Create Plugin Structure

\`\`\`bash
cd dev/plugins
laatoo plugin create studentmgmt -t full
cd studentmgmt
\`\`\`

**Next**: [Chapter 2: Data Model](02-data-model.md)
