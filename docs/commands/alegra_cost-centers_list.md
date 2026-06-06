---
title: alegra cost-centers list
---

## alegra cost-centers list

List cost-centers

```
alegra cost-centers list [flags]
```

### Options

```
      --all                      Fetch all pages
  -h, --help                     help for list
      --limit int                Max records per page (max 30)
      --name string              Filter by name (partial match)
      --order-direction string   Sort direction: ASC or DESC
      --order-field string       Field to sort by (id, name, code, status)
  -q, --query string             Free-text search
      --start int                Offset to start from (pagination)
      --status string            Filter by status: active or inactive
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

* [alegra cost-centers](alegra_cost-centers.md)	 - Manage cost centers

