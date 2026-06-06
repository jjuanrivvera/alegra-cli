---
title: alegra number-templates export
---

## alegra number-templates export

Export all number-templates to CSV or JSON

### Synopsis

Fetch every page of number-templates and write them to a file (or stdout).

```
alegra number-templates export [flags]
```

### Examples

```
  alegra number-templates export > number-templates.csv
  alegra number-templates export --format json --out number-templates.json
  alegra number-templates export --param status=open
```

### Options

```
      --format string       Export format: csv or json (default "csv")
  -h, --help                help for export
      --out string          Write to this file (default: stdout)
      --param stringArray   Filter by an API query parameter: key=value (repeatable)
  -q, --query string        Free-text search
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

