# Your First Laatoo Solution

This guide walks you through creating your first Laatoo solution from scratch. You'll create a simple application with a server plugin and basic UI.

## What We'll Build

A simple "Hello Laatoo" application with:
- A solution and application structure
- A server plugin with one service
- A UI plugin with one page
- Local deployment

## Step 1: Create a Solution

A solution is the top-level container for your entire system.

```bash
# Navigate to your workspace
cd ~/laatoo-workspace/solutions

# Create a new solution
laatoo solution create hello-laatoo

# Navigate into the solution
cd hello-laatoo
```

**What was created**:
```
hello-laatoo/
├── applications/          # Empty, we'll add an app next
├── config/
│   ├── config.yml         # Solution configuration
│   ├── modules/           # Server modules (user, role, etc.)
│   └── engines/           # Server engines (http, etc.)
├── data/                  # Shared data
├── dev/                   # Development
│   └── plugins/           # Plugin development
├── environment_local.yml  # Local environment variables
└── laatoo.yml            # Solution configuration
```

## Step 2: Create an Application

Applications are logical groupings of functionality within a solution.

```bash
# Create application directory
mkdir -p applications/webapp

# Create application config
mkdir -p applications/webapp/config
cat > applications/webapp/config/config.yml << EOF
name: webapp
description: My first web application
EOF

# Create default isolation
mkdir -p applications/webapp/isolations/default/config
cat > applications/webapp/isolations/default/config/config.yml << EOF
name: default
description: Default isolation
EOF
```

**What was created**:
```
applications/webapp/
├── config/
│   ├── config.yml         # Application configuration
│   └── modules/           # Plugins for this app (empty)
└── isolations/
    └── default/
        └── config/
            └── config.yml
```

## Step 3: Create a Server Plugin

Now let's create a plugin that provides a simple service.

```bash
# Create server plugin
cd dev/plugins
laatoo plugin create hello-server -t server

# Navigate into the plugin
cd hello-server
```

**Plugin structure created**:
```
hello-server/
├── config/
│   ├── config.yml         # Plugin metadata
│   └── services/          # Service definitions
├── files/                 # Static files
└── src/
    └── server/
        └── go/            # Go source code
```

### Define a Service

Create a service configuration:

```bash
# Create service definition
laatoo plugin addservice HelloService -p hello-server
```

This creates `config/services/HelloService.yml`:

```yaml
servicemethod: HelloService.HandleRequest
```

### Implement the Service

Create the service implementation in `src/server/go/helloservice.go`:

```go
package main

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/log"
)

type HelloService struct {
	core.Service
}

func (s *HelloService) Initialize(ctx core.ServerContext, conf config.Config) error {
	// Define service parameters
	err := s.AddStringParam(ctx, "name", "Name to greet")
	if err != nil {
		return err
	}
	
	s.SetDescription(ctx, "Simple hello service")
	return nil
}

func (s *HelloService) Invoke(ctx core.RequestContext) error {
	// Get parameter from request
	name, ok := ctx.GetStringParam("name")
	if !ok || name == "" {
		name = "World"
	}
	
	message := "Hello, " + name + "!"
	
	log.Info(ctx, "Greeting generated", "name", name)
	
	// Set response
	ctx.SetResponse(&core.Response{
		Status: core.StatusOK,
		Data: map[string]interface{}{
			"message":   message,
			"timestamp": ctx.Now(),
		},
	})
	
	return nil
}
```

### Create Service Factory

Create `src/server/go/factory.go`:

```go
package main

import (
	"laatoo.io/sdk/config"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/elements"
	"laatoo.io/sdk/server/errors"
)

type HelloFactory struct {
	core.ServiceFactory
}

func (f *HelloFactory) CreateService(
	ctx core.ServerContext,
	serviceName string,
	method string,
	conf config.Config,
) (core.Service, error) {
	
	service := &HelloService{}
	
	return service, nil
}

// Export factory for plugin system
var Factory core.ServiceFactory = &HelloFactory{}
```

### Create Plugin Manifest

Create `src/server/go/manifest.go` to expose your factory to the server:

```go
package main

import (
	"laatoo.io/sdk/server/core"
)

func Manifest(provider core.MetaDataProvider) []core.PluginComponent {
	return []core.PluginComponent{
		core.PluginComponent{Object: HelloFactory{}},
		core.PluginComponent{Object: HelloService{}},
	}
}
```

### Initialize Go Module

```bash
cd src/server/go
go mod init hello-server
go mod tidy
cd ../../..
```

### Build the Plugin

```bash
# Build the server plugin
laatoo plugin build hello-server
```

This compiles the Go code and creates `bin/` with the compiled plugin.

## Step 4: Create a Channel

Channels expose services via HTTP/gRPC. Let's create an HTTP endpoint:

```bash
laatoo plugin addchannel hello-api -p hello-server \
  --channelParent root \
  --service HelloService \
  --channelpath /api/hello
```

This creates `config/channels/hello-api.yml`:

```yaml
parent: root
method: POST
service: HelloService
path: /api/hello
```

Rebuild the plugin:
```bash
laatoo plugin build hello-server
```

