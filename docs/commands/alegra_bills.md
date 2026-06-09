---
title: alegra bills
---

## alegra bills

Manage provider bills (facturas de proveedor)

### Synopsis

Manage provider bills (facturas de proveedor) — purchases you owe. Supports attachments, comments, advance application, perceptions/retentions, and importing received Colombian e-invoices by CUFE.

### Options

```
  -h, --help   help for bills
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
* [alegra bills advances](alegra_bills_advances.md)	 - Apply provider advances to a bill
* [alegra bills attach](alegra_bills_attach.md)	 - Attach a file to a bill (JSON body with a base64 'file' field)
* [alegra bills attachment-delete](alegra_bills_attachment-delete.md)	 - Delete a bill attachment by file ID
* [alegra bills close](alegra_bills_close.md)	 - Close a bill with pending balance
* [alegra bills comment-delete](alegra_bills_comment-delete.md)	 - Delete a comment from a bill
* [alegra bills comment-update](alegra_bills_comment-update.md)	 - Edit a comment on a bill
* [alegra bills comments](alegra_bills_comments.md)	 - Add a comment to a bill
* [alegra bills create](alegra_bills_create.md)	 - Create a bill
* [alegra bills delete](alegra_bills_delete.md)	 - Delete a bill by ID
* [alegra bills export](alegra_bills_export.md)	 - Export all bills to CSV or JSON
* [alegra bills get](alegra_bills_get.md)	 - Get a single bill by ID
* [alegra bills import](alegra_bills_import.md)	 - Bulk-create bills from a CSV file
* [alegra bills import-by-cufe](alegra_bills_import-by-cufe.md)	 - Import a bill by CUFE (Colombia)
* [alegra bills list](alegra_bills_list.md)	 - List bills
* [alegra bills perceptions](alegra_bills_perceptions.md)	 - Replace a bill's perceptions
* [alegra bills retentions](alegra_bills_retentions.md)	 - Replace a bill's retentions
* [alegra bills update](alegra_bills_update.md)	 - Update a bill by ID

