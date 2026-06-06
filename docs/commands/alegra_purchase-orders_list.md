---
title: alegra purchase-orders list
---

## alegra purchase-orders list

List purchase-orders

```
alegra purchase-orders list [flags]
```

### Examples

```
  alegra purchase-orders list
  alegra purchase-orders list --limit 30 --all -o json
  alegra purchase-orders list --count
  alegra purchase-orders list --status <value>
  alegra purchase-orders list --param <api_param>=<value>
```

### Options

```
      --all                      Fetch all pages
      --client-id string         Filter by client ID
      --cost-center-id string    Filter by cost center ID (null/none/0 for none)
      --count                    Print only the total number of matching records
      --currency string          Filter by currency code
      --date string              Filter by date (YYYY-MM-DD)
      --delivery-date string     Filter by delivery date (YYYY-MM-DD)
  -h, --help                     help for list
      --limit int                Max records per page (max 30)
      --number string            Filter by number
      --order-direction string   Sort direction: ASC or DESC
      --order-field string       Field to sort by (id, name, date, deliveryDate, status)
      --param stringArray        Arbitrary API query parameter: key=value (repeatable; e.g. --param date_after=2026-01-01)
      --provider-name string     Filter by provider name
  -q, --query string             Free-text search
      --start int                Offset to start from (pagination)
      --status string            Filter by status
      --warehouse-id string      Filter by warehouse ID
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

* [alegra purchase-orders](alegra_purchase-orders.md)	 - Manage purchase orders (órdenes de compra)

