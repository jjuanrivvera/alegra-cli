---
title: alegra auth login
---

## alegra auth login

Store Alegra credentials for a profile (token kept in the OS keyring)

### Synopsis

Authenticate against the Alegra API.

Your email and the chosen base URL are saved to the config file; the API token
is stored securely in the OS keyring (not written to disk in plaintext).

Find your API token in Alegra: Configuración → Integraciones → API.

```
alegra auth login [flags]
```

### Options

```
      --base-url string   Override the API base URL for this profile
      --email string      Alegra account email
  -h, --help              help for login
      --token string      Alegra API token (omit to be prompted securely)
```

### Options inherited from parent commands

```
      --columns strings             Comma-separated columns for table/csv output
      --dry-run                     Print the equivalent curl request without sending it
  -o, --output string               Output format: table, json, yaml, csv (env: ALEGRA_OUTPUT)
      --profile string              Configuration profile to use (env: ALEGRA_PROFILE)
      --requests-per-second float   Client-side rate limit (default from config)
      --show-token                  In --dry-run, do not redact the Authorization header
  -v, --verbose                     Enable verbose (debug) logging to stderr
```

### SEE ALSO

* [alegra auth](alegra_auth.md)	 - Manage Alegra API authentication

