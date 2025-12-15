# Creating UI Plugins

Build user interfaces for Laatoo applications.

## Create UI Plugin

```bash
laatoo plugin create my-ui -t ui
```

## See Also

- [Design System](design-system.md)
- [Forms Guide](forms.md)

## Structure

```
my-ui/
├── config/
│   └── ui/
│       └── build.yml
└── src/
    └── ui/
        ├── js/
        │   └── index.js      # Initialize method
        ├── styles/            # CSS/SCSS
        └── registry/          # UI definitions
            ├── Pages/
            ├── Forms/
            ├── Views/
            ├── Blocks/
            ├── Actions/
            └── Menus/
```

## Build Configuration

```yaml
js:
  externals:
    - react
    - reactpages
  packages:
    react: "^18.2.0"
```

See [UI Plugin Chapter](../../tutorials/student-management/05-ui-plugin.md)
