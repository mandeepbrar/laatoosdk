# Building Services

Complete guide to creating services in Laatoo.

## Service Structure

```go
package main

import (
    "fmt"
    "laatoo.io/sdk/server/core"
    "laatoo.io/sdk/server/errors"
    "laatoo.io/sdk/server/log"
)

// Define constants for parameters
const (
    OPERATION = "operation"
)

// MyService struct embedding core.Service
type MyService struct {
    core.Service
    // Dependencies
    dataManager elements.DataManager
}

func (s *MyService) Initialize(ctx core.ServerContext, conf config.Config) error {
    // Setup and define parameters
    return nil
}

func (s *MyService) Start(ctx core.ServerContext) error {
    // Get dependencies
    return nil
}

// Invoke method to handle requests
func (svc *MyService) Invoke(ctx core.RequestContext) (err error) {
    // Get parameters from context
    operation, _ := ctx.GetStringParam(OPERATION)
    
    // Switch on operation
    switch operation {
    case "doSomething":
        return svc.doSomething(ctx)
    default:
        return errors.BadArg(ctx, OPERATION)
    }
}

// Helper method for specific operation
func (svc *MyService) doSomething(ctx core.RequestContext) error {
    // Logic implementation
    data := map[string]interface{}{
        "message": "Operation successful",
    }
    
    // Set response
    ctx.SetResponse(core.SuccessResponse(data))
    return nil
}
```

## See Tutorial
- [Server Plugin Chapter](../../tutorials/student-management/03-server-plugin.md)
