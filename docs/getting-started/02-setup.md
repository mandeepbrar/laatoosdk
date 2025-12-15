# Setup and Installation

This guide walks you through installing Laatoo and setting up your development environment.

## Prerequisites

### Required Software

1. **Go 1.19 or later**
   - Download from [golang.org](https://golang.org/dl/)
   - Verify installation: `go version`

2. **Node.js 16 or later** (for UI development)
   - Download from [nodejs.org](https://nodejs.org/)
   - Verify installation: `node --version` and `npm --version`

3. **Docker** (recommended for local development)
   - Download from [docker.com](https://www.docker.com/get-started)
   - Verify installation: `docker --version`

###Optional but Recommended

- **Git**: For version control
- **VSCode** or **GoLand**: IDEs with good Go support
- **PostgreSQL**: If using local database (instead of Docker)

## Installing Laatoo CLI

The Laatoo CLI (`laatoo` command) is the primary tool for creating, building, and running Laatoo solutions.


### Configuration

The Laatoo CLI uses a configuration file at `~/.laatoo/config.yaml`:

```yaml
# Default repository for plugins
repository: https://plugins.laatoo.io

# Default build mode
buildmode: local

# Development settings
dev:
  hotload: true
  verbose: false
```

Create this file if it doesn't exist:

```bash
mkdir -p ~/.laatoo
cat > ~/.laatoo/config.yaml << EOF
repository: ""
buildmode: local
dev:
  hotload: true
  verbose: false
EOF
```

## Setting Up Development Environment

### Directory Structure

Create a workspace for your Laatoo projects:

```bash
# Create workspace directory
mkdir -p ~/laatoo-workspace
cd ~/laatoo-workspace

# Create standard directories
mkdir -p {solutions,plugins,libs}
```

**Directory Purpose**:
- `solutions/`: Your Laatoo solutions
- `dev/plugins/`: Plugin development
- `libs/`: Shared libraries

### Environment Variables

Add these to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
# Laatoo workspace
export LAATOO_WORKSPACE=~/laatoo-workspace

# Add Laatoo to PATH (if installed from source)
export PATH=$PATH:$GOPATH/bin

# Optional: Default solution
export LAATOO_DEFAULT_SOLUTION=my-solution
```

Reload your shell:
```bash
source ~/.bashrc  # or ~/.zshrc
```

## Database Setup

Laatoo supports multiple databases. For local development, we recommend PostgreSQL.

### Using Docker (Recommended)

```bash
# Start PostgreSQL in Docker
docker run --name laatoo-postgres \
  -e POSTGRES_PASSWORD=laatoo \
  -e POSTGRES_USER=laatoo \
  -e POSTGRES_DB=laatoo \
  -p 5432:5432 \
  -d postgres:14

# Verify it's running
docker ps | grep laatoo-postgres
```

### Using Local PostgreSQL

If you prefer a local installation:

```bash
# Ubuntu/Debian
sudo apt-get install postgresql postgresql-contrib

# macOS
brew install postgresql

# Start PostgreSQL
sudo systemctl start postgresql  # Linux
brew services start postgresql   # macOS

# Create database and user
sudo -u postgres psql
CREATE DATABASE laatoo;
CREATE USER laatoo WITH PASSWORD 'laatoo';
GRANT ALL PRIVILEGES ON DATABASE laatoo TO laatoo;
\q
```

## Verifying Installation

### Check All Tools

Run these commands to verify everything is installed:

```bash
# Go
go version
# Output: go version go1.20.x ...

# Node.js
node --version
# Output: v18.x.x or later

# NPM
npm --version
# Output: 9.x.x or later

# Docker
docker --version
# Output: Docker version 20.x.x ...

# Laatoo CLI
laatoo version
# Output: Laatoo CLI version x.x.x

# PostgreSQL (if using Docker)
docker exec laatoo-postgres psql -U laatoo -c "SELECT version();"
```

### Test Laatoo CLI

```bash
# List available commands
laatoo help

# Check plugin commands
laatoo plugin help

# Check solution commands
laatoo solution help
```

Expected output should show all available commands and flags.

## IDE Setup

### Visual Studio Code

Recommended extensions:
- **Go** (golang.go)
- **React/JSX** (if doing UI development)
- **YAML** (redhat.vscode-yaml)
- **Docker** (ms-azuretools.vscode-docker)

Settings for Laatoo development (`.vscode/settings.json`):

```json
{
  "go.useLanguageServer": true,
  "go.gopath": "${env:GOPATH}",
  "go.toolsManagement.autoUpdate": true,
  "editor.formatOnSave": true,
  "[go]": {
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

### GoLand

1. Open **Settings** → **Go** → **GOPATH**
2. Add your workspace directory
3. Enable **Go Modules** support
4. Configure **File Watchers** for auto-formatting

## Troubleshooting

### Common Issues

**Issue**: `laatoo: command not found`

**Solution**: Ensure `$GOPATH/bin` is in your PATH:
```bash
echo $PATH | grep go
export PATH=$PATH:$(go env GOPATH)/bin
```

---

**Issue**: `Cannot connect to Docker daemon`

**Solution**: Start Docker service:
```bash
# Linux
sudo systemctl start docker

# macOS
open -a Docker
```

---

**Issue**: Go module errors when building plugins

**Solution**: Initialize Go modules in your plugin:
```bash
cd dev/plugins/myplugin/src/server/go
go mod init myplugin
go mod tidy
```

---

**Issue**: PostgreSQL connection refused

**Solution**: Check if PostgreSQL is running:
```bash
# If using Docker
docker ps | grep postgres

# If using local installation
sudo systemctl status postgresql  # Linux
brew services list | grep postgresql  # macOS
```

## Next Steps

Now that your environment is set up:

1. **[Create Your First Solution](03-first-solution.md)**: Build a simple Laatoo application
2. **[Student Management Tutorial](../tutorials/student-management/00-overview.md)**: Follow a complete example
3. **[Server Plugin Guide](../guides/server-development/creating-plugins.md)**: Deep dive into plugin development

## Additional Resources

- **Laatoo Documentation**: [Full docs](../)
- **Go Documentation**: [golang.org/doc](https://golang.org/doc)
- **React Documentation**: [reactjs.org](https://reactjs.org)
- **Docker Documentation**: [docs.docker.com](https://docs.docker.com)
