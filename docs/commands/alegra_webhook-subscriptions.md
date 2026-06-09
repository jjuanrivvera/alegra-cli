---
title: alegra webhook-subscriptions
---

## alegra webhook-subscriptions

Manage webhook subscriptions (event notifications)

### Synopsis

Manage webhook subscriptions: Alegra POSTs to your URL whenever an event fires (e.g. new-invoice, edit-client, delete-item). Create one per event+URL with --set event=new-invoice --set url=https://your.app/hook.

### Options

```
  -h, --help   help for webhook-subscriptions
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
* [alegra webhook-subscriptions create](alegra_webhook-subscriptions_create.md)	 - Create a webhook-subscription
* [alegra webhook-subscriptions delete](alegra_webhook-subscriptions_delete.md)	 - Delete a webhook-subscription by ID
* [alegra webhook-subscriptions export](alegra_webhook-subscriptions_export.md)	 - Export all webhook-subscriptions to CSV or JSON
* [alegra webhook-subscriptions get](alegra_webhook-subscriptions_get.md)	 - Get a single webhook-subscription by ID
* [alegra webhook-subscriptions import](alegra_webhook-subscriptions_import.md)	 - Bulk-create webhook-subscriptions from a CSV file
* [alegra webhook-subscriptions list](alegra_webhook-subscriptions_list.md)	 - List webhook-subscriptions
* [alegra webhook-subscriptions update](alegra_webhook-subscriptions_update.md)	 - Update a webhook-subscription by ID

