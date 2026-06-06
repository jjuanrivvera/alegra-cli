---
title: alegra alias
---

## alegra alias

Save and manage command aliases

### Synopsis

Aliases are shortcuts for longer commands, stored in your config.

  alegra alias set unpaid "invoices list --status open --all"
  alegra unpaid --client-id 12     # expands, then appends your extra args

Aliases never shadow built-in commands.

### Options

```
  -h, --help   help for alias
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

* [alegra](alegra.md)	 - Alegra accounting system CLI
* [alegra alias list](alegra_alias_list.md)	 - List aliases
* [alegra alias remove](alegra_alias_remove.md)	 - Remove an alias
* [alegra alias set](alegra_alias_set.md)	 - Create or update an alias

