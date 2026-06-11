---
title: Recetario
---

# Recetario

Recetas listas para copiar y pegar para tareas contables reales. Cada comando de
aquí funciona con la CLI actual. Asumen que ya estás autenticado (`alegra auth
login` o `ALEGRA_EMAIL`/`ALEGRA_TOKEN`).

!!! tip "Dos hábitos que valen la pena"
    - Agrega `--dry-run` a cualquier creación, actualización o acción para ver la
      petición exacta (y un `curl`) sin enviarla.
    - Pasa `-o json` por un pipe a [`jq`](https://jqlang.github.io/jq/) para
      hacer scripting.

## Contactos

```bash
# Listar solo clientes, todas las páginas
alegra contacts list --type client --all

# Buscar un contacto por correo o nombre (búsqueda de texto libre)
alegra contacts list -q "acme" --columns id,name,identification,email

# Crear una empresa colombiana (la DIAN necesita un objeto identification)
alegra contacts create \
  --set name="Acme S.A.S" \
  --set 'identification={"type":"NIT","number":"901123456"}' \
  --set 'type=["client"]' \
  --set kindOfPerson="LEGAL_ENTITY"

# Crear una persona que paga con cédula
alegra contacts create --set name="María Pérez" \
  --set 'identification={"type":"CC","number":"1020304050"}' \
  --set 'type=["client"]'

# Actualizar un solo campo (semántica PATCH — solo cambia lo que envías)
alegra contacts update 12 --set email="cuentas@acme.com"

# ¿Cuántos contactos tengo? (usa el total de la API, sin descarga completa)
alegra contacts list --count
```

## Ítems (productos y servicios)

```bash
# Crear un servicio con precio 50,000 e IVA del 19% (id de impuesto 3 en Colombia)
alegra items create --set name="Consultoría" --set price=50000 \
  --set 'tax=[{"id":3}]'

# Exportar el catálogo completo a una hoja de cálculo
alegra items list --all -o csv --columns id,name,reference,price,status > items.csv

# Las entradas de precio de un producto (una fila por lista de precios)
alegra items get 5 -o json | jq '.price'

# Stock por bodega de un ítem (--date para una foto histórica)
alegra items stock 5
alegra items stock 5 --date 2026-03-31

# Crear ítems en lote desde un CSV (una fila por ítem; los errores no detienen la corrida)
alegra items import -f catalog.csv --map 'SKU=reference,Name=name,Price=price'
```

## Facturas

```bash
# Facturas abiertas (sin pagar) de este mes — rango de fechas natural
alegra invoices list --status open --since this-month

# …o fechas explícitas
alegra invoices list --status open \
  --date-after 2026-06-01 --date-before 2026-06-30

# Contar cuentas por cobrar vencidas por fecha
alegra invoices list --status open --due-before "$(date -u +%F)" --count

# Una factura como JSON, extrae el total y el saldo
alegra invoices get 1234 -o json | jq '{number, total, balance, status}'

# Crear una factura en borrador desde un archivo (recomendado para las líneas)
alegra invoices create -f invoice.json

# Enviar una factura por correo al cliente
alegra invoices email 1234 --set 'emails=["cuentas@acme.com"]'
```

Un `invoice.json` mínimo (Colombia, electrónica):

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

> `stamp.generateStamp: true` es lo que realmente emite la factura electrónica de
> la DIAN y devuelve un CUFE. Quítalo (o usa un borrador) para mantenerla interna.
> Consulta la [guía de Facturación Electrónica](guides/electronic-invoicing.md).

## Pagos

```bash
# Pagos de ingreso recibidos este mes
alegra payments list --type in --date-after 2026-06-01

# Registrar un pago de ingreso en efectivo (asígnalo a facturas en el cuerpo)
alegra payments create -f payment.json

# Anular un pago (conserva el registro, elimina el efecto contable)
alegra payments void 88
```

## Facturas de compra y gastos

```bash
# Facturas de proveedor aún abiertas este mes
alegra bills list --status open --since this-month

# Registrar una factura de proveedor desde un archivo
alegra bills create -f bill.json

# Importar una factura electrónica colombiana recibida por su CUFE
alegra bills import-by-cufe --set cufe="<CUFE>"

# Aplicar un anticipo de proveedor, adjuntar el PDF, agregar un comentario
alegra bills advances 7 --set 'advances=[{"id":42}]'
alegra bills attach 7 --set 'file="<base64>"' --set name="factura.pdf"
alegra bills comments 7 --set text="Pagada por transferencia"

# Órdenes de compra: enviar por correo, luego anular
alegra purchase-orders email 3 --set 'emails=["proveedor@acme.com"]'
alegra purchase-orders void 3
```

## Vendedores, webhooks y datos de referencia

```bash
# Vendedores — se agregan en `reports sales-by-seller`
alegra sellers list

# Suscribirse a eventos de webhook (push, en vez de hacer polling)
alegra webhook-subscriptions create -f hook.json

# Catálogos offline por país — evita codificar valores mágicos (id de impuesto 3, NIT, CASH)
alegra catalog                                  # listar categorías de tu país
alegra catalog units
alegra catalog identification-types --country mexico
```

## Impuestos, plazos, listas de precios, bancos

```bash
alegra taxes list --columns id,name,percentage,type     # encontrar ids de IVA/retención
alegra terms list                                        # plazos de pago (días)
alegra bank-accounts list
alegra number-templates list -o json | jq '.[] | {id, name, documentType}'  # resoluciones
```

## Reportes

```bash
alegra reports sales-by-client --from 2026-01-01 --to 2026-03-31
alegra reports sales-by-seller --from 2026-01-01 --to 2026-12-31 -o csv > sellers.csv
```

## Jugadas avanzadas de filtrado

```bash
# Cualquier parámetro de consulta de Alegra que la CLI no exponga como flag (escotilla de escape)
alegra contacts list --param identification=901123456 --param order_field=name

# Ordenar + paginar explícitamente
alegra invoices list --order-field date --order-direction DESC --limit 30 --start 30

# Solo el conteo de lo que sea
alegra bills list --status open --count
```

## Salida y scripting

```bash
# Tabla bonita (por defecto), JSON, YAML o CSV
alegra contacts list -o table
alegra contacts list -o json | jq -r '.[].email'
alegra items list --all -o csv > items.csv

# Elige exactamente las columnas que quieres
alegra invoices list --columns number,date,status,total,balance

# Guarda una petición que corres constantemente como una función de shell
unpaid() { alegra invoices list --status open --all -o json; }
unpaid | jq 'map(.balance | tonumber) | add'   # total pendiente
```

Siguiente: las [guías de extremo a extremo](guides/invoice-to-cash.md) encadenan
todo esto en flujos de trabajo reales.
