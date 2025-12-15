# Chapter 5: UI Plugin - Setup

Create the user interface plugin for Student Management System.

## Overview

The UI plugin provides:
- Student-facing pages (exam list, take exam, view results)
- Teacher-facing pages (manage students, create exams, view all results)
- Admin dashboard

In this chapter, we'll set up the UI plugin structure and create the main navigation.

## Step 1: Create UI Plugin

```bash
cd dev/plugins
laatoo plugin create studentmgmt-ui -t ui
cd studentmgmt-ui
```

This creates:
```
studentmgmt-ui/
├── config/
│   ├── config.yml
│   └── ui/
│       └── build.yml       # Build configuration
└── src/
    └── ui/
        ├── js/             # JavaScript code
        │   └── index.js    # Initialize method
        ├── styles/         # CSS/SCSS styles
        └── registry/       # UI definitions
            ├── Pages/      # Page definitions
            ├── Forms/      # Form definitions
            ├── Views/      # View definitions
            ├── Blocks/     # UI blocks
            ├── Actions/    # Actions
            └── Menus/      # Menu definitions
```

## Step 2: Configure Build Dependencies

Edit `config/ui/build.yml`:

```yaml
js:
  externals:
    - react
    - react-dom
    - react-redux
    - redux-saga
    - reactpages           # Laatoo UI framework
    - reactwebcommon       # Common components
    - jsui                 # Core UI utilities
    - core-js
    
  packages:
    react: "^18.1.0"
    react-dom: "^18.1.0"
    react-redux: "^7.1.1"
    redux-saga: "^1.1.0"
    core-js: "^3.2.1"
    
css:
  includePaths:
    - src/ui/styles
```

## Step 3: Create Main Menu

Create `src/ui/registry/Menus/main-menu.yml`:

```yaml
# Main navigation menu
vertical: true
items:
  - title: Dashboard
    page: /dashboard
    iconClass: fa fa-home
    roles: [Student, Teacher, Admin]
    
  - title: Students
    page: /students
    iconClass: fa fa-users
    roles: [Teacher, Admin]
    
  - title: Exams
    page: /exams
    iconClass: fa fa-file-text
    roles: [Teacher, Admin]
    
  - title: My Exams
    page: /my-exams
    iconClass: fa fa-pencil
    roles: [Student]
    
  - title: Results
    page: /results
    iconClass: fa fa-chart-bar
    roles: [Teacher, Admin]
    
  - title: My Results
    page: /my-results
    iconClass: fa fa-trophy
    roles: [Student]
```

## Step 4: Create Page Layout

Create `src/ui/registry/Pages/dashboard.yml`:

```yaml
route: "/dashboard"
authenticate: true
component:
  type: layout
  layout: 2col
  leftcol:
    type: menu
    id: main-menu
    className: sidebar-menu
  rightcol:
    type: block
    id: dashboard-content
    className: main-content
```

## Step 5: Create Dashboard Block

Create `src/ui/registry/Blocks/dashboard-content.xml`:

```xml
<Block className="dashboard-container">
  <h1 module="html">Student Management Dashboard</h1>
  
  <!-- Statistics Cards -->
  <Block className="stats-grid">
    <Block className="stat-card">
      <h3 module="html">Total Students</h3>
      <Block id="totalStudents" className="stat-value">0</Block>
    </Block>
    
    <Block className="stat-card">
      <h3 module="html">Active Exams</h3>
      <Block id="activeExams" className="stat-value">0</Block>
    </Block>
    
    <Block className="stat-card">
      <h3 module="html">Completed Exams</h3>
      <Block id="completedExams" className="stat-value">0</Block>
    </Block>
    
    <Block className="stat-card">
      <h3 module="html">Average Score</h3>
      <Block id="avgScore" className="stat-value">0%</Block>
    </Block>
  </Block>
  
  <!-- Recent Activity -->
  <Block className="recent-activity">
    <h2 module="html">Recent Activity</h2>
    <Block id="activityList"></Block>
  </Block>
</Block>
```

## Step 6: Create Initialize Script

Create `src/ui/js/index.js`:

```javascript
// Initialize method - called by reactapplication when plugin loads
var module;

function Initialize(appName, ins, mod, settings, def, req) {
  module = this;
  
  // Access to properties and localization
  module.properties = Application.Properties[ins];
  module.settings = settings;
  
  console.log('Student Management UI plugin initialized');
}

export {
  Initialize
}
```

> **Note**: The Initialize method is the entry point for your UI plugin. It's called when the plugin is loaded by the `reactapplication` module. See [Initialize Method Guide](../../guides/ui-development/initialize.md) for complete documentation.

## Step 7: Create Styles

Create `src/ui/styles/main.scss`:

```scss
// Dashboard styles
.dashboard-container {
  padding: 20px;
  
  h1 {
    color: #333;
    margin-bottom: 30px;
  }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.stat-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  
  h3 {
    margin: 0 0 10px 0;
    font-size: 14px;
    opacity: 0.9;
  }
  
  .stat-value {
    font-size: 32px;
    font-weight: bold;
  }
}

.recent-activity {
  h2 {
    color: #333;
    margin-bottom: 20px;
  }
}

.sidebar-menu {
  background: #2c3e50;
  height: 100vh;
  padding: 20px 0;
  
  .menu-item {
    color: white;
    padding: 15px 20px;
    cursor: pointer;
    transition: background 0.3s;
    
    &:hover {
      background: #34495e;
    }
    
    &.active {
      background: #3498db;
    }
    
    i {
      margin-right: 10px;
    }
  }
}
```

## Step 8: Build the UI Plugin

```bash
# Install npm dependencies
laatoo plugin build studentmgmt-ui --getbuildpackages

# This will:
# 1. Install npm packages
# 2. Bundle JavaScript
# 3. Compile SCSS to CSS
# 4. Create bin/ directory
```

## Step 9: Install UI Plugin

```bash
# Copy to application
cp -r bin ../../applications/education/config/modules/studentmgmt-ui
```

## Step 10: Test the UI

```bash
# Restart solution
cd ../../../..
laatoo solution run student-mgmt-system
```

Open browser to `http://localhost:8080/dashboard`

You should see:
- Sidebar menu with navigation items
- Dashboard with statistics cards
- Recent activity section

## Understanding UI Structure

### Pages

Pages are route-based and define the overall layout:

```yaml
route: "/dashboard"        # URL path
authenticate: true         # Requires login
component:
  type: layout
  layout: 2col            # Two-column layout
  leftcol:                # Left column content
    type: menu
    id: main-menu
  rightcol:               # Right column content
    type: block
    id: dashboard-content
```

### Blocks

Blocks are reusable UI components:

```xml
<Block className="my-block">
  <h1 module="html">Title</h1>
  <Block id="content"></Block>
</Block>
```

### Initialize Script

The Initialize method is called when the plugin loads. See the [Initialize Method Guide](../../guides/ui-development/initialize.md) for complete documentation on:
- Method signature and parameters
- Accessing properties and settings
- Using global Window functions

## Best Practices

### 1. Use Semantic Class Names

```xml
<!-- Good -->
<Block className="student-card">
<Block className="exam-list">

<!-- Avoid -->
<Block className="blk1">
<Block className="container2">
```

### 2. Separate Concerns

- **Pages**: Define routes and layout
- **Blocks**: UI components
- **Forms**: Data entry
- **Views**: Data display
- **Actions**: User interactions
- **Services**: Business logic

### 3. Handle Errors

```javascript
try {
  const data = await context.callService('Service', params);
  // Handle success
} catch (error) {
  console.error('Error:', error);
  // Show error message to user
}
```

### 4. Use CSS Variables

```scss
:root {
  --primary-color: #667eea;
  --secondary-color: #764ba2;
  --text-color: #333;
}

.stat-card {
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
}
```

## Summary

You've created:
- ✅ UI plugin structure
- ✅ Main navigation menu
- ✅ Dashboard page with layout
- ✅ Dashboard block with statistics
- ✅ Initialize script with data loading
- ✅ Styles for UI components

**Next**: [Chapter 6: Forms](06-forms.md) - Create forms for student enrollment and exam creation
