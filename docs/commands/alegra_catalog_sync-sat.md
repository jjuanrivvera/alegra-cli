---
title: alegra catalog sync-sat
---

## alegra catalog sync-sat

Download the SAT product-keys catalog (México) to the local cache

### Synopsis

Download México's SAT c_ClaveProdServ catalog (~52k product/service keys,
~7MB) into the shared local cache. Alegra exposes no endpoint for it, so the
CLI sources it from the SAT's published data (phpcfdi/resources-sat-catalogs
mirror). Needed once; re-run to pick up new SAT catalog versions. No Alegra
credentials required.

```
alegra catalog sync-sat [flags]
```

### Options

```
  -h, --help   help for sync-sat
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

* [alegra catalog](alegra_catalog.md)	 - Country reference catalogs (units, identification types, taxes, ...)

