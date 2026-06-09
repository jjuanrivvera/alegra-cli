---
title: Configuration & Profiles
---

# Configuration & Profiles

The config file lives at `~/.alegra-cli/config.yaml` (override with
`ALEGRA_CONFIG`). API tokens are **not** stored here — they live in the OS
keyring.

```yaml
defaultProfile: prod
profiles:
  prod:
    email: you@biz.com
    baseUrl: https://api.alegra.com/api/v1
  sandbox:
    email: dev@biz.com
settings:
  defaultOutputFormat: table
  requestsPerSecond: 5
  logLevel: info
```

## Managing profiles

```bash
alegra config set-profile --name prod --email you@biz.com
alegra auth login --profile prod        # store the prod token
alegra config use prod                  # set the default profile
alegra config list-profiles
alegra config view                      # tokens shown as (keyring), never printed
alegra config path
```

## Selecting a profile per command

```bash
alegra --profile sandbox invoices list
ALEGRA_PROFILE=sandbox alegra invoices list
```

## Precedence

For every setting, the order is: **command-line flag → environment variable →
profile in config → built-in default**. Credentials additionally fall back to
the OS keyring when not provided by env or config.

## Environment variables

| Variable | Meaning |
| --- | --- |
| `ALEGRA_EMAIL` | Basic auth user |
| `ALEGRA_TOKEN` | Basic auth password (API token) |
| `ALEGRA_BEARER_TOKEN` | OAuth bearer token |
| `ALEGRA_BASE_URL` | API base URL |
| `ALEGRA_PROFILE` | Active profile |
| `ALEGRA_OUTPUT` | Default output format |
| `ALEGRA_LOG_LEVEL` | Log level (debug, info, warn, error) |
| `ALEGRA_CONFIG` | Config file path |

## Country & pre-flight validation

Alegra is one API localized per country. `alegra init` and `alegra auth login`
**auto-detect** your country from `/company` (`applicationVersion`: colombia,
mexico, costaRica, peru, …) and cache it on the profile, so country-specific
behavior just works.

`create` runs a country-aware **pre-flight validation** (CO/MX/PE/CR) that catches
common mistakes before sending. The country is resolved in this order:

1. the `--country` flag on the command,
2. the profile's auto-detected country,
3. the offline hint set with `alegra config set-country <country>`.

```bash
alegra config set-country colombia   # offline hint when no account is detected
alegra invoices create -f inv.json --no-validate   # skip the pre-flight check
```

The same per-country reference data (units, identification types, tax types, …)
is available offline via [`alegra catalog`](catalog.md).
