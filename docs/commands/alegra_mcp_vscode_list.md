---
title: alegra mcp vscode list
---

## alegra mcp vscode list

Show VSCode MCP servers

### Synopsis

Show all MCP servers configured in VSCode

```
alegra mcp vscode list [flags]
```

### Options

```
      --config-path string   Path to VSCode config file
  -h, --help                 help for list
      --workspace            List from workspace settings (.vscode/mcp.json) instead of user settings
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

* [alegra mcp vscode](alegra_mcp_vscode.md)	 - Manage VSCode MCP servers

