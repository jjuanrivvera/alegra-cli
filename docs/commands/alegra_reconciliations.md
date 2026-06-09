---
title: alegra reconciliations
---

## alegra reconciliations

Manage bank reconciliations

### Synopsis

Manage bank reconciliations. Create performs a create-or-update (upsert); there is no separate update endpoint.

### Options

```
  -h, --help   help for reconciliations
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
* [alegra reconciliations create](alegra_reconciliations_create.md)	 - Create a reconciliation
* [alegra reconciliations delete](alegra_reconciliations_delete.md)	 - Delete a reconciliation by ID
* [alegra reconciliations export](alegra_reconciliations_export.md)	 - Export all reconciliations to CSV or JSON
* [alegra reconciliations get](alegra_reconciliations_get.md)	 - Get a single reconciliation by ID
* [alegra reconciliations import](alegra_reconciliations_import.md)	 - Bulk-create reconciliations from a CSV file
* [alegra reconciliations list](alegra_reconciliations_list.md)	 - List reconciliations

