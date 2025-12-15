# API Channels

Expose services via HTTP, gRPC, or WebSocket.

## HTTP Channel

```yaml
parent: root
method: POST
service: MyService
path: /api/myservice
authenticate: true
```

## Channel Types
- HTTP REST
- gRPC
- WebSocket
- GraphQL

See [Tutorial](../../tutorials/student-management/03-server-plugin.md) for examples.
