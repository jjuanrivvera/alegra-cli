---
title: alegra variant-attributes list
---

## alegra variant-attributes list

List variant-attributes

```
alegra variant-attributes list [flags]
```

### Options

```
      --all                      Fetch all pages
  -h, --help                     help for list
      --limit int                Max records per page (max 30)
      --name string              Filter variant attributes by name
      --options string           Filter variant attributes by options
      --order-direction string   Sort direction: ASC or DESC
      --order-field string       Field to sort by (name)
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

* [alegra variant-attributes](alegra_variant-attributes.md)	 - Manage item variant attributes (e.g. Color, Talla)

