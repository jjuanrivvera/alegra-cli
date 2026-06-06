---
title: Reporting & Month-End
---

# Reporting & Month-End

Recipes for the numbers you need at closing time. These lean on `--count` (cheap
totals via the API), `--all` (auto-pagination), and `jq` (local math).

## Quick totals

```bash
# How many invoices were issued, how many still open
alegra invoices list --count
alegra invoices list --status open --count

# Outstanding receivables (sum of open balances)
alegra invoices list --status open --all -o json | jq 'map(.balance|tonumber)|add'

# Payables outstanding
alegra bills list --status open --all -o json | jq 'map(.balance|tonumber)|add'
```

## This month's activity

```bash
FROM=2026-06-01; TO=2026-06-30

# Sales invoices this month
alegra invoices list --date-after $FROM --date-before $TO --all \
  -o csv --columns number,date,status,total,balance > sales-$FROM.csv

# Income received this month
alegra payments list --type in --date-after $FROM --date-before $TO --all \
  -o json | jq 'map(.amount|tonumber)|add'

# Expenses paid this month
alegra payments list --type out --date-after $FROM --date-before $TO --all \
  -o json | jq 'map(.amount|tonumber)|add'
```

## Sales reports

```bash
# Who bought the most this quarter
alegra reports sales-by-client --from 2026-01-01 --to 2026-03-31 \
  --order-field total --order-direction DESC

# Sales by seller, to CSV
alegra reports sales-by-seller --from 2026-01-01 --to 2026-12-31 -o csv > sellers.csv
```

!!! note "Plan-gated reports"
    `reports income-statement` and `reports account-statement` exist but are only
    available on certain Alegra plans — they return **HTTP 402/403** otherwise.

## Aging (overdue receivables), today

Until a dedicated `ar aging` command lands, compute buckets with `jq`:

```bash
TODAY=$(date -u +%F)
alegra invoices list --status open --all -o json | jq --arg today "$TODAY" '
  map({number, dueDate, balance: (.balance|tonumber)})
  | group_by(.dueDate < $today)
  | map({ overdue: .[0].dueDate < $today,
          count: length,
          amount: (map(.balance)|add) })
'
```

## A one-line month-end snapshot

```bash
FROM=2026-06-01; TO=2026-06-30
echo "Invoiced: $(alegra invoices list --date-after $FROM --date-before $TO --count)"
echo "Open AR:  $(alegra invoices list --status open --all -o json | jq 'map(.balance|tonumber)|add')"
echo "Open AP:  $(alegra bills list --status open --all -o json | jq 'map(.balance|tonumber)|add')"
```

Wrap that in a script and run it from [automation](automation.md) on a schedule.
