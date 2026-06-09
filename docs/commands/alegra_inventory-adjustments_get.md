---
title: alegra inventory-adjustments get
---

## alegra inventory-adjustments get

Get a single inventory-adjustment by ID

```
alegra inventory-adjustments get <id> [flags]
```

### Examples

```
  alegra inventory-adjustments get <id>
  alegra inventory-adjustments get <id> -o json
```

### Options

```
  -h, --help   help for get
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

* [alegra inventory-adjustments](alegra_inventory-adjustments.md)	 - Manage inventory adjustments (manual stock corrections)

