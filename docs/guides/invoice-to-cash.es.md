---
title: De la Factura al Cobro
---

# Factura → Cobro (un ciclo de ventas completo)

Esta guía recorre un flujo completo de pedido a cobro con la CLI: configura el
cliente y el producto, factúrale, envíasela y luego registra el pago.

## 1. Encuentra o crea el cliente

```bash
# ¿Ya existe?
alegra contacts list -q "acme" --columns id,name,identification

# Si no, créalo (ejemplo de Colombia)
CID=$(alegra contacts create \
  --set name="Acme S.A.S" \
  --set 'identification={"type":"NIT","number":"901123456"}' \
  --set 'type=["client"]' --set kindOfPerson="LEGAL_ENTITY" \
  -o json | jq -r '.id')
echo "client $CID"
```

## 2. Encuentra o crea el producto/servicio

```bash
IID=$(alegra items create --set name="Consultoría mensual" --set price=500000 \
  --set 'tax=[{"id":3}]' -o json | jq -r '.id')
```

## 3. Crea la factura

Arma el cuerpo en un archivo (las líneas no encajan bien con `--set`):

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

## 4. Envíasela al cliente por correo

```bash
alegra invoices email "$INV" --set 'emails=["cuentas@acme.com"]'
```

## 5. Registra el pago cuando llegue

Cuando Acme pague, registra un pago de **ingreso** y asígnalo a la factura. Pon
la asignación en un archivo:

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
alegra invoices get "$INV" -o json | jq '{status, balance}'   # el saldo debería ser 0, estado closed
```

## 6. ¿Necesitas revertirla?

No edites una factura ya timbrada. Según lo que necesites:

```bash
# Anular (void) o reabrir el documento mismo
alegra invoices void "$INV"
alegra invoices open "$INV"

# Un pago registrado también se puede anular / reabrir
alegra payments void <payment-id>

# Para corregir el monto, emite una nota de crédito
alegra credit-notes create -f credit-note.json
```

`alegra invoices preview "$INV"` devuelve una URL de vista previa en PDF si
quieres revisarla primero con tus propios ojos.

---

### Tips

- Usa `--dry-run` primero en los pasos 3–5 para inspeccionar cada petición.
- `paymentForm`: `CASH` (pagada ahora, define `dueDate = date`) vs `CREDIT` (a plazo).
- ¿Facturación recurrente? Modélala una vez y usa `alegra recurring-invoices create`.
- Consulta [Facturación Electrónica](electronic-invoicing.md) para los detalles del
  timbrado y los campos obligatorios por país.
