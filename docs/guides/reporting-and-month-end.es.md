---
title: Reportes y cierre de mes
---

# Reportes y cierre de mes

Recetas para los números que necesitas a la hora del cierre. Se apoyan en
`--count` (totales baratos vía la API), `--all` (paginación automática) y `jq`
(cálculos locales).

## Totales rápidos

```bash
# Cuántas facturas se emitieron, cuántas siguen abiertas
alegra invoices list --count
alegra invoices list --status open --count

# Cuentas por cobrar pendientes (suma de saldos abiertos)
alegra invoices list --status open --all -o json | jq 'map(.balance|tonumber)|add'

# Cuentas por pagar pendientes
alegra bills list --status open --all -o json | jq 'map(.balance|tonumber)|add'
```

## Actividad de este mes

```bash
FROM=2026-06-01; TO=2026-06-30

# Facturas de venta de este mes
alegra invoices list --date-after $FROM --date-before $TO --all \
  -o csv --columns number,date,status,total,balance > sales-$FROM.csv

# Ingresos recibidos este mes
alegra payments list --type in --date-after $FROM --date-before $TO --all \
  -o json | jq 'map(.amount|tonumber)|add'

# Gastos pagados este mes
alegra payments list --type out --date-after $FROM --date-before $TO --all \
  -o json | jq 'map(.amount|tonumber)|add'
```

## Reportes de ventas

```bash
# Quién compró más este trimestre
alegra reports sales-by-client --from 2026-01-01 --to 2026-03-31 \
  --order-field total --order-direction DESC

# Ventas por vendedor, a CSV (gestiona a los vendedores con `alegra sellers`)
alegra reports sales-by-seller --from 2026-01-01 --to 2026-12-31 -o csv > sellers.csv
```

## Valuación de inventario al cierre

```bash
# Existencias por bodega de un ítem, a la fecha de fin de período
alegra items stock 5 --date 2026-06-30

# ¿Necesitas una unidad / impuesto / código de ID válido para un ajuste? Búscalo sin conexión:
alegra catalog units
```

!!! note "Reportes según el plan"
    `reports income-statement` y `reports account-statement` existen pero solo
    están disponibles en ciertos planes de Alegra — de lo contrario devuelven
    **HTTP 402/403**.

## Antigüedad de cartera (cuentas por cobrar vencidas), hoy

Hasta que llegue un comando dedicado `ar aging`, calcula los buckets con `jq`:

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

## Una instantánea de cierre de mes en una línea

```bash
FROM=2026-06-01; TO=2026-06-30
echo "Invoiced: $(alegra invoices list --date-after $FROM --date-before $TO --count)"
echo "Open AR:  $(alegra invoices list --status open --all -o json | jq 'map(.balance|tonumber)|add')"
echo "Open AP:  $(alegra bills list --status open --all -o json | jq 'map(.balance|tonumber)|add')"
```

Envuelve eso en un script y córrelo desde [automatización](automation.md) de forma
programada.
