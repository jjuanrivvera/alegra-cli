---
title: alegra currencies
---

## alegra currencies

Manage currencies (monedas)

### Synopsis

Manage currencies. Currencies are identified by their ISO code (e.g. USD), not a numeric id; get/update take the code.

### Options

```
  -h, --help   help for currencies
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
* [alegra currencies create](alegra_currencies_create.md)	 - Create a currency
* [alegra currencies delete](alegra_currencies_delete.md)	 - Delete a currency by ID
* [alegra currencies export](alegra_currencies_export.md)	 - Export all currencies to CSV or JSON
* [alegra currencies get](alegra_currencies_get.md)	 - Get a single currency by ID
* [alegra currencies import](alegra_currencies_import.md)	 - Bulk-create currencies from a CSV file
* [alegra currencies list](alegra_currencies_list.md)	 - List currencies
* [alegra currencies update](alegra_currencies_update.md)	 - Update a currency by ID

