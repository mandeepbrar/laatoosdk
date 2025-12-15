# Background Tasks

Queue and process asynchronous tasks.

## Push Task

```go
err := ctx.PushTask("my-queue", map[string]interface{}{
    "data": value,
})
```

## Configure Queue

```yaml
taskqueues:
  my-queue:
    workers: 5
    maxConcurrent: 10
```

See [Workflows Chapter](../../tutorials/student-management/04-workflows.md)
