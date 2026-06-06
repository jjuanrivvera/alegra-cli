---
title: alegra invoices stamp
---

## alegra invoices stamp

Stamp up to 10 invoices (low-level; prefer `emit`)

```
alegra invoices stamp [flags]
```

### Options

```
  -d, --data string       Request body as a JSON string
  -f, --file string       Read JSON request body from a file (use - for stdin)
  -h, --help              help for stamp
      --set stringArray   Set a top-level field: key=value (value parsed as JSON when valid). Repeatable. For nested documents (e.g. invoice items[]) use --file.
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

* [alegra invoices](alegra_invoices.md)	 - Manage sales invoices

