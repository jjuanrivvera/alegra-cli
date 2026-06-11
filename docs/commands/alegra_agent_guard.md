---
title: alegra agent guard
---

## alegra agent guard

Generate agent-safety config that blocks destructive alegra operations

### Synopsis

guard generates the permission rules and hooks that stop an AI agent from
running destructive alegra operations, derived from the live command tree so the
list is always complete.

By default it hard-blocks the irreversible actions (delete, void, emit, stamp,
close, cancel) and makes ordinary writes (create, update, import) require
approval; read operations stay allowed. Pass --all-writes to block writes too.

Output is printed for review by default; pass --write to install it. See the
Agent Safety guide: https://jjuanrivvera.github.io/alegra-cli/user-guide/agent-safety/

```
alegra agent guard [flags]
```

### Examples

```
  alegra agent guard --host claude-code
  alegra agent guard --host codex
  alegra agent guard --host opencode --all-writes
  alegra agent guard --host claude-code --write
```

### Options

```
      --all-writes    Also block create/update/import (default: those require approval)
  -h, --help          help for guard
      --host string   Target agent host: claude-code, codex, opencode
      --write         Write the config/hook files instead of printing them
```

### Options inherited from parent commands

```
      --base-url string             Override the API base URL (env: ALEGRA_BASE_URL)
      --columns strings             Comma-separated columns for table/csv output
      --dry-run                     Print the equivalent curl request without sending it
      --no-color                    Disable colored output (also respects the NO_COLOR env var)
  -o, --output string               Output format: table, json, yaml, csv (env: ALEGRA_OUTPUT)
      --profile string              Configuration profile to use (env: ALEGRA_PROFILE)
      --requests-per-second float   Client-side rate limit (default from config)
      --show-token                  In --dry-run, do not redact the Authorization header
  -v, --verbose                     Enable verbose (debug) logging to stderr
```

### SEE ALSO

* [alegra agent](alegra_agent.md)	 - Helpers for running alegra under an AI agent

