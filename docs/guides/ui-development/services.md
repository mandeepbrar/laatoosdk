# UI Services

Call backend services from UI.

## Call Service

```javascript
const response = await context.callService('MyService', {
  param1: value1,
  param2: value2
});

if (response.data) {
  // Handle success
}
```

See [UI Plugin Chapter](../../tutorials/student-management/05-ui-plugin.md)
