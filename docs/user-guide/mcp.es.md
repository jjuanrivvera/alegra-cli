---
title: Servidor MCP de Alegra (alegra mcp)
description: Usa alegra-cli como servidor MCP de Alegra. alegra mcp convierte cada comando en una herramienta MCP para Claude, Cursor y otros agentes de IA.
---

# Servidor MCP

alegra-cli puede exponer todo su árbol de comandos como un servidor
[Model Context Protocol](https://modelcontextprotocol.io/), para que los agentes
de IA (Claude, etc.) puedan manejar Alegra a través de herramientas bien
descritas — una por comando.

```bash
alegra mcp --help
```

Cada comando de la CLI se convierte en una herramienta MCP llamada
`alegra_<resource>_<action>` (por ejemplo, `alegra_invoices_list`,
`alegra_contacts_create`) cuyo esquema de entrada se deriva de las flags y los
argumentos de ese comando.

Las credenciales se resuelven exactamente igual que en la CLI (perfil →
keyring del sistema/env), y las flags sensibles `--show-token` y `--profile` se
excluyen de los esquemas de herramientas expuestos.

!!! warning
    El servidor MCP puede crear, modificar y eliminar registros contables reales.
    Ejecútalo contra un perfil de sandbox mientras desarrollas, y revisa las
    acciones que tu agente tiene permitido realizar.

## Conéctalo a tu editor automáticamente

En lugar de editar los archivos de configuración a mano, deja que la CLI los
escriba por ti. Hay subcomandos para los hosts MCP más comunes:

```bash
alegra mcp claude enable     # agrega el servidor a la config de Claude Desktop
alegra mcp cursor enable     # Cursor
alegra mcp vscode enable     # VS Code

alegra mcp claude list       # muestra los servidores configurados
alegra mcp claude disable    # vuelve a quitarlo
```

Cada `enable` escribe la config MCP del host para lanzar `alegra mcp start` (el
servidor stdio) con tu configuración actual.

## Transportes

| Comando | Transporte | Uso |
| --- | --- | --- |
| `alegra mcp start` | stdio | El predeterminado; lo que lanza la integración con el editor. También: `claude mcp add alegra -- alegra mcp start`. |
| `alegra mcp stream` | HTTP | Servidor de larga duración para agentes remotos/compartidos. Flags: `--host`, `--port` (por defecto 8080), `--log-level`. |
| `alegra mcp tools` | — | Imprime el esquema completo de herramientas como JSON (una herramienta por comando) — útil para inspeccionar o alimentar otro sistema. |

```bash
# Ejecuta el servidor HTTP en un puerto personalizado
alegra mcp stream --port 9000 --log-level debug

# Inspecciona las herramientas expuestas
alegra mcp tools | jq '.[].name'
```

Luego pídele al agente que, por ejemplo, "liste las facturas abiertas de este mes
en Alegra".

## Mantenlo seguro

Cada tool expuesta está anotada como solo lectura o destructiva, así que un host
que respete las anotaciones MCP controla las escrituras por su cuenta. Para
realmente bloquear las tools destructivas (`void`, `emit`, `delete`, …) del lado
del agente, mira [Seguridad para agentes](agent-safety.md).
