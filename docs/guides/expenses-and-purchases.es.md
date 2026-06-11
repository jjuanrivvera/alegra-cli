---
title: Gastos y Compras
---

# Gastos y Compras (el lado de las cuentas por pagar)

El ciclo de compra refleja el de ventas: proveedores en vez de clientes,
facturas de compra en vez de facturas de venta, pagos salientes en vez de
entrantes.

## Proveedores (contactos de tipo proveedor)

```bash
# Crear un proveedor
PID=$(alegra contacts create --set name="Proveedor X S.A.S" \
  --set 'identification={"type":"NIT","number":"900888777"}' \
  --set 'type=["provider"]' --set kindOfPerson="LEGAL_ENTITY" -o json | jq -r '.id')

# Listar proveedores
alegra contacts list --type provider --all
```

## Facturas de compra (facturas de proveedor)

```bash
# Registrar una factura de compra (purchases.items contiene las líneas)
cat > bill.json <<JSON
{
  "provider": { "id": $PID },
  "date": "$(date -u +%F)",
  "dueDate": "$(date -u +%F)",
  "purchases": { "items": [ { "id": 1, "price": 100000, "quantity": 1 } ] }
}
JSON
alegra bills create -f bill.json

# Listar facturas de compra abiertas, de la más antigua a la más reciente
alegra bills list --status open --order-field date --order-direction ASC --all

# ¿Cuánto debo? (conteo + saldo pendiente)
alegra bills list --status open --count
alegra bills list --status open --all -o json | jq 'map(.balance|tonumber)|add'

# Colombia: importar la factura electrónica de un proveedor directo desde su CUFE
alegra bills import-by-cufe --set cufe="<the-cufe>"

# Cerrar una factura de compra con saldo restante (la da de baja)
alegra bills close 7 --set date="$(date -u +%F)" --set 'category={"id":5}'

# Aplicar un anticipo de proveedor a una factura de compra; adjuntar un PDF; comentar
alegra bills advances 7 --set 'advances=[{"id":42}]'
alegra bills attach 7 --set 'file="<base64>"' --set name="factura.pdf"
alegra bills comments 7 --set text="Pagada parcialmente"

# Ajustar percepciones / retenciones en una factura de compra (reemplaza el conjunto)
alegra bills retentions 7 --set 'retentions=[{"id":12}]'
alegra bills perceptions 7 --set 'perceptions=[]'
```

## Órdenes de compra

```bash
alegra purchase-orders create -f po.json
alegra purchase-orders email 3 --set 'emails=["ventas@proveedor.com"]'
alegra purchase-orders void 3
```

## Pagar a proveedores (pagos salientes)

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

# Todos los gastos pagados este mes
alegra payments list --type out --date-after 2026-06-01 --all
```

## Retenciones e impuestos

Las compras colombianas suelen llevar retenciones (retefuente, ICA). Busca los
ids y adjúntalos en el cuerpo de la factura de compra/pago:

```bash
alegra retentions list --columns id,name,percentage
alegra taxes list --columns id,name,percentage,type
```

Consulta la [guía de Reportes](reporting-and-month-end.md) para consolidar todo
esto al cierre de mes.
