# Data Views

Display and browse data.

## View Definition

Views can fetch data in two ways:

### Option 1: Direct Service Call

```yaml
serviceName: MyEntity.Query
pagination: true
pageSize: 20

item:
  type: custom
  blockId: my-card
  
filters:
  - name: Status
    field: Status
    type: select
```

### Option 2: Using Dataset

```yaml
dataset: my_dataset
pagination: true
pageSize: 20

item:
  type: block
  id: my-card
  
filters:
  Status: [[ctx.someValue]]
```

**When to use each**:
- Use `serviceName:` for simple, direct queries
- Use `dataset:` for reusable queries with complex logic, transformations, or multiple consumers


See [Views Chapter](../../tutorials/student-management/07-views-lists.md)
