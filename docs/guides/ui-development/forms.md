# Forms & Validation

Create data entry forms.

## Form Definition

Forms are defined using XML in `src/ui/registry/Forms/`.

```xml
<Form form="my_form" module="formikforms" action="submitForm" successRedirect="/success">
  <Block className="row">
    <Block className="col-md-6">
      <Field type="text" name="Name" label="Name" required="true"/>
    </Block>
    <Block className="col-md-6">
      <Field type="email" name="Email" label="Email" validation="email"/>
    </Block>
  </Block>
  <Button type="submit">Submit</Button>
</Form>
```

## Page Integration

Forms must be included in a page to be accessible.

```yaml
route: /my-form-page
component:
  type: layout
  layout: 1col
  items:
    - type: form
      id: my_form
```

## Field Types
- string, text, email, password
- number, date, datetime
- select, checkbox, radio
- file, textarea

See [Forms Chapter](../../tutorials/student-management/06-forms.md)
