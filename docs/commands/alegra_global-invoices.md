---
title: alegra global-invoices
---

## alegra global-invoices

Manage global invoices (facturas globales)

### Synopsis

Manage Alegra global invoices, which group sale tickets into a single CFDI for the general public.

### Options

```
  -h, --help   help for global-invoices
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

* [alegra](alegra.md)	 - Alegra accounting system CLI
* [alegra global-invoices create](alegra_global-invoices_create.md)	 - Create a global-invoice
* [alegra global-invoices delete](alegra_global-invoices_delete.md)	 - Delete a global-invoice by ID
* [alegra global-invoices export](alegra_global-invoices_export.md)	 - Export all global-invoices to CSV or JSON
* [alegra global-invoices get](alegra_global-invoices_get.md)	 - Get a single global-invoice by ID
* [alegra global-invoices import](alegra_global-invoices_import.md)	 - Bulk-create global-invoices from a CSV file
* [alegra global-invoices list](alegra_global-invoices_list.md)	 - List global-invoices
* [alegra global-invoices update](alegra_global-invoices_update.md)	 - Update a global-invoice by ID

