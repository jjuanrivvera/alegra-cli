---
title: alegra number-templates import
---

## alegra number-templates import

Bulk-create number-templates from a CSV file

### Synopsis

Create one number-template per CSV row.

The header row names the fields; use --map to rename columns to API fields and
dotted paths for nested objects (e.g. --map 'NIT=identification.number'). Apply
constant fields to every row with --set. Rows are processed independently;
failures are reported and do not stop the run.

```
alegra number-templates import [flags]
```

### Examples

```
  alegra number-templates import --file rows.csv
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
  -o, --output string               Output format: table, json, yaml, csv (env: ALEGRA_OUTPUT)
      --profile string              Configuration profile to use (env: ALEGRA_PROFILE)
      --requests-per-second float   Client-side rate limit (default from config)
      --show-token                  In --dry-run, do not redact the Authorization header
  -v, --verbose                     Enable verbose (debug) logging to stderr
```

### SEE ALSO

* [alegra number-templates](alegra_number-templates.md)	 - Manage document numberings (numeraciones de facturación)

