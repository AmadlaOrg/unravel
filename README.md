# unravel

System discovery CLI for the [Amadla](https://github.com/AmadlaOrg) ecosystem. Discovers existing system state and outputs it as HERY entities.

## Usage

```bash
# Discover all system state (runs all unravel-* plugins)
unravel discover

# Discover from a specific plugin
unravel discover --from system

# Filter to a specific entity type
unravel discover --type network

# Drift detection pipeline
unravel discover | judge audit

# List discovered plugins
unravel plugins
```

## Plugin Protocol

Unravel discovers `unravel-*` binaries on PATH. Each plugin implements:

- `info` — JSON metadata (name, version, backend, description, supports)
- `discover` — outputs HERY entities as JSON to stdout
- `discover --type <entity-type>` — filtered discovery

Exit codes: `0` success, `1` failure, `2` usage error. Data to stdout, diagnostics to stderr.

## Design

- **Stateless** — no daemon, no caching. Discovers and outputs.
- **UNIX philosophy** — pipe output to files, judge, lighthouse, or any tool.
- **Plugin-based** — `unravel-*` plugins extend discovery to new backends.

## License

MIT
