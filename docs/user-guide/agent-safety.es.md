---
title: Seguridad para agentes
description: Evita que un agente de IA ejecute operaciones contables destructivas, con hooks y configuración por host para Claude Code, Codex y OpenCode.
---

# Seguridad para agentes: controlar operaciones destructivas

Dejar que un agente de IA toque tu contabilidad es útil pero riesgoso: puede
crear facturas, registrar pagos, eliminar registros y emitir facturas
electrónicas ante la DIAN/SAT. Esta página muestra cómo asegurarte de que un
agente **no pueda** ejecutar las operaciones que no quieres — ya sea que maneje
`alegra` por un shell o por el servidor [`alegra mcp`](mcp.md).

Ten claro qué hace cada capa:

| Capa | Mecanismo | Qué hace | ¿Hace cumplir? |
| --- | --- | --- | --- |
| 1. Anotaciones de tools | `alegra mcp` marca cada tool como solo lectura o destructiva | Permite que un host que respete anotaciones controle las escrituras por su cuenta | **Advisoria** — depende del host |
| 2. Config del host | reglas deny / ask **y un hook PreToolUse** | El hook inspecciona el comando real y no se esquiva con trucos del shell — el bloqueo definitivo | **Sí** — esta es la barrera real |
| 3. Built-ins del CLI | `--dry-run`, confirmación de `delete` | Previsualización y confirmación manual en una terminal | Solo shell (no la superficie MCP) |

El resumen honesto: **la capa 1 hace que los buenos hosts se porten bien por
defecto, pero no es un límite de seguridad. La capa 2 es la que de verdad
bloquea.** Usa ambas.

!!! tip "El hook es la capa que no se puede esquivar"
    Los globs de `deny`/`ask` son fáciles de leer, pero un comando de shell bien
    armado puede esquivarlos (comillas, subshells, `&&`). Un **hook PreToolUse**
    corre tu código contra el comando o tool *real* antes de ejecutarlo, así que
    es la capa que **bloquea de forma definitiva** — en las superficies shell y
    MCP. `alegra agent guard --host claude-code` te lo genera. (Los hooks son una
    función de Claude Code; Codex bloquea con un sandbox de solo lectura,
    OpenCode con reglas `deny`.)

## La vía rápida: `alegra agent guard`

`alegra` te genera la config de la capa 2, con las operaciones destructivas
derivadas del árbol de comandos en vivo (así la lista siempre está completa):

```bash
alegra agent guard --host claude-code   # imprime settings.json + el hook PreToolUse
alegra agent guard --host codex         # imprime config.toml
alegra agent guard --host opencode      # imprime opencode.json
```

Por defecto **bloquea de forma definitiva las acciones irreversibles** (`delete`,
`void`, `emit`, `stamp`, `close`, y las acciones `*-delete`) y hace que las
escrituras normales (`create`, `update`, `import`) **requieran aprobación**; las
lecturas quedan permitidas. Flags:

- `--all-writes` — bloquea toda escritura, no solo las irreversibles.
- `--write` — instala los archivos en vez de imprimirlos (nunca sobreescribe una
  config existente; esa la imprime para que la mergees).

Revisa la salida, pégala en tu host y listo. El resto de la página explica qué
genera y cómo ajustarla a mano.

## Qué cuenta como destructivo

Las operaciones de lectura de `alegra` son `list`, `get` y `export` (más
`catalog`, `reports`, `doctor`, `version`). Todo lo demás escribe: `create`,
`update`, `delete`, `import`, y las acciones de recurso `void`, `emit`, `stamp`,
`email`, `open`, `close`, `transfer`, …

## Capa 1 — anotaciones MCP integradas

`alegra mcp` anuncia la naturaleza de cada tool con las anotaciones estándar de
MCP: las tools de lectura llevan `readOnlyHint: true` y las de
escritura/acciones llevan `destructiveHint: true`. Lo obtienes gratis, sin
configurar nada.

Lo que eso te da depende por completo del host:

- Un host que respeta anotaciones (p. ej. **Codex**) pedirá aprobación para las
  tools destructivas y dejará correr las de solo lectura, automáticamente.
