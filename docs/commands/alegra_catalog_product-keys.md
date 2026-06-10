---
title: alegra catalog product-keys
---

## alegra catalog product-keys

Search SAT product/service keys (México, claveProdServ)

### Synopsis

Search México's SAT c_ClaveProdServ catalog offline. Matching is case- and
accent-insensitive across the key, its description, and the SAT's published
similar-names list. Requires a one-time `alegra catalog sync-sat`.

```
alegra catalog product-keys [query] [flags]
```

### Examples

```
  alegra catalog product-keys refrigerador
  alegra catalog product-keys 10101506 -o json
  alegra catalog product-keys "servicios de programacion" --limit 5
```

### Options

```
  -h, --help        help for product-keys
      --limit int   Maximum results to show (0 = no limit) (default 25)
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

