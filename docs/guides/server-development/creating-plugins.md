# Creating Server Plugins

Guide to building server plugins.

## Create Plugin

```bash
laatoo plugin create myplugin -t server
```

## Structure

```
myplugin/
├── config/
│   ├── config.yml
│   └── server/
│       ├── objects/
│       ├── services/
│       ├── factories/
│       └── modules/
└── src/server/go/
```

## Plugin Configuration

```yaml
name: myplugin
version: 1.0.0
description: My plugin

entities:
  - MyEntity
```

See [Tutorial](../../tutorials/student-management/01-setup-solution.md)
