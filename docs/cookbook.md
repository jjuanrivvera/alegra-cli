---
title: Cookbook
---

# Cookbook

Copy-paste recipes for real accounting tasks. Every command here works against
today's CLI. They assume you're authenticated (`alegra auth login` or
`ALEGRA_EMAIL`/`ALEGRA_TOKEN`).

!!! tip "Two habits that pay off"
    - Add `--dry-run` to any create/update/action to see the exact request (and a
      `curl`) without sending it.
    - Pipe `-o json` into [`jq`](https://jqlang.github.io/jq/) for scripting.

## Contacts

```bash
# List clients only, every page
alegra contacts list --type client --all

# Find a contact by email or name (free-text search)
alegra contacts list -q "acme" --columns id,name,identification,email

# Create a Colombian company (DIAN needs an identification object)
alegra contacts create \
  --set name="Acme S.A.S" \
  --set 'identification={"type":"NIT","number":"901123456"}' \
  --set 'type=["client"]' \
  --set kindOfPerson="LEGAL_ENTITY"

# Create a person paying with a cédula
alegra contacts create --set name="María Pérez" \
  --set 'identification={"type":"CC","number":"1020304050"}' \
  --set 'type=["client"]'

# Update just one field (PATCH semantics — only what you send changes)
alegra contacts update 12 --set email="cuentas@acme.com"

# How many contacts do I have? (uses the API total, no full download)
alegra contacts list --count
```

## Items (products & services)

```bash
# Create a service priced at 50,000 with 19% IVA (tax id 3 in Colombia)
alegra items create --set name="Consultoría" --set price=50000 \
  --set 'tax=[{"id":3}]'

# Export the full catalog to a spreadsheet
alegra items list --all -o csv --columns id,name,reference,price,status > items.csv

# A product's current price across price lists
alegra items get 5 -o json | jq '.price'
```

## Invoices

```bash
# This month's open (unpaid) invoices
alegra invoices list --status open \
  --date-after 2026-06-01 --date-before 2026-06-30

# Count overdue-by-date receivables
alegra invoices list --status open --due-before "$(date -u +%F)" --count

# One invoice as JSON, pull the total and balance
alegra invoices get 1234 -o json | jq '{number, total, balance, status}'

# Create a draft invoice from a file (recommended for line items)
alegra invoices create -f invoice.json

# Email an invoice to the client
alegra invoices email 1234 --set 'emails=["cuentas@acme.com"]'
```

A minimal `invoice.json` (Colombia, electronic):

```json
{
  "client": { "id": 12 },
  "numberTemplate": { "id": 7 },
  "date": "2026-06-06",
  "dueDate": "2026-06-21",
  "paymentForm": "CASH",
  "items": [
    { "id": 5, "price": 50000, "quantity": 1, "tax": [{ "id": 3 }] }
  ],
  "stamp": { "generateStamp": true }
}
```

> `stamp.generateStamp: true` is what actually emits the DIAN factura
> electrónica and returns a CUFE. Drop it (or use a draft) to keep it internal.
> See the [Electronic Invoicing guide](guides/electronic-invoicing.md).

## Payments

```bash
# Income payments received this month
alegra payments list --type in --date-after 2026-06-01

# Register a cash income payment (allocate to invoices in the body)
alegra payments create -f payment.json

# Void a payment (keeps the record, removes accounting effect)
alegra payments void 88
```

## Taxes, terms, price lists, banks

```bash
alegra taxes list --columns id,name,percentage,type     # find IVA/retención ids
alegra terms list                                        # payment terms (días)
alegra bank-accounts list
alegra number-templates list -o json | jq '.[] | {id, name, documentType}'  # resolutions
```

## Reports

```bash
alegra reports sales-by-client --from 2026-01-01 --to 2026-03-31
alegra reports sales-by-seller --from 2026-01-01 --to 2026-12-31 -o csv > sellers.csv
```

## Filtering power moves

```bash
# Any Alegra query param the CLI doesn't expose as a flag (escape hatch)
alegra contacts list --param identification=901123456 --param order_field=name

# Sort + paginate explicitly
alegra invoices list --order-field date --order-direction DESC --limit 30 --start 30

# Just the count of anything
alegra bills list --status open --count
```

## Output & scripting

```bash
# Pretty table (default), JSON, YAML, or CSV
alegra contacts list -o table
alegra contacts list -o json | jq -r '.[].email'
alegra items list --all -o csv > items.csv

# Pick exactly the columns you want
alegra invoices list --columns number,date,status,total,balance

# Save a request you run constantly as a shell function
unpaid() { alegra invoices list --status open --all -o json; }
unpaid | jq 'map(.balance | tonumber) | add'   # total outstanding
```

Next: the [end-to-end guides](guides/invoice-to-cash.md) string these into real
workflows.
