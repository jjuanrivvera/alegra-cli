---
title: alegra items list
---

## alegra items list

List items

```
alegra items list [flags]
```

### Examples

```
  alegra items list
  alegra items list --limit 30 --all -o json
  alegra items list --count
  alegra items list --name <value>
  alegra items list --param <api_param>=<value>
```

### Options

```
      --all                      Fetch all pages
      --count                    Print only the total number of matching records
  -h, --help                     help for list
      --limit int                Max records per page (max 30)
      --name string              Filter by name
      --order-direction string   Sort direction: ASC or DESC
      --order-field string       Field to sort by (id, name, reference, price)
      --param stringArray        Arbitrary API query parameter: key=value (repeatable; e.g. --param date_after=2026-01-01)
  -q, --query string             Free-text search
      --start int                Offset to start from (pagination)
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

* [alegra items](alegra_items.md)	 - Manage items (products and services)

