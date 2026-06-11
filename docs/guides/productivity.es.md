---
title: Funciones de productividad
---

# Funciones de productividad

Pequeños detalles que hacen que el CLI sea ágil para el día a día.

## `alegra doctor` — ¿está todo bien?

Un solo comando de solo lectura revisa la configuración, las credenciales, la
autenticación, la empresa/país, tu presupuesto de límite de velocidad, las
resoluciones de numeración y el acceso según el plan:

```bash
$ alegra doctor
✔ config        ~/.alegra-cli/config.yaml (profile: default)
✔ credentials   keyring token
✔ auth          Juan Felipe Rivera (admin)
✔ company       Invitas · country: colombia · regime: Responsable de IVA
✔ rate limit    149/150 remaining (resets in 41s)
✔ numbering     14 resolution(s) configured
✔ plan          core endpoints reachable
```

(Una línea `⚠ plan` aparece solo cuando un endpoint que tu plan no incluye
devuelve `402` — ver la [Referencia de errores](../reference/errors.md).)

Córrelo primero cada vez que algo se comporte mal.

## Alias — guarda los comandos que repites

```bash
alegra alias set unpaid "invoices list --status open --all"
alegra unpaid --client-id 12        # expande, luego anexa tus argumentos extra
alegra alias list
alegra alias remove unpaid
```

Los alias nunca tapan los comandos integrados, y los argumentos finales se anexan
a la expansión.

## Rangos de fecha naturales

Cada `list` acepta `--since`/`--until` con valores amigables:

```bash
alegra invoices list --since this-month
alegra invoices list --since last-month --until last-month
alegra payments  list --type in --since 7d            # últimos 7 días
alegra bills     list --since 2026-01-01 --until 2026-03-31
```

Aceptados: `YYYY-MM-DD`, `today`, `yesterday`, `tomorrow`, `this-month`,
`last-month`, `this-year`, `last-year`, `this-quarter`, y `Nd`/`Nw`/`Nm`/`Ny`.

## Conteos sin descargar

```bash
alegra invoices list --status open --count     # usa el total de la API
```

## Filtros arbitrarios (escape hatch)

Cualquier parámetro de consulta de Alegra que el CLI no exponga como flag:

```bash
alegra invoices list --param client_name="Acme" --param numberTemplate_fullNumber="FE-1"
```

## Validación previa

`create` revisa tu cuerpo contra las reglas del país antes de enviarlo (ver
[Facturación electrónica](electronic-invoicing.md)). Configura el país una vez:

```bash
alegra config set-country colombia
```

Omite una validación con `--no-validate`; crea borradores internos con `--draft`.

## Errores más amigables

Las fallas se explican solas y sugieren una solución — `402` se vuelve "tu plan no
incluye esto", `429` se vuelve una pista de límite de velocidad, y los códigos de
timbrado se mapean a soluciones. Ver la [Referencia de errores](../reference/errors.md).

## Manéjalo desde un agente de IA

El CLI completo es un servidor MCP: `claude mcp add alegra -- alegra mcp start` (o
`alegra mcp claude enable`). Ver [Servidor MCP](../user-guide/mcp.md).
