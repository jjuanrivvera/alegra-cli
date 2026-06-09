---
title: alegra mcp vscode enable
---

## alegra mcp vscode enable

Add server to VSCode config

### Synopsis

Add this application as an MCP server in VSCode

```
alegra mcp vscode enable [flags]
```

### Options

```
      --config-path string   Path to VSCode config file
  -e, --env stringToString   Environment variables (e.g., --env KEY1=value1 --env KEY2=value2) (default [])
  -h, --help                 help for enable
      --log-level string     Log level (debug, info, warn, error)
      --server-name string   Name for the MCP server (default: derived from executable name)
      --workspace            Add to workspace settings (.vscode/mcp.json) instead of user settings
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