- Un host que controla por nombre de tool (p. ej. **Claude Code**) ignora la
  anotación para las decisiones de permiso — ahí configuras reglas (capa 2).

Las anotaciones son advisorias según la spec de MCP: un host que las ignore
ejecuta todo. Nunca reemplazan la capa 2.

!!! note "Nombres de las tools"
    Las tools MCP se llaman `alegra_<recurso>_<subcomando>` — p. ej.
    `alegra_invoices_void`, `alegra_contacts_delete`, `alegra_invoices_list`.
    Los hosts les agregan namespace; en Claude Code una tool es
    `mcp__<server>__alegra_<recurso>_<subcomando>`.

## Capa 2 — enforcement por host

### Claude Code

Claude Code controla por **nombre de tool/comando**, y **`deny` siempre le gana a
`allow`**. Pon esto en el `.claude/settings.json` del proyecto (compartido con
el equipo) o en `~/.claude/settings.json` (todos los proyectos).

**Superficie shell** — bloquea los comandos `alegra` destructivos:

```json
{
  "permissions": {
    "deny": [
      "Bash(alegra * delete:*)",
      "Bash(alegra * create:*)",
      "Bash(alegra * update:*)",
      "Bash(alegra invoices emit:*)",
      "Bash(alegra invoices void:*)",
      "Bash(alegra invoices stamp:*)",
      "Bash(alegra invoices email:*)"
    ]
  }
}
```

Usa `"ask"` en vez de `"deny"` para pedir aprobación en lugar de bloquear de
plano. Un comando denegado nunca corre; en una sesión headless/CI `ask` también
falla cerrado (sin humano que apruebe → bloqueado).

**El bloqueo definitivo es un hook PreToolUse.** Los globs de `deny` de arriba
ayudan, pero un comando de shell bien armado puede esquivarlos (comillas,
subshells, `&&`); un hook corre tu código contra el comando real antes de
ejecutarlo, así que no se puede esquivar — y el mismo hook cubre también la
superficie MCP. Esta es la pieza más importante en Claude Code. Crea
`.claude/hooks/block-alegra-writes.sh`:

```bash
#!/usr/bin/env bash
# Lee el payload del hook desde stdin y deniega comandos alegra destructivos.
input=$(cat)
cmd=$(printf '%s' "$input" | jq -r '.tool_input.command // empty')

if printf '%s' "$cmd" | grep -qE 'alegra\b.*\b(delete|create|update|import|emit|void|stamp|email|open|close|transfer)\b'; then
  jq -n '{
    hookSpecificOutput: {
      hookEventName: "PreToolUse",
      permissionDecision: "deny",
      permissionDecisionReason: "Operación alegra destructiva bloqueada por política."
    }
  }'
fi
exit 0
```

