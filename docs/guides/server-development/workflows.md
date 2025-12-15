# Workflow Engine

Build automated business processes.

## Workflow Structure

```yaml
name: MyWorkflow
workflow:
  sequence:
    elements:
      - activity:
          name: Step1
          service: Service1
          arguments:
            param: value
          result: result1
      - activity:
          name: Step2
          service: Service2
```

See [Workflows Chapter](../../tutorials/student-management/04-workflows.md) for complete examples.
