---
title: Expenses & Purchases
---

# Expenses & Purchases (the payables side)

The buying cycle mirrors sales: suppliers instead of clients, bills instead of
invoices, outgoing payments instead of incoming ones.

## Suppliers (provider contacts)

```bash
# Create a provider
PID=$(alegra contacts create --set name="Proveedor X S.A.S" \
  --set 'identification={"type":"NIT","number":"900888777"}' \
  --set 'type=["provider"]' --set kindOfPerson="LEGAL_ENTITY" -o json | jq -r '.id')

# List providers
alegra contacts list --type provider --all
```

## Bills (facturas de proveedor)

```bash
# Record a bill (purchases.items holds the lines)
cat > bill.json <<JSON
{
  "provider": { "id": $PID },
  "date": "$(date -u +%F)",
  "dueDate": "$(date -u +%F)",
  "purchases": { "items": [ { "id": 1, "price": 100000, "quantity": 1 } ] }
}
JSON
alegra bills create -f bill.json

# List open bills, oldest first
alegra bills list --status open --order-field date --order-direction ASC --all

# What do I owe? (count + outstanding)
alegra bills list --status open --count
alegra bills list --status open --all -o json | jq 'map(.balance|tonumber)|add'

# Colombia: import a supplier's electronic bill straight from its CUFE
alegra bills import-by-cufe --set cufe="<the-cufe>"

# Close a bill with a remaining balance (writes it off)
alegra bills close 7 --set date="$(date -u +%F)" --set 'category={"id":5}'

# Apply a provider advance to a bill; attach a PDF; comment
alegra bills advances 7 --set 'advances=[{"id":42}]'
alegra bills attach 7 --set 'file="<base64>"' --set name="factura.pdf"
alegra bills comments 7 --set text="Pagada parcialmente"

# Adjust perceptions / retentions on a bill (replace the set)
alegra bills retentions 7 --set 'retentions=[{"id":12}]'
alegra bills perceptions 7 --set 'perceptions=[]'
```

## Purchase orders

```bash
alegra purchase-orders create -f po.json
alegra purchase-orders email 3 --set 'emails=["ventas@proveedor.com"]'
alegra purchase-orders void 3
```

## Paying suppliers (outgoing payments)

```bash
cat > out.json <<JSON
{
  "type": "out",
  "date": "$(date -u +%F)",
  "bankAccount": { "id": 2 },
  "provider": { "id": $PID },
  "bills": [ { "id": 7, "amount": 100000 } ]
}
JSON
alegra payments create -f out.json

# All expenses paid this month
alegra payments list --type out --date-after 2026-06-01 --all
```

## Retentions & taxes

Colombian purchases often carry retenciones (retefuente, ICA). Look up the ids
and attach them in the bill/payment body:

```bash
alegra retentions list --columns id,name,percentage
alegra taxes list --columns id,name,percentage,type
```

See the [Reporting guide](reporting-and-month-end.md) to roll all this up at
month-end.
