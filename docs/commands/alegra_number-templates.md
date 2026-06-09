---
title: alegra number-templates
---

## alegra number-templates

Manage document numberings (numeraciones de facturación)

### Synopsis

Manage document numberings (numeraciones de facturación) — the prefixes, ranges, and DIAN/SAT resolutions that number invoices and other documents.

### Options

```
  -h, --help   help for number-templates
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

* [alegra](alegra.md)	 - Alegra accounting system CLI
* [alegra number-templates create](alegra_number-templates_create.md)	 - Create a number-template
* [alegra number-templates delete](alegra_number-templates_delete.md)	 - Delete a number-template by ID
* [alegra number-templates export](alegra_number-templates_export.md)	 - Export all number-templates to CSV or JSON
* [alegra number-templates get](alegra_number-templates_get.md)	 - Get a single number-template by ID
* [alegra number-templates import](alegra_number-templates_import.md)	 - Bulk-create number-templates from a CSV file
* [alegra number-templates list](alegra_number-templates_list.md)	 - List number-templates
* [alegra number-templates update](alegra_number-templates_update.md)	 - Update a number-template by ID

