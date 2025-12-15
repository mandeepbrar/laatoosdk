# Laatoo Platform SDK

Laatoo is a modular, plugin-based application development platform that enables rapid development of scalable,multi-tenant enterprise applications. Built in Go with a React-based UI framework, Laatoo provides a complete ecosystem for building backend services, workflows, and modern web applications.

## What is Laatoo?

Laatoo is designed to solve the challenge of building complex, multi-tenant enterprise systems by providing:

- **Plugin-Based Architecture**: Build modular, reusable components that can be shared across projects
- **Multi-Tenancy by Default**: Built-in isolation and tenant management
- **Hierarchical Configuration**: Solution → Application → Isolation hierarchy for maximum flexibility
- **Comprehensive Server SDK**: Services, workflows, tasks, security, data access, and more
- **Modern UI Framework**: React-based UI with forms, views, blocks, pages, and actions
- **Developer-Friendly CLI**: Command-line tools for scaffolding, building, and deployment

## Key Features

### Server Development
- **Service-Oriented Architecture**: Create reusable services with clean interfaces
- **Entity Modeling**: Define data models with automatic CRUD generation
- **Workflow Engine**: Build complex business processes using YAML-based workflows
- **Background Tasks**: Queue and process asynchronous tasks
- **Real-Time Messaging**: Pub/sub messaging system
- **Flexible Security**: JWT-based authentication with role-based access control
- **Data Management**: Support for multiple databases (PostgreSQL, Firestore, MongoDB)

### UI Development
- **Component-Based**: Build UIs using reusable blocks and components
- **Forms**: Declarative form definitions with validation
- **Views**: Flexible data views with filtering and pagination
- **Pages**: Route-based page system
- **Actions**: Declarative actions for complex interactions
- **Datasets**: Powerful data querying with OData support
- **Design System**: Pre-built components following modern design principles

### Multi-Tenancy
- **Isolation**: Complete data and configuration isolation per tenant
- **Hierarchical Scoping**: Modules can be scoped to solution, application, or tenant level
- **Flexible Deployment**: Run multiple tenants on shared or dedicated infrastructure

### DevOps
- **CLI Tools**: Complete command-line interface for all operations
- **Hot Reloading**: Develop with fast feedback loops
- **Docker Support**: Containerized deployment out of the box
- **Environment Management**: Separate configurations for dev, staging, production

## Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│                    Solution                         │
│  ┌─────────────────────────────────────────────┐   │
│  │            Application                       │   │
│  │  ┌─────────────────────────────────────┐   │   │
│  │  │         Isolation (Tenant)          │   │   │
│  │  │  ┌──────────────┬────────────────┐ │   │   │
│  │  │  │   Server     │   UI Plugins   │ │   │   │
│  │  │  │   Plugins    │                │ │   │   │
│  │  │  │              │                │ │   │   │
│  │  │  │  • Services  │   • Pages      │ │   │   │
│  │  │  │  • Entities  │   • Forms      │ │   │   │
│  │  │  │  • Workflows │   • Views      │ │   │   │
│  │  │  │  • Tasks     │   • Blocks     │ │   │   │
│  │  │  │  • Security  │   • Actions    │ │   │   │
│  │  │  └──────────────┴────────────────┘ │   │   │
│  │  └─────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
```

## Quick Start

### Prerequisites
- Go 1.19+ (for server development)
- Node.js 16+ (for UI development)
- Docker (for running solutions locally)

### Installation

```bash
# Install Laatoo CLI
go install laatoo.io/cli/laatoo@latest

# Verify installation
laatoo version
```

### Create Your First Solution

```bash
# Create a new solution
laatoo solution create my-first-app

# Navigate to solution directory
cd my-first-app

# Create an application
cd applications
mkdir myapp
cd myapp

# Create a server plugin
laatoo plugin create studentplugin -t server

# Build the plugin
laatoo plugin build studentplugin

# Run the solution
cd ../..
laatoo solution run my-first-app
```

## Learning Path

### 1. Getting Started
Start here if you're new to Laatoo:
- [Introduction](docs/getting-started/01-introduction.md) - Understand Laatoo's architecture
- [Setup](docs/getting-started/02-setup.md) - Install and configure
- [First Solution](docs/getting-started/03-first-solution.md) - Create your first app

### 2. Complete Tutorial
Follow the Student Management System tutorial to build a real application:
- [Overview](docs/tutorials/student-management/00-overview.md)
- [Setup Solution](docs/tutorials/student-management/01-setup-solution.md)
- [Data Model](docs/tutorials/student-management/02-data-model.md)
- [Server Plugin](docs/tutorials/student-management/03-server-plugin.md)
- [Workflows](docs/tutorials/student-management/04-workflows.md)
- [UI Plugin](docs/tutorials/student-management/05-ui-plugin.md)
- [Forms](docs/tutorials/student-management/06-forms.md)
- [Views & Lists](docs/tutorials/student-management/07-views-lists.md)
- [Dashboards](docs/tutorials/student-management/08-dashboards.md)
- [Deployment](docs/tutorials/student-management/09-deployment.md)

### 3. Development Guides
Deep dive into specific topics:

**Server Development:**
- [Creating Plugins](docs/guides/server-development/creating-plugins.md)
- [Entities](docs/guides/server-development/entities.md)
- [Services](docs/guides/server-development/services.md)
- [Channels](docs/guides/server-development/channels.md)
- [Workflows](docs/guides/server-development/workflows.md)
- [Tasks](docs/guides/server-development/tasks.md)
- [Security](docs/guides/server-development/security.md)

**UI Development:**
- [Creating UI Plugins](docs/guides/ui-development/creating-ui-plugins.md)
- [Pages](docs/guides/ui-development/pages.md)
- [Forms](docs/guides/ui-development/forms.md)
- [Blocks](docs/guides/ui-development/blocks.md)
- [Views](docs/guides/ui-development/views.md)
- [Datasets](docs/guides/ui-development/datasets.md)
- [Actions](docs/guides/ui-development/actions.md)
- [Services](docs/guides/ui-development/services.md)

### 4. Reference Documentation
- [CLI Commands](docs/reference/cli-commands.md)
- [SDK Interfaces](docs/reference/sdk-interfaces.md)
- [Configuration](docs/reference/configuration.md)

## Example: Student Management System

Throughout the documentation, we use a Student Management System as a practical example. This system includes:

**Features:**
- Student registration and profiles
- Exam creation and management
- Exam submission and grading
- Results dashboard
- Role-based access (Students, Teachers, Admins)

**Technical Components:**
- **Server Plugin**: Entities (Student, Exam, Result), Services, Workflows
- **UI Plugin**: Registration forms, student lists, exam views, results dashboard
- **Security**: Multi-role authentication and authorization
- **Workflows**: Automated exam processing and result generation

## Use Cases

Laatoo is ideal for:

- **SaaS Applications**: Multi-tenant systems with complex data models
- **Enterprise Portals**: Internal applications with multiple user roles
- **Business Process Management**: Workflow-driven applications
- **Data Management Systems**: Applications requiring flexible data models
- **Microservices**: Building modular, service-oriented architectures

## Next Steps

Ready to start building? Head to the [Introduction](docs/getting-started/01-introduction.md) guide to learn more about Laatoo's architecture, or jump straight to creating your [First Solution](docs/getting-started/03-first-solution.md).

---
