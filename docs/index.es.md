---
title: Inicio
---

# Alegra CLI

Una interfaz de línea de comandos rápida y scriptable para la API de contabilidad
de [Alegra](https://www.alegra.com/).

Gestiona contactos, facturas, ítems, pagos, impuestos y toda la superficie de
recursos de Alegra desde tu terminal — con salida `table`/`json`/`yaml`/`csv`,
perfiles con nombre, un modo dry-run y un [servidor MCP](user-guide/mcp.md)
integrado.

!!! note "No oficial"
    Esta es una herramienta de la comunidad, sin afiliación con Alegra. Usa la API
    pública en `https://api.alegra.com/api/v1`.

## Enlaces rápidos

- [Instalación](getting-started/installation.md) · [Autenticación](getting-started/authentication.md) · [Inicio rápido](getting-started/quickstart.md)
- **[Cookbook](cookbook.md)** — recetas listas para copiar y pegar en tareas del día a día
- Guías: [De factura a cobro](guides/invoice-to-cash.md) · [Gastos y compras](guides/expenses-and-purchases.md) · [Facturación electrónica](guides/electronic-invoicing.md) · [Reportes y cierre de mes](guides/reporting-and-month-end.md) · [Automatización](guides/automation.md)
- Referencia: [Comandos](commands/index.md) · [Errores y límites de tasa](reference/errors.md) · [Preguntas frecuentes](reference/faq.md)

## De un vistazo

```bash
alegra auth login
alegra contacts list --type client --all
alegra invoices get 12 -o json
alegra invoices create -f new-invoice.json
alegra invoices void 12
```
