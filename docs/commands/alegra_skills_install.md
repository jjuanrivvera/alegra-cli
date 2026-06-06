---
title: alegra skills install
---

## alegra skills install

Write the bundled skill into an agent's skills directory

```
alegra skills install [flags]
```

### Options

```
      --agent string   Target agent: claude, codex, copilot, cursor, gemini, opencode, windsurf (default "claude")
      --dir string     Explicit destination base directory (overrides --agent/--global)
  -g, --global         Install for all projects (user-level) instead of the current project
  -h, --help           help for install
```

### Options inherited from parent commands

```
      --base-url string             Override the API base URL (env: ALEGRA_BASE_URL)
      --columns strings             Comma-separated columns for table/csv output
      --dry-run                     Print the equivalent curl request without sending it
  -o, --output string               Output format: table, json, yaml, csv (env: ALEGRA_OUTPUT)
      --profile string              Configuration profile to use (env: ALEGRA_PROFILE)
      --requests-per-second float   Client-side rate limit (default from config)
      --show-token                  In --dry-run, do not redact the Authorization header
  -v, --verbose                     Enable verbose (debug) logging to stderr
```

### SEE ALSO

* [alegra skills](alegra_skills.md)	 - Install this CLI's AI-agent skill into Claude, Cursor, and other agents

