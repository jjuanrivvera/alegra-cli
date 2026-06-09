---
title: alegra warehouse-transfers
---

## alegra warehouse-transfers

Manage inventory transfers between warehouses

### Synopsis

Manage warehouse transfers — movements of item stock from one warehouse (bodega) to another.

### Options

```
  -h, --help   help for warehouse-transfers
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
* [alegra warehouse-transfers create](alegra_warehouse-transfers_create.md)	 - Create a warehouse-transfer
* [alegra warehouse-transfers delete](alegra_warehouse-transfers_delete.md)	 - Delete a warehouse-transfer by ID
* [alegra warehouse-transfers export](alegra_warehouse-transfers_export.md)	 - Export all warehouse-transfers to CSV or JSON
* [alegra warehouse-transfers get](alegra_warehouse-transfers_get.md)	 - Get a single warehouse-transfer by ID
* [alegra warehouse-transfers import](alegra_warehouse-transfers_import.md)	 - Bulk-create warehouse-transfers from a CSV file
* [alegra warehouse-transfers list](alegra_warehouse-transfers_list.md)	 - List warehouse-transfers
* [alegra warehouse-transfers update](alegra_warehouse-transfers_update.md)	 - Update a warehouse-transfer by ID

