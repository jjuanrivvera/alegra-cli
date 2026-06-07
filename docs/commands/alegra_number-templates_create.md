---
title: alegra number-templates create
---

## alegra number-templates create

Create a number-template

### Synopsis

Create a number-template.

Provide the body with --file <path> (recommended for nested documents),
--data '<json>', or one or more --set key=value pairs for flat fields.
The body is pre-flight validated for your country; use --no-validate to skip.

```
alegra number-templates create [flags]
```

### Examples

```
  alegra number-templates create -f number-template.json
  alegra number-templates create --set name="Example"
  echo '{...}' | alegra number-templates create -f -
```

### Options

```
      --country string    Country for pre-flight validation (default: auto-detected from the account)
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

* [alegra number-templates](alegra_number-templates.md)	 - Manage document numberings (numeraciones de facturación)

