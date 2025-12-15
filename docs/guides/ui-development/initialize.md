# Initialize Method

The `Initialize` method is the entry point for Laatoo UI plugins. It is called by the `reactapplication` module when your plugin is loaded during application startup.

## Purpose

The `Initialize` method allows your UI plugin to:
- Register custom components and handlers
- Access plugin settings and properties
- Set up initial state and configuration
- Register boot-time methods
- Initialize module-level variables

## Method Signature

```javascript
function Initialize(appName, ins, mod, settings, def, req) {
  module = this;
  module.properties = Application.Properties[ins];
  module.settings = settings;
  
  // Your initialization code here
}
```

## Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `appName` | string | Application name |
| `ins` | string | Instance/module name (plugin name) |
| `mod` | object | Module reference |
| `settings` | object | Plugin settings from configuration |
| `def` | object | Default configuration |
| `req` | function | Require function for dynamic imports |

## The `this` Context

The `this` context in `Initialize` represents the module being initialized. You should:

1. **Assign to a module variable**: `module = this`
2. **Access properties**: `module.properties = Application.Properties[ins]`
3. **Store settings**: `module.settings = settings`

## Properties and Localization

### Accessing Properties

Properties are localized strings and configuration values defined in your plugin's properties files.

```javascript
function Initialize(appName, ins, mod, settings, def, req) {
  module = this;
  // Access properties for this plugin instance
  module.properties = Application.Properties[ins];
  
  // Now you can use properties in your code
  // e.g., module.properties.welcomeMessage
}
```

Properties file example (`properties/default.yml` in your plugin):
```yaml
welcomeMessage: "Welcome to Student Management"
buttonLabels:
  save: "Save Student"
  cancel: "Cancel"
  delete: "Delete"
```

### Using Properties in Components

Once properties are loaded, they can be accessed throughout your plugin:

```javascript
// In your UI components
const welcomeText = module.properties.welcomeMessage;
const saveLabel = module.properties.buttonLabels.save;
```

## Settings

The `settings` parameter contains runtime settings and configuration for your plugin. These come from:
- Plugin configuration files
- Application-level settings
- Runtime environment variables

```javascript
function Initialize(appName, ins, mod, settings, def, req) {
  module = this;
  module.settings = settings;
  
  // Access settings
  if (module.settings.enableFeatureX) {
    // Enable feature
  }
}
```

## Common Initialization Tasks

### 1. Registering Handlers

Use the `_r()` global function to register handlers and components:

```javascript
function Initialize(appName, ins, mod, settings, def, req) {
  // Register a boot method
  _r("Bootmethods", "loadStudentData", function(store, Application) {
    // Load initial data
  });
  
  // Register a custom handler
  _r("CustomHandlers", "myHandler", myHandlerFunction);
}
```

### 2. Setting Up Module State

```javascript
var module;
var studentCache = {};

function Initialize(appName, ins, mod, settings, def, req) {
  module = this;
  module.properties = Application.Properties[ins];
  module.settings = settings;
  
  // Initialize module-level state
  studentCache = {};
}
```

### 3. Accessing Global Window Functions

The Laatoo platform provides global functions on the `Window` object:

```javascript
// Execute a service
Window.executeService("serviceName", payload, urlParams, successCallback, errorCallback);

// Show a message
Window.showMessage({ Default: "Operation completed" });

// Show an error
Window.showError({ Default: "An error occurred" });

// Execute an action
Window.executeAction("actionName", payload, actionParams);

// Show/close dialogs
Window.showDialog(title, component, onClose, actions);
Window.closeDialog();
```

## Example: Student Management Initialize

```javascript
var module;

function Initialize(appName, ins, mod, settings, def, req) {
  module = this;
  module.properties = Application.Properties[ins];
  module.settings = settings;
  
  console.log('Student Management Plugin initialized');
  
  // Register a boot method to load initial data
  _r("Bootmethods", "loadStudentStats", function(store, Application) {
    Window.executeService(
      "studentmgmt.Student.Query",
      { limit: 5 },
      null,
      function(response) {
        console.log('Loaded recent students:', response.data);
        // Update store or state with initial data
      },
      function(error) {
        console.error('Failed to load students:', error);
      }
    );
  });
}

export {
  Initialize
}
```

## File Location

The `Initialize` method must be exported from your plugin's main JavaScript file:

```
dev/plugins/studentmgmt/
└── src/
    └── ui/
        └── js/
            └── index.js  ← Initialize method here
```

## Build Configuration

Ensure your `config/ui/build.yml` includes the js directory:

```yaml
js:
  externals:
    - react
    - reactpages
  entry: src/ui/js/index.js
```

## Execution Flow

1. **Application Starts**: `reactapplication` plugin loads
2. **Plugins Discovered**: All UI plugins are discovered
3. **Initialize Called**: Each plugin's `Initialize` method is called
4. **Registration Phase**: Plugins register their components, handlers
5. **Boot Methods Run**: Registered boot methods execute
6. **Application Ready**: UI is rendered and ready for user interaction

## Best Practices

1. **Keep it Lightweight**: `Initialize` should complete quickly
2. **Async Operations**: Use boot methods for async data loading
3. **Error Handling**: Wrap initialization code in try-catch
4. **Logging**: Log initialization for debugging
5. **Module Scope**: Store references at module level for access across your plugin
6. **Properties**: Always set up properties access for localization support

## Common Pitfalls

### ❌ Wrong: Using incorrect signature
```javascript
export const Initialize = (context) => {
  // This won't work!
}
```

### ✅ Correct: Using standard signature
```javascript
function Initialize(appName, ins, mod, settings, def, req) {
  module = this;
  module.properties = Application.Properties[ins];
  module.settings = settings;
}
```

### ❌ Wrong: Not setting module reference
```javascript
function Initialize(appName, ins, mod, settings, def, req) {
  // Missing module = this
  // Properties won't be accessible
}
```

### ✅ Correct: Setting module reference
```javascript
var module;

function Initialize(appName, ins, mod, settings, def, req) {
  module = this;
  module.properties = Application.Properties[ins];
}
```

## See Also

- [UI Plugin Structure](creating-ui-plugins.md)
- [Properties and Localization](../server-development/properties.md)
- [Actions](actions.md) - Using Window.executeAction
- [Services](services.md) - Calling backend services
