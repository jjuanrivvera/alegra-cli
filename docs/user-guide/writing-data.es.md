---
title: Crear y actualizar
---

# Crear y actualizar registros

Los cuerpos de las peticiones de Alegra suelen estar muy anidados (una factura
tiene ítems de línea, impuestos, referencias de cliente, etc.), por eso los
comandos de creación/actualización aceptan un **cuerpo JSON** de tres formas.
Se pueden combinar: un cuerpo base desde `--data`/`--file`, con los overrides de
`--set` fusionados encima.

## Desde un archivo (o stdin)

```bash
alegra invoices create -f invoice.json
cat invoice.json | alegra invoices create -f -
```

## Desde una cadena JSON en línea

```bash
alegra contacts create -d '{"name":"Acme","type":["client"]}'
```

## Desde pares `--set key=value`

Los valores se parsean como JSON cuando son válidos (números, booleanos,
arreglos, objetos, null); de lo contrario se tratan como cadenas:

```bash
alegra contacts create \
  --set name="Acme S.A.S" \
  --set 'type=["client"]' \
  --set 'address={"city":"Cali"}'
```

## Actualizar

`update` envía un cuerpo parcial (Alegra aplica semántica PATCH — solo cambian
los campos que envías):

```bash
alegra invoices update 12 --set observations="Pagada por transferencia"
alegra contacts update 99 -f patch.json
```

## Eliminar

```bash
alegra contacts delete 99       # pide confirmación
alegra contacts delete 99 -y    # omite el prompt
```

## Acciones de recurso

Muchos documentos admiten acciones extra además del CRUD:

```bash
alegra invoices void 12
alegra invoices open 12
alegra invoices email 12 --set 'emails=["client@acme.com"]'
alegra invoices stamp --set 'invoices=[{"id":12},{"id":13}]'
alegra bills close 7 --set date=2026-06-05 --set 'category={"id":5}'
alegra bank-accounts transfer 3 --set 'destinationAccount={"id":4}' --set amount=100000
```

## Importación y exportación masiva (CSV/JSON)

Casi todos los recursos admiten operaciones masivas.

**Import** crea un registro por cada fila del CSV. La fila de encabezado nombra
los campos; usa `--map` para renombrar una columna a un campo de la API (las
rutas con puntos construyen objetos anidados) y `--set` para aplicar una
constante a cada fila. Las filas se procesan de forma independiente — una fila
fallida se reporta y no detiene la ejecución:

```bash
alegra contacts import -f clients.csv \
  --map 'Name=name,NIT=identification.number' \
  --set 'identification.type=NIT' --set 'type=["client"]'

alegra items import -f catalog.csv --map 'SKU=reference,Name=name'
```

Previsualiza cada fila sin escribir nada con `--dry-run`.

**Export** trae *todas* las páginas y las escribe a un archivo o a stdout, en CSV
(por defecto) o JSON, respetando `--columns` y cualquier filtro `--param`:

```bash
alegra invoices export --param status=open > receivables.csv
alegra items export --format json --out items.json
alegra contacts export --columns id,name,email,status
```

## Dry run

Previsualiza la petición HTTP exacta (y un `curl` listo para copiar y pegar) sin
enviarla:

```bash
alegra invoices create -f invoice.json --dry-run
```

Agrega `--show-token` para incluir el encabezado de autenticación en el curl
impreso (desactivado por defecto).

!!! tip "Validación previa al envío"
    `create` valida el cuerpo para tu país detectado (CO/MX/PE/CR) antes de
    enviarlo y detecta errores comunes a tiempo. Sáltatela con `--no-validate`, o
    define el país explícitamente con `--country` (consulta
    [Configuración](configuration.md)).
