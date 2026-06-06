---
title: alegra additional-charges create
---

## alegra additional-charges create

Create a additional-charge

### Synopsis

Create a additional-charge.

Provide the body with --file <path> (recommended for nested documents),
--data '<json>', or one or more --set key=value pairs for flat fields.

```
alegra additional-charges create [flags]
```

### Examples

```
  alegra additional-charges create -f additional-charge.json
  alegra additional-charges create --set name="Example"
  echo '{...}' | alegra additional-charges create -f -
```

### Options

```
  -d, --data string       Request body as a JSON string
  -f, --file string       Read JSON request body from a file (use - for stdin)
  -h, --help              help for create
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

* [alegra additional-charges](alegra_additional-charges.md)	 - Manage additional charges (tips and parafiscal contributions)

