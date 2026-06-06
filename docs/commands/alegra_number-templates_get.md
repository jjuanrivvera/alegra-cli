---
title: alegra number-templates get
---

## alegra number-templates get

Get a single number-template by ID

```
alegra number-templates get <id> [flags]
```

### Examples

```
  alegra number-templates get <id>
  alegra number-templates get <id> -o json
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
  -o, --output string               Output format: table, json, yaml, csv (env: ALEGRA_OUTPUT)
      --profile string              Configuration profile to use (env: ALEGRA_PROFILE)
      --requests-per-second float   Client-side rate limit (default from config)
      --show-token                  In --dry-run, do not redact the Authorization header
  -v, --verbose                     Enable verbose (debug) logging to stderr
```

### SEE ALSO

* [alegra number-templates](alegra_number-templates.md)	 - Manage document numberings (numeraciones de facturación)

