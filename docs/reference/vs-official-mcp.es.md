# alegra-cli frente al servidor MCP oficial de Alegra

Alegra publica un **servidor MCP** alojado (`mcp.alegra.com`) para que los agentes de IA
puedan llamar a la API de contabilidad sobre el Model Context Protocol. `alegra-cli` sigue
otro camino: está construido **pensando primero en el agente**, como una herramienta de
línea de comandos que los agentes de IA manejan directamente — y *además* trae su propio
servidor MCP (`alegra mcp`) y funciona para personas en una terminal. Esta página compara
ambos para que elijas el que mejor te encaje.

!!! note "Foto del momento"
    Refleja el **conjunto de herramientas desplegado en vivo del servidor MCP oficial (49
    herramientas de solo lectura)** y `alegra-cli` **v0.8.1**, verificado el **2026-06-11**
    conectándose a `mcp.alegra.com`. Ambos envuelven la misma API v1 de Alegra y ambos
    requieren acceso de red a ella. Alegra puede cambiar su MCP con el tiempo — trátalo como
    una foto de un momento puntual.

## De un vistazo

| | MCP oficial de Alegra | alegra-cli |
| --- | --- | --- |
| Interfaz | Servidor MCP alojado (`mcp.alegra.com`, streamable HTTP) | Terminal + **skill** de agente + servidor `alegra mcp` |
| Autenticación | OAuth 2.0 (login en el navegador, authorization-code + PKCE; scope `owner`) | Token de API de Alegra (Basic), guardado en el **keyring del sistema** |
| Cómo se conecta un agente | Un host MCP interactivo cuyo callback Alegra puso en su lista de permitidos — **claude.ai** web o **ChatGPT** (ejecuta el flujo OAuth) | Un comando en **Claude Code**, Cursor o Codex vía una [skill](../user-guide/agent-skill.md); **o** [`alegra mcp`](../user-guide/mcp.md) para hosts MCP |
| Headless / CI / cron | No — OAuth necesita un login interactivo en el navegador (sin `client_credentials` ni device grant) | Sí — el token de API funciona desatendido |
| Instalación | Ninguna — alojado por Alegra | Instalas un binario (Homebrew, Scoop, Docker, …) |
| Salida | Resultados de herramienta en JSON estructurado | tabla / JSON / YAML / CSV — combinable con `jq` y pipes (menos tokens para un agente) |
| Controles de seguridad para agentes | Depende del host | `--dry-run` en cualquier comando, confirmación interactiva en `delete`, restringible con hooks del shell |
| Uso humano (terminal) | — | De primera clase |
| Cobertura de recursos | 17 áreas (49 herramientas) | 45+ recursos |
| Escrituras (create / update / delete) | No expuestas — solo lectura | Sí |
| Emisión de facturas (stamp / void / email) | No expuesta | Sí |

## Construido pensando primero en el agente

