---
title: alegra item-categories create
---

## alegra item-categories create

Create a item-category

### Synopsis

Create a item-category.

Provide the body with --file <path> (recommended for nested documents),
--data '<json>', or one or more --set key=value pairs for flat fields.
The body is pre-flight validated for your country; use --no-validate to skip.

```
alegra item-categories create [flags]
```

### Examples

```
  alegra item-categories create -f item-category.json
  alegra item-categories create --set name="Example"
  echo '{...}' | alegra item-categories create -f -
```

### Options

```
      --country string    Country for pre-flight validation (default: config/company)
  -d, --data string       Request body as a JSON string
      --draft             Create as a draft (strip any electronic stamp from the body)
  -f, --file string       Read JSON request body from a file (use - for stdin)
  -h, --help              help for create
      --no-validate       Skip client-side pre-flight validation
      --set stringArray   Set a top-level field: key=value (value parsed as JSON when valid). Repeatable. For nested documents (e.g. invoice items[]) use --file.
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

* [alegra item-categories](alegra_item-categories.md)	 - Manage item categories