Regístralo en `.claude/settings.json`:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          { "type": "command", "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/block-alegra-writes.sh" }
        ]
      }
    ]
  }
}
```

Un hook PreToolUse que imprime `permissionDecision: "deny"` (con exit code 0)
bloquea la llamada; el exit code 2 también bloquea y muestra stderr al modelo.
Es la opción más fuerte — corre tu lógica, no un glob.

**Superficie MCP** — cuando alegra corre como servidor MCP local
(`claude mcp add alegra -- alegra mcp start`), deniega las tools de escritura
por nombre. Como `deny` le gana a `allow`, deniega los verbos de escritura en
vez de intentar hacer allow-list de las lecturas:

```json
{
  "permissions": {
    "deny": [
      "mcp__alegra__alegra_.*_(create|update|delete|import|void|emit|stamp|email|open|close|transfer)"
    ]
  }
}
```

Los patrones de permiso MCP se evalúan como expresiones regulares, así que el
`.*` y el grupo `(...)` funcionan tal cual; las tools de lectura (`..._list`,
`..._get`, `..._export`) quedan intactas. Un hook PreToolUse también puede
controlar tools MCP — pon `"matcher": "mcp__alegra__.*"` e inspecciona
`.tool_name` en el script para semántica de allow-list (deniega cualquier cosa
que no sea un verbo de lectura).

### Codex

Codex controla con dos ajustes en `config.toml` (`~/.codex/config.toml`):
`sandbox_mode` y `approval_policy`.

**Superficie shell** — la barrera más fuerte y simple es quitar el acceso de
escritura:

```toml
sandbox_mode    = "read-only"   # Codex puede leer; editar/comandos/red requieren aprobación
approval_policy = "untrusted"   # solo las lecturas seguras corren solas; lo demás pregunta
```

Con `sandbox_mode = "read-only"`, Codex no puede correr un comando `alegra` que
cambie estado sin una aprobación explícita — sin lista de comandos que
mantener. Si quieres que las escrituras sean posibles pero siempre revisadas,
deja `approval_policy` en un modo que pregunte para que cada escritura se pause
ante un humano.

**Superficie MCP** — aquí es donde la capa 1 rinde. La doc de Codex dice que las
llamadas a tools MCP con anotación destructive **siempre requieren aprobación**.
Como `alegra mcp` marca `void`/`emit`/`delete`/`update`/… como destructivas,
Codex pregunta antes de ejecutarlas automáticamente, mientras que las de solo
lectura corren sin fricción. No necesitas config por tool; deja `approval_policy`
en un modo que pregunte para que la aprobación se respete.

Nota honesta: en Codex la barrera es **aprobación** (un humano dice que sí), no
un bloqueo duro — está bien para uso interactivo, pero implica que una corrida
desatendida de Codex debería usar `sandbox_mode = "read-only"` para que no haya
nada que aprobar.

### OpenCode

OpenCode controla con un bloque `permission` en `opencode.json`. Cada regla es
`"allow"`, `"ask"` o `"deny"`, y **gana el último patrón que coincide**, así que
pon el comodín primero y las denegaciones específicas después.

**Superficie shell** — `bash` acepta globs por patrón:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "permission": {
    "bash": {
      "*": "allow",
      "alegra * delete*": "deny",
      "alegra * create*": "ask",
      "alegra * update*": "ask",
      "alegra invoices emit*": "deny",
      "alegra invoices void*": "deny",
      "alegra invoices stamp*": "deny"
    }
  }
}
```

`deny` bloquea el comando; `ask` pide aprobación antes de ejecutarlo.

**Superficie MCP** — las mismas claves de permiso coinciden con los nombres de
las tools MCP como patrones comodín, así que puedes denegar las tools de
escritura de un servidor:

```json
{
  "permission": {
    "alegra_*_delete": "deny",
    "alegra_*_create": "ask",
    "alegra_*_update": "ask",
    "alegra_invoices_void": "deny",
    "alegra_invoices_emit": "deny",
    "alegra_invoices_stamp": "deny"
  }
}
```

Confirma el prefijo exacto que tu OpenCode usa para el servidor y haz coincidir
el patrón (el nombre de la tool es `alegra_<recurso>_<subcomando>`).

## Capa 3 — built-ins del CLI

Cuando el agente usa un **shell**, dos built-ins ayudan incluso sin config del
host:

- `--dry-run` en cualquier comando imprime la petición exacta y no envía nada.
- `delete` pide confirmación salvo que pases `-y`.

Estos no cubren la **superficie MCP**: una llamada a tool MCP no tiene terminal,
así que la confirmación interactiva de `delete` no aplica. En la superficie MCP,
apóyate en las capas 1 y 2.

## Configuración recomendada

Defensa en profundidad:

1. Deja que la capa 1 haga su trabajo — viene activa por defecto.
2. Agrega la capa 2 de tu host: una regla `deny` (o `ask`) **y** un hook en
   Claude Code; `sandbox_mode = "read-only"` en Codex; denegaciones de
   `permission` en OpenCode.
3. Mantén `--dry-run` en tus hábitos y scripts.

Pruébalo: pídele al agente que ejecute `alegra invoices void 1` (o que llame la
tool `alegra_invoices_void`) y confirma que queda bloqueado o pausado para
aprobación antes de que algo llegue a la API.
