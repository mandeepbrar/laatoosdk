# SDK Interfaces Reference

Key interfaces to implement.

## Server Interfaces

### Service
```go
type Service interface {
    ConfigurableObject
    Initialize(ServerContext, config.Config) error
    Start(ServerContext) error
    Invoke(RequestContext) error
    Stop(ServerContext) error
}
```

### ServiceFactory
```go
type ServiceFactory interface {
    CreateService(ServerContext, string, config.Config) (Service, error)
}
```

### DataComponent
```go
type DataComponent interface {
    GetById(RequestContext, string, string) (interface{}, error)
    Query(RequestContext, map[string]interface{}, string, int, int) ([]interface{}, error)
    Create(RequestContext, interface{}) (interface{}, error)
    Update(RequestContext, string, interface{}, string) (interface{}, error)
    Delete(RequestContext, string) error
}
```

## Context Interfaces

### RequestContext
- GetStringParam(name) (string, bool)
- GetIntParam(name) (int, bool)
- GetParamValue(name) (interface{}, bool)
- SetResponse(*Response)
- Now() time.Time
- GetUser() User
- PushTask(queue, data) error

### ServerContext
- GetService(name) (Service, error)
- GetServerElement(type) interface{}

See SDK documentation for complete details.