## Step 5: Install Plugin in Application

```bash
# Create modules directory
mkdir -p ../../applications/webapp/config/modules

# Copy plugin to application
cp -r bin ../../applications/webapp/config/modules/hello-server
```

## Step 6: Configure Database Connection

Edit `applications/webapp/isolations/default/config/config.yml`:

```yaml
name: default
description: Default isolation

# Database configuration
database:
  type: postgres
  host: localhost
  port: 5432
  database: laatoo
  user: laatoo
  password: laatoo
```

## Step 7: Run the Solution

```bash
# Return to solution root
cd ../../../..

# Run the solution
laatoo solution run hello-laatoo
```

The server will start and you should see:
```
[INFO] Starting Laatoo Solution: hello-laatoo
[INFO] Loading application: webapp
[INFO] Loading isolation: default
[INFO] Loading module: hello-server
[INFO] Starting service: HelloService
[INFO] Starting channel: hello-api on /api/hello
[INFO] Server ready on http://localhost:8080
```

## Step 8: Test the Service

In another terminal:

```bash
# Test the service
curl -X POST http://localhost:8080/api/hello \
  -H "Content-Type: application/json" \
  -d '{"name": "Laatoo"}'
```

Expected response:
```json
{
  "data": {
    "message": "Hello, Laatoo!",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

## Step 9: Create a UI Plugin (Optional)

Let's add a simple UI to call our service.

```bash
# Create UI plugin
cd dev/plugins
laatoo plugin create hello-ui -t ui
cd hello-ui
```

### Create a Page

```bash
laatoo ui addpage home -p hello-ui
```

Edit `config/ui/pages/home.yml`:

```yaml
route: "/"
authenticate: false
component:
  type: block
  id: homepage
```

### Create a Block

Create `config/ui/blocks/homepage.xml`:

```xml
<Block className="hello_container">
  <h1 module="html">Hello Laatoo!</h1>
  <Block className="greeting">
    <Action name="sayHello">
      <Content>Say Hello</Content>
    </Action>
    <Block id="result"></Block>
  </Block>
</Block>
```

### Create an Action

Edit `config/ui/actions/sayHello.yml`:

```yaml
actiontype: callmethod
actionparams:
  method: sayHello
```

### Create the UI Logic

Create `src/ui/Initialize.js`:

```javascript
export const Initialize = (context) => {
  return {
    sayHello: async () => {
      const response = await context.callService('HelloService', {
        name: 'from UI'
      });
      
      const result = document.getElementById('result');
      result.innerHTML = `<p>${response.data.message}</p>`;
    }
  };
};
```

### Build and Install UI Plugin

```bash
# Build UI plugin
laatoo plugin build hello-ui --getbuildpackages

# Install in application
cp -r bin ../../applications/webapp/config/modules/hello-ui
```

## Step 10: View the UI

1. Restart the solution:
```bash
cd ../../../..
laatoo solution run hello-laatoo
```

2. Open browser to `http://localhost:8080/`

3. Click "Say Hello" button

You should see "Hello, from UI!" displayed.

## Understanding What You Built

### Solution Structure

```
hello-laatoo/
├── applications/webapp/     # Your application
│   ├── config/
│   │   └── modules/         # Installed plugins
│   │       ├── hello-server/
│   │       └── hello-ui/
│   └── isolations/default/  # Default tenant
└── dev/plugins/             # Plugin development
    ├── hello-server/        # Server logic
    └── hello-ui/            # UI components
```

### Request Flow

```
Browser → Channel (/api/hello) → HelloService → Response
```

### Key Concepts Learned

1. **Solutions** contain applications
2. **Applications** provide logical separation
3. **Isolations** enable multi-tenancy
4. **Server Plugins** provide services
5. **Services** implement business logic
6. **Channels** expose services as APIs
7. **UI Plugins** provide user interfaces

## Next Steps

### Enhance Your Application

1. **Add Database**: [Entities Guide](../guides/server-development/entities.md)
2. **Add Authentication**: [Security Guide](../guides/server-development/security.md)
3. **Add More UI**: [UI Development Guide](../guides/ui-development/creating-ui-plugins.md)

### Learn More

1. **Complete Tutorial**: [Student Management System](../tutorials/student-management/00-overview.md)
2. **Server Development**: [Creating Plugins](../guides/server-development/creating-plugins.md)
3. **UI Development**: [Creating UI Plugins](../guides/ui-development/creating-ui-plugins.md)

## Troubleshooting

**Issue**: Service not found

**Solution**: Ensure plugin is built and copied to `config/modules/`

---

**Issue**: Database connection failed

**Solution**: Check PostgreSQL is running and credentials match `config.yml`

---

**Issue**: UI assets not loading

**Solution**: Rebuild UI plugin with `--getbuildpackages` flag

## Summary

Congratulations! You've created your first Laatoo solution with:
- ✅ Solution and application structure
- ✅ Server plugin with a service
- ✅ HTTP API endpoint
- ✅ UI plugin with a page
- ✅ Full end-to-end working application

You're now ready to explore more advanced features and build real-world applications with Laatoo!
