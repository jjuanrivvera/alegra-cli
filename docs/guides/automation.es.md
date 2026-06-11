---
title: Automatización y scripting
---

# Automatización y scripting

`alegra` está hecho para automatizarse: JSON de entrada/salida, códigos de salida
que significan algo, `--dry-run` en todas partes, y un limitador de velocidad
integrado (se adapta a los encabezados `X-Rate-Limit-*`) para que los bucles se
mantengan bajo el tope por minuto de tu plan.

## Usa un perfil + entorno para corridas no interactivas

En CI o cron no hay prompt del keyring — usa variables de entorno:

```bash
export ALEGRA_EMAIL="you@biz.com"
export ALEGRA_TOKEN="$ALEGRA_API_TOKEN"   # desde tu almacén de secretos
export ALEGRA_OUTPUT=json                  # salida por defecto apta para máquinas
alegra auth status
```

## Importa contactos en masa desde un CSV

Usa el importador integrado. Dado un `clients.csv` con un encabezado `Name,NIT`:

```bash
alegra contacts import -f clients.csv \
  --map 'Name=name,NIT=identification.number' \
  --set 'identification.type=NIT' --set 'type=["client"]' --set kindOfPerson=LEGAL_ENTITY
```

- `--map` renombra columnas del CSV a campos de la API; las rutas con punto
  construyen objetos anidados.
- `--set` aplica campos constantes a cada fila.
- Las filas son independientes — las fallas se reportan y no detienen la corrida.
- Agrega `--dry-run` para previsualizar los cuerpos JSON sin enviarlos.

Todo recurso que admita creación tiene `import`.

## Exporta cualquier cosa a una hoja de cálculo

```bash
alegra items export > items.csv                          # CSV, todas las páginas
alegra invoices export --param status=open --format json # JSON
alegra contacts export --out contacts.csv
```

`export` pagina automáticamente y usa CSV por defecto; reduce columnas con
`--columns`.

## Instantánea diaria vía cron

```cron
# 8am cada día hábil: envíate por correo una instantánea tipo P&L en una línea
0 8 * * 1-5  ALEGRA_EMAIL=you@biz.com ALEGRA_TOKEN=xxx /usr/local/bin/alegra invoices list --status open --count | mail -s "Open invoices" you@biz.com
```

## Empuja en vez de sondear: suscripciones a webhooks

El sondeo con cron es simple, pero para flujos orientados a eventos suscríbete a
los webhooks de Alegra para que *él* te llame a ti cuando algo cambie:

```bash
alegra webhook-subscriptions create -f hook.json   # { event, url, ... }
alegra webhook-subscriptions list
```

## Previsualiza en CI, nunca emitas por accidente

`--dry-run` imprime la petición exacta (y un `curl`) y no envía nada — perfecto
para validaciones de pull request que verifican payloads:

```bash
alegra invoices create -f invoice.json --dry-run
```

## Manéjalo desde un agente de IA (MCP)

El CLI completo se expone como un servidor Model Context Protocol:

```bash
claude mcp add alegra -- alegra mcp start    # o: alegra mcp claude enable
```

Luego pídele a tu agente cosas como *"lista las facturas abiertas de este mes en
Alegra y suma los saldos."* Cada comando se convierte en una herramienta
(`alegra_invoices_list`, …). Apúntalo a un perfil de sandbox/prueba mientras
experimentas. Ver [Servidor MCP](../user-guide/mcp.md).

## Códigos de salida y errores

`alegra` sale con un código distinto de cero al fallar e imprime una sola línea
`Error: …` a stderr. Para el significado de los estados HTTP (p. ej., `402` =
restringido por plan, `429` = límite de velocidad), ver la
[Referencia de errores](../reference/errors.md).

```bash
if ! alegra auth status >/dev/null 2>&1; then
  echo "Alegra auth is broken" >&2; exit 1
fi
```

## Sé un buen ciudadano de la API

- Prefiere `--count` antes que `--all | wc -l` cuando solo necesitas un número.
- Usa `--columns`/`-o csv` para traer solo lo que necesitas.
- Backfills largos: el limitador te mantiene bajo 150/min, pero considera
  correrlos fuera de las horas pico.