`alegra-cli` se diseñó asumiendo que quien llama es, ante todo, un agente de IA y no una
persona escribiendo comandos. El razonamiento está expuesto en el post
[CLIs over MCPs](https://jjuanrivvera.com/es/blog/clis-sobre-mcps/); en resumen:

- **Se conecta a agentes de shell en un solo paso.** Un agente de código (Claude Code,
  Cursor, Codex) instala el binario, configura un token de API y ejecuta comandos — sin
  flujo OAuth, sin navegador, sin configurar un host. El MCP oficial usa OAuth 2.0 con una
  **lista de permitidos de redirect-URI cerrada**: al sondear su endpoint de registro,
  solo se aceptan los callbacks de los partners de integración de Alegra (`claude.ai` y
  ChatGPT) — los de loopback (`localhost`/`127.0.0.1`), los schemes de URI custom
  (`vscode://`, `cursor://`) y los HTTPS arbitrarios se rechazan todos. Por eso ningún
  cliente MCP de terminal, CLI o self-hosted (Claude Code, Cursor, OpenCode, OpenClaw,
  Hermes, …) puede completar el flujo. Con `alegra-cli`, el mismo token de API también
  funciona en CI, cron y scripts, donde un login OAuth interactivo no es posible.
- **Salida combinable.** Un agente puede elegir `--columns`, redirigir a `jq` con un pipe,
  o usar `grep` — enviando muchos menos tokens al modelo que un resultado de herramienta en
  JSON crudo.
- **La seguridad es aplicable.** `--dry-run` previsualiza la petición exacta en cualquier
  comando, `delete` pide confirmación, y un shell puede bloquear comandos destructivos con
  hooks — de modo que un agente autónomo *no puede* borrar aunque lo intente. Mira
  [Seguridad para agentes](../user-guide/agent-safety.md) para el control por host (Claude
  Code, Codex, OpenCode).
- **Dos vías de entrada para agentes.** Un agente de código con acceso al shell usa la
  **skill** (aprende las reglas de oro y cuándo llamar a cada comando); un host del
  protocolo MCP usa **`alegra mcp`**, que expone el mismo árbol de comandos como
  herramientas MCP.

!!! note "Sobre las credenciales"
    Ambas herramientas manejan las credenciales de forma responsable — no son un punto de
    diferencia. `alegra-cli` guarda un token de API de larga vida en el keyring del sistema;
    el MCP oficial usa tokens OAuth de corta vida que el modelo nunca ve. La verdadera
    diferencia es el alcance: un token de API corre desatendido (CI, cron, agentes de
    shell), mientras que OAuth exige un login interactivo de una sola vez en el navegador.

## Cobertura de recursos

Ambos envuelven la misma API v1 de Alegra pero exponen porciones distintas de ella. La
diferencia clave: las 49 herramientas del MCP oficial son todas de **solo lectura**
(`get…`/`list…`), mientras que `alegra-cli` lee *y* escribe.

### Leídos por ambos, escritos solo por alegra-cli
Contactos; ítems (con stock por bodega vía `items stock`) y la familia de inventario
(categorías de ítems, atributos de variantes, bodegas, transferencias, listas de precios,
ajustes y numeraciones, campos personalizados); cuentas bancarias y conciliaciones;
facturas; la familia de gastos (facturas de compra, notas débito de proveedor, órdenes de
compra, pagos salientes); reportes de ventas; y los catálogos de referencia de unidades /
claves de producto del SAT. El MCP oficial los expone en **solo lectura**; `alegra-cli`
además los crea, actualiza y elimina.

### Solo en alegra-cli
Recursos que el MCP oficial no expone en absoluto: libros diarios, centros de costo,
impuestos, retenciones, monedas, vendedores, cotizaciones, notas crédito, notas débito de
cliente (ingresos), remisiones, recibos de transporte, facturas globales (CFDI), facturas y
pagos recurrentes, numeraciones de documentos, términos de pago, cargos adicionales, pagos
entrantes y suscripciones a webhooks. Además de toda operación de escritura y la emisión de
facturas electrónicas (ver más abajo).

**Claves de producto del SAT** (`claveProdServ`, ~52k entradas específicas de México): el
MCP oficial expone una herramienta de solo lectura `config_getProductKeys`. `alegra-cli`
sincroniza el catálogo publicado por el SAT bajo demanda (`alegra catalog sync-sat`,
ofrecido por `alegra init` en cuentas mexicanas) y lo busca sin conexión con
`alegra catalog product-keys <query>`, sin necesidad de conexión.

### Solo en el MCP oficial
- Los ayudantes de mesa de ayuda del **Support Center** y el **Task Manager** (tareas) —
  fuera del alcance de una CLI de contabilidad.

## Diferencias de capacidad que importan

**Facturación electrónica.** Las herramientas de facturas del MCP oficial son de solo
lectura (`getInvoices`, `getInvoiceById`/`ByNumber`, más un generador de documentos
exportables). No puede crear, actualizar, `stamp`, `void` ni `email` una factura. Emitir una
factura electrónica a la DIAN/SAT, anularla o enviarla por correo se hace con `alegra-cli`
(`alegra invoices stamp|void|email …`) o con la API REST directamente. Si la emisión
electrónica es parte de tu flujo de trabajo, la CLI cubre el ciclo completo.

**Pagos.** El MCP oficial expone únicamente pagos **salientes**, en solo lectura, dentro de
sus herramientas de gastos. `alegra-cli` tiene un único recurso `payments` que cubre tanto
entrantes como salientes, con lectura, escritura y `stamp`/`void`/`open`.

**Reportes.** Los dos conjuntos de reportes son complementarios:

| Reporte | MCP oficial | alegra-cli |
| --- | --- | --- |
| Ventas por cliente / totales | ✓ | ✓ |
| Ventas por vendedor | ✓ | ✓ |
| Ventas por vendedor totales · ventas generales | ✓ | — |
| Estado de cuenta · estado de resultados | — | ✓ |

## ¿Cuál debería usar?

Sirven a nichos distintos; ambos son válidos.

| Tu situación | Mejor opción |
| --- | --- |
| Consulta de solo lectura de tus datos desde claude.ai o ChatGPT | **MCP oficial** (hecho a propósito para eso) |
| Agente de código con acceso al shell (Claude Code, Cursor, Codex) | **alegra-cli + skill** |
| Un host MCP que el servidor oficial no soporta, o necesitas herramientas de escritura | El `alegra mcp` de **alegra-cli** |
| Scripts, CI/CD, jobs de cron | **alegra-cli** (pipes, códigos de salida, variables de entorno) |
| Agente autónomo tocando datos de producción | **alegra-cli** (dry-run, confirmaciones, hooks) |
| Emisión de facturas electrónicas | **alegra-cli** |
| Trabajar de forma interactiva en una terminal | **alegra-cli** |

Como el servidor MCP de `alegra-cli` se genera a partir de su árbol de comandos, cada
operación contable que puedas ejecutar en la terminal está disponible para un agente (los
comandos locales de setup y credenciales quedan fuera de la superficie de tools). Mira
[Skill para agentes](../user-guide/agent-skill.md), [Servidor MCP](../user-guide/mcp.md) y
[Seguridad para agentes](../user-guide/agent-safety.md).
