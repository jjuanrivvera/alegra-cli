---
title: alegra config set-country
---

## alegra config set-country

Set an offline fallback country hint for pre-flight validation

### Synopsis

Set a global, offline fallback country hint for client-side pre-flight
validation (e.g. colombia, mexico).

The account's real country is the source of truth and is auto-detected per
profile on 'auth login' (refreshed by 'doctor'); that detected value takes
precedence over this hint. Use set-country only when you need validation to know
the country before logging in (e.g. offline '--dry-run').

```
alegra config set-country <country> [flags]
```

### Options

```
  -h, --help   help for set-country
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

* [alegra config](alegra_config.md)	 - Manage alegra-cli configuration and profiles

