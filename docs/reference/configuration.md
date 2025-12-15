# Configuration Reference

All configuration file formats.

## Entity Configuration

```yaml
name: EntityName
attributes:
  - name: FieldName
    type: string|int|float|boolean|date|datetime|json
    primary: true|false
    required: true|false
    unique: true|false
    default: value
indexes:
  - name: idx_name
    columns: [field1, field2]
```

## Service Configuration

```yaml
servicemethod: plugin.ServiceName.Invoke
description: Service description
permissions:
  - permission:resource
```

## Channel Configuration

```yaml
parent: root
method: GET|POST|PUT|DELETE
service: ServiceName
path: /api/path
authenticate: true|false
roles: [Role1, Role2]
```

## Workflow Configuration

```yaml
name: WorkflowName
taskqueues:
  - queue-name
workflow:
  sequence:
    elements:
      - activity:
          name: ActivityName
          service: ServiceName
          arguments:
            param: value
```

See tutorial chapters for complete examples.
