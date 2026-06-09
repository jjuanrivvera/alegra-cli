---
title: alegra variant-attributes import
---

## alegra variant-attributes import

Bulk-create variant-attributes from a CSV file

### Synopsis

Create one variant-attribute per CSV row.

The header row names the fields; use --map to rename columns to API fields and
dotted paths for nested objects (e.g. --map 'NIT=identification.number'). Apply
constant fields to every row with --set. Rows are processed independently;
failures are reported and do not stop the run.

```
alegra variant-attributes import [flags]
```

### Examples

```
  alegra variant-attributes import --file rows.csv
  alegra contacts import -f clients.csv \
    --map 'Name=name,NIT=identification.number' \
    --set 'identification.type=NIT' --set 'type=["client"]'
```

### Options

```
  -f, --file string       CSV file to import (required)
  -h, --help              help for import
      --map stringArray   Map a CSV column to a field path: column=field.path (repeatable)
      --set stringArray   Constant field applied to every row: key=value (repeatable)
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

* [alegra variant-attributes](alegra_variant-attributes.md)	 - Manage item variant attributes (e.g. Color, Talla)

