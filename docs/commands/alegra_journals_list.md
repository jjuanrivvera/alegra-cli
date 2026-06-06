---
title: alegra journals list
---

## alegra journals list

List journals

```
alegra journals list [flags]
```

### Examples

```
  alegra journals list
  alegra journals list --limit 30 --all -o json
  alegra journals list --count
  alegra journals list --date <value>
  alegra journals list --param <api_param>=<value>
```

### Options

```
      --all                      Fetch all pages
      --client-id string         Filter by client ID
      --client-name string       Filter by client name
      --count                    Print only the total number of matching records
      --date string              Filter by date (YYYY-MM-DD)
  -h, --help                     help for list
      --limit int                Max records per page (max 30)
      --number string            Filter by numbering
      --observations string      Filter by observations
      --order-direction string   Sort direction: ASC or DESC
      --order-field string       Field to sort by (date, name, reference, observations)
      --param stringArray        Arbitrary API query parameter: key=value (repeatable; e.g. --param date_after=2026-01-01)
  -q, --query string             Free-text search
      --reference string         Filter by reference
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

* [alegra journals](alegra_journals.md)	 - Manage accounting journal entries (comprobantes contables)

