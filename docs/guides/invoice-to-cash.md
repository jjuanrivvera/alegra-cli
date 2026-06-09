---
title: Invoice to Cash
---

# Invoice → Cash (a full sales cycle)

This guide walks a complete order-to-cash flow with the CLI: set up the client
and product, invoice them, send it, then record the payment.

## 1. Find or create the client

```bash
# Already there?
alegra contacts list -q "acme" --columns id,name,identification

# If not, create (Colombia example)
CID=$(alegra contacts create \
  --set name="Acme S.A.S" \
  --set 'identification={"type":"NIT","number":"901123456"}' \
  --set 'type=["client"]' --set kindOfPerson="LEGAL_ENTITY" \
  -o json | jq -r '.id')
echo "client $CID"
```

## 2. Find or create the product/service

```bash
IID=$(alegra items create --set name="Consultoría mensual" --set price=500000 \
  --set 'tax=[{"id":3}]' -o json | jq -r '.id')
```

## 3. Create the invoice

Build the body in a file (line items don't fit `--set` well):

```bash
cat > invoice.json <<JSON
{
  "client": { "id": $CID },
  "numberTemplate": { "id": 7 },
  "date": "$(date -u +%F)",
  "dueDate": "$(date -u -v+15d +%F 2>/dev/null || date -u -d '+15 days' +%F)",
  "paymentForm": "CREDIT",
  "items": [ { "id": $IID, "price": 500000, "quantity": 1, "tax": [{"id":3}] } ],
  "stamp": { "generateStamp": true }
}
JSON

INV=$(alegra invoices create -f invoice.json -o json | jq -r '.id')
alegra invoices get "$INV" -o json | jq '{number, status, total, balance}'
```

## 4. Email it to the client

```bash
alegra invoices email "$INV" --set 'emails=["cuentas@acme.com"]'
```

## 5. Record the payment when it arrives

When Acme pays, register an **income** payment and allocate it to the invoice.
Put the allocation in a file:

```bash
cat > payment.json <<JSON
{
  "type": "in",
  "date": "$(date -u +%F)",
  "bankAccount": { "id": 2 },
  "client": { "id": $CID },
  "invoices": [ { "id": $INV, "amount": 595000 } ],
  "paymentMethod": "transfer"
}
JSON

alegra payments create -f payment.json
alegra invoices get "$INV" -o json | jq '{status, balance}'   # balance should be 0, status closed
```

## 6. Need to reverse it?

Don't edit a stamped invoice. Depending on what you need:

```bash
# Annul (void) or reopen the document itself
alegra invoices void "$INV"
alegra invoices open "$INV"

# A recorded payment can be voided / reopened too
alegra payments void <payment-id>

# To correct the amount, issue a credit note
alegra credit-notes create -f credit-note.json
```

`alegra invoices preview "$INV"` returns a PDF preview URL if you want to eyeball
it first.

---

### Tips

- Use `--dry-run` on steps 3–5 first to inspect each request.
- `paymentForm`: `CASH` (paid now, set `dueDate = date`) vs `CREDIT` (on terms).
- Recurring billing? Model it once and use `alegra recurring-invoices create`.
- See [Electronic Invoicing](electronic-invoicing.md) for the stamping details and
  country-specific required fields.
