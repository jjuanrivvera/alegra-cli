---
title: alegra recurring-invoices
---

## alegra recurring-invoices

Manage recurring invoices (facturas recurrentes)

### Synopsis

Manage recurring invoices: templates that automatically generate sales invoices on a schedule.

### Options

```
  -h, --help   help for recurring-invoices
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
* [alegra recurring-invoices create](alegra_recurring-invoices_create.md)	 - Create a recurring-invoice
* [alegra recurring-invoices delete](alegra_recurring-invoices_delete.md)	 - Delete a recurring-invoice by ID
* [alegra recurring-invoices export](alegra_recurring-invoices_export.md)	 - Export all recurring-invoices to CSV or JSON
* [alegra recurring-invoices get](alegra_recurring-invoices_get.md)	 - Get a single recurring-invoice by ID
* [alegra recurring-invoices import](alegra_recurring-invoices_import.md)	 - Bulk-create recurring-invoices from a CSV file
* [alegra recurring-invoices list](alegra_recurring-invoices_list.md)	 - List recurring-invoices
* [alegra recurring-invoices update](alegra_recurring-invoices_update.md)	 - Update a recurring-invoice by ID

