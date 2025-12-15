# CLI Command Reference

Complete Laatoo CLI reference.

## Global Flags
- `--config` - Config file path
- `--verbose` - Verbose output
- `--withdebugger` - Attach debugger

## Plugin Commands

### Create Plugin
```bash
laatoo plugin create [name] -t [type]
```
Types: server, ui, min, full

### Build Plugin
```bash
laatoo plugin build [name] --buildmode [mode]
```
Modes: local, prod, dockerdebug

### Add Elements
```bash
laatoo plugin addservice [name] -p [plugin]
laatoo plugin addchannel [name] -p [plugin]
laatoo plugin addrule [name] -p [plugin]
```

## Solution Commands

### Create Solution
```bash
laatoo solution create [name]
```

### Run Solution
```bash
laatoo solution run [name] -e [env] -m [mode]
```

## UI Commands

```bash
laatoo ui addpage [name] -p [plugin]
laatoo ui addform [name] -p [plugin]
laatoo ui addview [name] -p [plugin]
laatoo ui addblock [name] -p [plugin]
laatoo ui addaction [name] -p [plugin]
```

See [LAATOO_CLI_GUIDE.md](/home/mandeep/goprogs/src/laatoo/docs/LAATOO_CLI_GUIDE.md) for full details.
