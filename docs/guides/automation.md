---
title: Automation & Scripting
---

# Automation & Scripting

`alegra` is built to be scripted: JSON in/out, exit codes that mean something,
`--dry-run` everywhere, and a built-in rate limiter so loops stay under Alegra's
**150 requests/minute** ceiling.

## Use a profile + env for non-interactive runs

In CI or cron there's no keyring prompt — use environment variables:

```bash
export ALEGRA_EMAIL="you@biz.com"
export ALEGRA_TOKEN="$ALEGRA_API_TOKEN"   # from your secret store
export ALEGRA_OUTPUT=json                  # machine-friendly default
alegra auth status
```

## Bulk import contacts from a CSV

Use the built-in importer. Given `clients.csv` with a `Name,NIT` header:

```bash
alegra contacts import -f clients.csv \
  --map 'Name=name,NIT=identification.number' \
  --set 'identification.type=NIT' --set 'type=["client"]' --set kindOfPerson=LEGAL_ENTITY
```

- `--map` renames CSV columns to API fields; dotted paths build nested objects.
- `--set` applies constant fields to every row.
- Rows are independent — failures are reported and don't stop the run.
- Add `--dry-run` to preview the JSON bodies without sending.

Every resource that supports creation has `import`.

## Export anything to a spreadsheet

```bash
alegra items export > items.csv                          # CSV, all pages
alegra invoices export --param status=open --format json # JSON
alegra contacts export --out contacts.csv
```

`export` auto-paginates and defaults to CSV; narrow columns with `--columns`.

## Daily snapshot via cron

```cron
# 8am every weekday: email yourself a one-line P&L-ish snapshot
0 8 * * 1-5  ALEGRA_EMAIL=you@biz.com ALEGRA_TOKEN=xxx /usr/local/bin/alegra invoices list --status open --count | mail -s "Open invoices" you@biz.com
```

## Preview in CI, never emit by accident

`--dry-run` prints the exact request (and a `curl`) and sends nothing — perfect
for pull-request checks that validate payloads:

```bash
alegra invoices create -f invoice.json --dry-run
```

## Drive it from an AI agent (MCP)

The whole CLI is exposed as a Model Context Protocol server:

```bash
claude mcp add alegra -- alegra mcp
```

Then ask your agent things like *"list this month's open invoices in Alegra and
total the balances."* Each command becomes a tool (`alegra_invoices_list`, …).
Point it at a sandbox/test profile while you experiment. See
[MCP Server](../user-guide/mcp.md).

## Exit codes & errors

`alegra` exits non-zero on failure and prints a single `Error: …` line to stderr.
For the meaning of HTTP statuses (e.g. `402` = plan-gated, `429` = rate limited),
see the [Error reference](../reference/errors.md).

```bash
if ! alegra auth status >/dev/null 2>&1; then
  echo "Alegra auth is broken" >&2; exit 1
fi
```

## Be a good API citizen

- Prefer `--count` over `--all | wc -l` when you only need a number.
- Use `--columns`/`-o csv` to fetch only what you need.
- Long backfills: the limiter keeps you under 150/min, but consider running them
  off-peak.
