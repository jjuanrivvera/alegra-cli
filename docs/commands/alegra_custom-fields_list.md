---
title: alegra custom-fields list
---

## alegra custom-fields list

List custom-fields

```
alegra custom-fields list [flags]
```

### Examples

```
  alegra custom-fields list
  alegra custom-fields list --limit 30 --all -o json
  alegra custom-fields list --count
  alegra custom-fields list --name <value>
  alegra custom-fields list --param <api_param>=<value>
```

### Options

```
      --all                      Fetch all pages
      --count                    Print only the total number of matching records
      --description string       Filter by text contained in the field description
  -h, --help                     help for list
      --limit int                Max records per page (max 30)
      --name string              Filter by text contained in the field name
      --order-direction string   Sort direction: ASC or DESC
      --order-field string       Field to sort by (id, name)
      --param stringArray        Arbitrary API query parameter: key=value (repeatable; e.g. --param date_after=2026-01-01)
  -q, --query string             Free-text search
      --resource-type string     Filter by associated resource type (e.g. item)
      --since string             Start of date range (YYYY-MM-DD, today, this-month, last-month, 7d, 3m, ...)
      --start int                Offset to start from (pagination)
      --status string            Filter by status: active or inactive
      --type string              Filter by field type: text, int, date, boolean, optionsList
      --until string             End of date range (same formats as --since)
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

* [alegra custom-fields](alegra_custom-fields.md)	 - Manage custom fields (campos adicionales)

