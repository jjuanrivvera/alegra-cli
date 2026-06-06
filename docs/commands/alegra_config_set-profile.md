---
title: alegra config set-profile
---

## alegra config set-profile

Create or update a profile (use `alegra auth login` to set the token)

```
alegra config set-profile [flags]
```

### Options

```
      --base-url string   API base URL
      --email string      Account email
  -h, --help              help for set-profile
      --name string       Profile name (default: active profile)
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

* [alegra config](alegra_config.md)	 - Manage alegra-cli configuration and profiles

