---
title: alegra invoices
---

## alegra invoices

Manage sales invoices

### Synopsis

Manage sales invoices (facturas de venta).

An invoice becomes fiscal/electronic when it carries a numberTemplate.id and stamp.generateStamp:true (emission returns a CUFE in CO / UUID in MX). Prefer the `emit` subcommand for the safe, batched, idempotent flow. Invoices are append-only: reverse or correct a stamped one with a credit note, never by editing.

### Options

```
  -h, --help   help for invoices
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
* [alegra invoices create](alegra_invoices_create.md)	 - Create a invoice
* [alegra invoices delete](alegra_invoices_delete.md)	 - Delete a invoice by ID
* [alegra invoices email](alegra_invoices_email.md)	 - Email an invoice
* [alegra invoices emit](alegra_invoices_emit.md)	 - Emit (stamp) draft invoices electronically, in batches of 10
* [alegra invoices export](alegra_invoices_export.md)	 - Export all invoices to CSV or JSON
* [alegra invoices get](alegra_invoices_get.md)	 - Get a single invoice by ID
* [alegra invoices import](alegra_invoices_import.md)	 - Bulk-create invoices from a CSV file
* [alegra invoices list](alegra_invoices_list.md)	 - List invoices
* [alegra invoices open](alegra_invoices_open.md)	 - Revert a voided invoice
* [alegra invoices preview](alegra_invoices_preview.md)	 - Generate a preview PDF URL
* [alegra invoices stamp](alegra_invoices_stamp.md)	 - Stamp up to 10 invoices (low-level; prefer `emit`)
* [alegra invoices update](alegra_invoices_update.md)	 - Update a invoice by ID
* [alegra invoices void](alegra_invoices_void.md)	 - Void an invoice

