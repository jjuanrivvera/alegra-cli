---
title: alegra invoices emit
---

## alegra invoices emit

Emit (stamp) draft/open invoices electronically, in batches of 10

### Synopsis

emit sends invoices to the tax authority for stamping.

Pass invoice ids explicitly, or use --all to emit every draft invoice. Already-
emitted ids are skipped (a local idempotency guard) unless --force is given —
emission is NOT idempotent on Alegra's side, so this prevents duplicates.

```
alegra invoices emit [id...] [flags]
```

### Examples

```
  alegra invoices emit 1234 1235
  alegra invoices emit --all            # every draft
  alegra invoices emit --all --dry-run
```

### Options

```
      --all     Emit every draft invoice
      --force   Re-emit even if locally recorded as emitted
  -h, --help    help for emit
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

* [alegra invoices](alegra_invoices.md)	 - Manage sales invoices

