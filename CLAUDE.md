# CLAUDE.md

## AI Skills

Follow the practices defined in `~/Projects/SiteNetSoft/ai-skills/`:
- `dev-practices/golang/` — Go style, error handling, functions, testing, linting
- `dev-practices/git/` — Git authorship rules, multi-repo workspace patterns

## Project Overview

Unravel is a system discovery CLI tool for the Amadla ecosystem. It discovers existing system state and outputs it as HERY entities. It scans PATH for `unravel-*` plugins and delegates discovery to them.

**Ecosystem context:** Unravel outputs entities matching the current taxonomy (no "Entity" prefix): Application, OS, Package, Security/*, System/*, User, Service, etc. Entity schemas live in `Entities/<Type>/`. Part of the broader pipeline (`raise` -> `lay` -> `enjoin` -> `weaver` -> `waiter`), unravel operates independently to discover and report existing state.

## Build Commands

```bash
make build    # Build for current platform
make test     # Run tests
make clean    # Remove build artifacts
```

## Architecture

**UNIX Plugin Protocol:**
```
unravel discover                          ->  all unravel-* plugins: discover
unravel discover --from <backend>         ->  unravel-<backend> discover
unravel discover --type <entity-type>     ->  unravel-<backend> discover --type <entity-type>
unravel plugins                           ->  scans PATH for unravel-* binaries
```

**Package Structure:**
- `main.go` - CLI entry point (Cobra)
- `plugin/` - Plugin discovery and execution (PATH scanning, subprocess delegation)
- `cmd/` - CLI commands (discover, plugins)

**Plugin Protocol (unravel-* binaries):**
- `info` subcommand -> JSON metadata (name, version, backend, description, supports)
- `discover` subcommand -> outputs HERY entities as JSON to stdout
- `discover --type <entity-type>` -> filtered discovery
- Exit codes: 0 success, 1 failure, 2 usage error
- Data to stdout, diagnostics to stderr

**Stateless:** No daemon mode, no caching. Discovers and outputs. Users pipe to a file if they want to cache.
