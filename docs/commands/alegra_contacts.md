---
title: alegra contacts
---

## alegra contacts

Manage contacts (clients and providers)

### Synopsis

Create, list, update, and delete Alegra contacts — your clients and providers.

Body notes: identification is an object {type, number} (e.g. {"type":"NIT","number":"901123456"}); type is an array like ["client"] or ["provider"]; kindOfPerson is LEGAL_ENTITY or PERSON_ENTITY.

### Options

```
  -h, --help   help for contacts
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
* [alegra contacts create](alegra_contacts_create.md)	 - Create a contact
* [alegra contacts delete](alegra_contacts_delete.md)	 - Delete a contact by ID
* [alegra contacts export](alegra_contacts_export.md)	 - Export all contacts to CSV or JSON
* [alegra contacts get](alegra_contacts_get.md)	 - Get a single contact by ID
* [alegra contacts import](alegra_contacts_import.md)	 - Bulk-create contacts from a CSV file
* [alegra contacts list](alegra_contacts_list.md)	 - List contacts
* [alegra contacts update](alegra_contacts_update.md)	 - Update a contact by ID

