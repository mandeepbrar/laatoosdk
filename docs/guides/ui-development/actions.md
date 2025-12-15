# Actions & Interactions

Handle user interactions.

## Action Definition

```yaml
actiontype: callmethod
actionparams:
  method: myMethod
  successMessage: "Success!"
```

## Form Submission Action

Actions can be used to handle form submissions by calling a backend service.

```yaml
actiontype: executeservice
servicename: studentmgmt.Student.Create
successaction: navigateToDashboard
params:
  data: {{jsreplace "ctx.formData"}}
```

## Action Types
- callmethod
- callservice
- navigate  
- openinteraction

See [Forms Chapter](../../tutorials/student-management/06-forms.md)
