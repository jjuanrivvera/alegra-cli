---
title: alegra company
---

## alegra company

View and update the account's company (empresa)

### Synopsis

Manage the singleton company registered on your Alegra account.

The Alegra API exposes the company only at GET /company and PUT /company; there
is no list or per-id access. Use "company get" to view it and "company update"
to edit it with a JSON body.

### Options

```
  -h, --help   help for company
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
* [alegra company get](alegra_company_get.md)	 - Show the company information
* [alegra company update](alegra_company_update.md)	 - Update the company information

