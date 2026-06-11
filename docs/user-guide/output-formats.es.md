---
title: Formatos de salida
---

# Formatos de salida

Elige con `-o`/`--output` (o con `ALEGRA_OUTPUT`, o con `settings.defaultOutputFormat`).

| Formato | Uso |
| --- | --- |
| `table` | Legible para humanos (por defecto). Una lista se muestra como una cuadrícula; un registro único como clave/valor. |
| `json` | JSON con formato legible: ideal para `jq`. |
| `yaml` | YAML. |
| `csv` | Compatible con hojas de cálculo; una fila por registro. |

```bash
alegra invoices list -o json
alegra items list -o csv > items.csv
alegra contacts get 5 -o yaml
```

## Elegir columnas

`--columns` sobrescribe las columnas por defecto (por clave JSON) para `table` y `csv`:

```bash
alegra contacts list --columns id,name,email,status
alegra items list -o csv --columns id,name,price
```

Cada recurso trae columnas por defecto sensatas; `json`/`yaml` siempre incluyen
el registro completo.

## Consejos para scripting

```bash
# Count open invoices — uses the API's metadata total, no full fetch
alegra invoices list --status open --count

# Count locally after fetching every page (when you need the records anyway)
alegra invoices list --status open --all -o json | jq 'length'

# Extract IDs
alegra contacts list --all -o json | jq -r '.[].id'
```

Prefiere `--count` en vez de `--all | jq 'length'` cuando solo necesitas el
número: le pide a Alegra el total en lugar de descargar todas las páginas.
