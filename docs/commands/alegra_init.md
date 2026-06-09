---
title: alegra init
---

## alegra init

Guided setup: authenticate, detect your country, and save a profile

### Synopsis

init walks you through first-time setup: it takes your Alegra email and
API token, verifies them, auto-detects your account's country (the API version),
and saves a profile (token in the OS keyring).

Find your API token in Alegra: Configuración → Integraciones → API.

```
alegra init [flags]
```

### Options

```
      --base-url string   Override the API base URL
      --email string      Alegra account email (skips the prompt)
  -h, --help              help for init
      --token string      Alegra API token (skips the prompt)
```

### Options inherited from parent commands

```
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

