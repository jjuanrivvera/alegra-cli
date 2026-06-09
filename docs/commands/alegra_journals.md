---
title: alegra journals
---

## alegra journals

Manage accounting journal entries (comprobantes contables)

### Synopsis

Manage manual accounting journal entries (comprobantes contables) — direct debit/credit postings to ledger accounts. Use `alegra journals balance` for balances grouped by period.

### Options

```
  -h, --help   help for journals
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
* [alegra journals balance](alegra_journals_balance.md)	 - Retrieve journal balances grouped by month or day
* [alegra journals create](alegra_journals_create.md)	 - Create a journal
* [alegra journals delete](alegra_journals_delete.md)	 - Delete a journal by ID
* [alegra journals export](alegra_journals_export.md)	 - Export all journals to CSV or JSON
* [alegra journals get](alegra_journals_get.md)	 - Get a single journal by ID
* [alegra journals import](alegra_journals_import.md)	 - Bulk-create journals from a CSV file
* [alegra journals list](alegra_journals_list.md)	 - List journals
* [alegra journals update](alegra_journals_update.md)	 - Update a journal by ID

