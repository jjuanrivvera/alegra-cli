---
title: Skill de agente
---

# Úsala desde tu agente de IA (skill)

alegra-cli incluye una **skill de agente** — un `SKILL.md` que le enseña a los
agentes de IA de programación (Claude Code, Cursor, Codex, Gemini CLI, Windsurf,
GitHub Copilot y más) cómo y cuándo manejar la CLI `alegra`. Una vez instalada,
puedes simplemente pedirle a tu agente cosas como *"crea una factura en borrador
para el cliente 12 y envíasela por correo"* y ya conoce los comandos, las flags y
las reglas de seguridad.

## Instala la skill

La forma multi-agente (recomendada) — detecta cada agente que tengas y se
mantiene al día:

```bash
npx skills add jjuanrivvera/alegra-cli
```

Integrada en la CLI (sin necesidad de Node) — escribe la skill empaquetada
directamente:

```bash
alegra skills install                 # proyecto actual (./.claude/skills)
alegra skills install --global        # a nivel de usuario (~/.claude/skills)
alegra skills install --agent cursor --global
alegra skills path                    # muestra dónde se instalaría
alegra skills print                   # imprime el SKILL.md
```

Plugin nativo de Claude Code:

```text
/plugin marketplace add jjuanrivvera/alegra-cli
/plugin install alegra-cli@alegra
```

## Requisitos previos

La skill **envuelve** el binario `alegra` — no lo empaqueta. Así que:

1. Instala la CLI: `brew install jjuanrivvera/alegra-cli/alegra-cli`
   (o `go install github.com/jjuanrivvera/alegra-cli/cmd/alegra@latest`).
2. Autentícate: `alegra auth login` (o define `ALEGRA_EMAIL`/`ALEGRA_TOKEN`).
3. Verifica: `alegra doctor`.

## Skill vs MCP

- **Skill** (esta página) — instrucciones en markdown para agentes con acceso a
  shell. El agente ejecuta los comandos `alegra …` directamente. La menor
  fricción; ideal para uso por CLI.
- **[Servidor MCP](mcp.md)** — `alegra mcp` expone cada comando como una
  herramienta estructurada, para clientes que prefieren esquemas de herramientas
  (Claude Desktop, etc.).

Usa la que tu agente admita; pueden coexistir.

## Qué aprende el agente

La skill enseña el flujo de autenticar → descubrir → actuar → verificar, las
reglas de oro (previsualiza las escrituras con `--dry-run`, parsea con `-o json`,
las facturas son de solo agregado, prefiere `--count`), el mapa de
recursos/acciones, cómo escribir cuerpos anidados, la emisión de facturación
electrónica, la importación/exportación masiva por CSV y el manejo de errores.
Junto a ella se incluye un cheatsheet condensado en
`references/alegra-commands.md`.
