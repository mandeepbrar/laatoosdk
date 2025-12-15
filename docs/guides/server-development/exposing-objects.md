# Exposing Objects in Laatoo Plugins

Laatoo plugins can expose Go objects to the server, making them available for configuration and use by other components. This is done through the `manifest.go` file and the `config/server/objects` directory.

## The Manifest File

Every server-side Go plugin must have a `manifest.go` file in its source root (`src/server/go/`). This file exports a `Manifest` function that returns a list of `core.PluginComponent` objects.

### Example `manifest.go`

```go
package main

import (
	"laatoo.io/sdk/server/core"
)

// Manifest exposes components to the server
func Manifest(provider core.MetaDataProvider) []core.PluginComponent {
	return []core.PluginComponent{
        // Expose a Service Factory
        core.PluginComponent{Object: MyServiceFactory{}},
        
        // Expose Services
        core.PluginComponent{Object: MyService{}},
        
        // Expose a standalone object
        core.PluginComponent{Object: &MyCustomObject{}},
    }
}
```

## Configuring Exposed Objects

Once an object is exposed via the manifest, you can configure instances of it using YAML files in the `config/server/objects` directory of your plugin.

### Directory Structure

```
my-plugin/
├── config/
│   └── server/
│       └── objects/
│           └── myplugin.MyCustomObject.yml
├── src/
│   └── server/
│       └── go/
│           ├── manifest.go
│           └── myobject.go
```

### Object Configuration

The configuration file specifies the object type (matching the struct name exposed in the manifest) and its settings.

**`config/server/objects/myplugin.MyCustomObject.yml`**:

```yaml
type: component
description: My custom object description
configurations:
  apiKey:
    description: API key for external service
    type: string
    required: true
  timeout:
    description: Request timeout in seconds
    type: int
    default: 30
```

## Using Exposed Objects

Exposed objects are initialized by the server on startup. They can implement various interfaces to integrate with the system:

- **`core.ServiceFactory`**: To create services.
- **`core.ServerElement`**: To extend server functionality.
- **`core.Initializer`**: To perform startup logic.

### Example Object Implementation

```go
type MyCustomObject struct {
    ApiKey string
    Timeout int
}

// Initialize is called by the server with the configuration
func (o *MyCustomObject) Initialize(ctx core.ServerContext, conf config.Config) error {
    o.ApiKey, _ = conf.GetString(ctx, "settings.apiKey")
    o.Timeout, _ = conf.GetInt(ctx, "settings.timeout")
    return nil
}
```
