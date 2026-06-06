---
title: Productivity Features
---

# Productivity Features

Small things that make the CLI fast to live in.

## `alegra doctor` — is everything OK?

One read-only command checks config, credentials, auth, company/country, your
rate-limit budget, numbering resolutions, and plan access:

```bash
$ alegra doctor
✔ config        ~/.alegra-cli/config.yaml (profile: default)
✔ credentials   keyring token
✔ auth          Juan Felipe Rivera (admin)
✔ company       Invitas · country: colombia · regime: Responsable de IVA
✔ rate limit    149/150 remaining (resets in 41s)
✔ numbering     14 resolution(s) configured
⚠ plan          'reconciliations' returned 403 — not in your plan
```

Run it first whenever something misbehaves.

## Aliases — save the commands you repeat

```bash
alegra alias set unpaid "invoices list --status open --all"
alegra unpaid --client-id 12        # expands, then appends your extra args
alegra alias list
alegra alias remove unpaid
```

Aliases never shadow built-in commands, and trailing arguments are appended to
the expansion.

## Natural date ranges

Every `list` accepts `--since`/`--until` with friendly values:

```bash
alegra invoices list --since this-month
alegra invoices list --since last-month --until last-month
alegra payments  list --type in --since 7d            # last 7 days
alegra bills     list --since 2026-01-01 --until 2026-03-31
```

Accepted: `YYYY-MM-DD`, `today`, `yesterday`, `tomorrow`, `this-month`,
`last-month`, `this-year`, `last-year`, `this-quarter`, and `Nd`/`Nw`/`Nm`/`Ny`.

## Counts without downloading

```bash
alegra invoices list --status open --count     # uses the API total
```

## Arbitrary filters (escape hatch)

Any Alegra query parameter the CLI doesn't expose as a flag:

```bash
alegra invoices list --param client_name="Acme" --param numberTemplate_fullNumber="FE-1"
```

## Pre-flight validation

`create` checks your body against country rules before sending (see
[Electronic Invoicing](electronic-invoicing.md)). Set the country once:

```bash
alegra config set-country colombia
```

Skip a check with `--no-validate`; create internal drafts with `--draft`.

## Friendlier errors

Failures explain themselves and suggest a fix — `402` becomes "your plan doesn't
include this", `429` becomes a rate-limit hint, and stamping codes map to
remedies. See the [Error reference](../reference/errors.md).

## Drive it from an AI agent

The whole CLI is an MCP server: `claude mcp add alegra -- alegra mcp`. See
[MCP Server](../user-guide/mcp.md).
