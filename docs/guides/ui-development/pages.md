# Pages & Routes

Create routed pages in your application.

## Page Definition

```yaml
route: "/mypage"
authenticate: true
roles: [Admin]
component:
  type: layout
  layout: 2col
  leftcol:
    type: menu
    id: main-menu
  rightcol:
    type: block
    id: content
```

## Layouts
- `1col` - Single column
- `2col` - Two columns (sidebar + content)
- `3col` - Three columns
- `custom` - Custom layout

See [UI Plugin Chapter](../../tutorials/student-management/05-ui-plugin.md)
