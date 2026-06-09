---
title: alegra config
---

## alegra config

Manage alegra-cli configuration and profiles

### Options

```
  -h, --help   help for config
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
* [alegra config list-profiles](alegra_config_list-profiles.md)	 - List configured profiles
* [alegra config path](alegra_config_path.md)	 - Print the config file path
* [alegra config set-country](alegra_config_set-country.md)	 - Set an offline fallback country hint for pre-flight validation
* [alegra config set-profile](alegra_config_set-profile.md)	 - Create or update a profile (use `alegra auth login` to set the token)
* [alegra config use](alegra_config_use.md)	 - Set the default profile
* [alegra config view](alegra_config_view.md)	 - Show the current configuration (tokens redacted)

