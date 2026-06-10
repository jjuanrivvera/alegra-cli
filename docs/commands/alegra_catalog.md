---
title: alegra catalog
---

## alegra catalog

Country reference catalogs (units, identification types, taxes, ...)

### Synopsis

Show Alegra's per-country reference catalogs — units of measure and reference
enums such as identification types, tax types, payment methods, and document
types. This is the data the official MCP's units/reference tools serve; Alegra
exposes no public REST endpoint for it, so the CLI embeds it (generated from
Alegra's published country parameter pages).

Run without an argument to list the categories for your country; pass a category
key to list its values. The country is auto-detected from your account; override
it with --country (works offline, no login required).

```
alegra catalog [category] [flags]
```

### Examples

```
  alegra catalog                       # list categories for your country
  alegra catalog units                 # units of measure
  alegra catalog identification-types  # valid contact ID types
  alegra catalog units --country mexico -o json
```

### Options

```
      --country string   Country to look up (default: auto-detected from the account; e.g. colombia, mexico, costaRica, peru)
  -h, --help             help for catalog
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
* [alegra catalog product-keys](alegra_catalog_product-keys.md)	 - Search SAT product/service keys (México, claveProdServ)
* [alegra catalog sync-sat](alegra_catalog_sync-sat.md)	 - Download the SAT product-keys catalog (México) to the local cache

