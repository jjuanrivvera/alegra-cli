---
title: Preguntas frecuentes y resolución de problemas
---

# Preguntas frecuentes y resolución de problemas

### ¿Cómo me autentico en un script / CI?
Configura variables de entorno — sin necesidad del keyring:
```bash
export ALEGRA_EMAIL="you@biz.com" ALEGRA_TOKEN="…"
alegra auth status
```
Las variables de entorno tienen prioridad sobre el archivo de configuración y el keyring. Mira [Configuración](../user-guide/configuration.md).

### ¿De dónde saco mi token de API?
Alegra → **Configuración → API / Integraciones con otros sistemas**. La página
muestra el email a usar y te deja generar un token. La autenticación es HTTP Basic
(`base64(email:token)`) — la CLI se encarga de la codificación.

### `alegra invoices create` falla con HTTP 400. ¿Por qué?
Casi siempre es un campo requerido por el país que falta. Colombia necesita un
`numberTemplate.id` (resolución) y un cliente con el tipo de `identification`
correcto (NIT para empresas); México necesita régimen/uso CFDI y un CP que coincida. Ejecuta
con `--dry-run` para inspeccionar el cuerpo y mira [Facturación electrónica](../guides/electronic-invoicing.md).

### ¿Cómo paso un filtro para el que la CLI no tiene flag?
Usa la salida de escape:
```bash
alegra invoices list --param client_name="Acme" --param numberTemplate_fullNumber="FE-123"
```
Cualquier parámetro de consulta de Alegra funciona con `--param key=value` (repetible).

### ¿Por qué `list` solo devuelve 30 filas?
Alegra limita `limit` a **30** por página. Usa `--all` para traer todas las páginas
automáticamente, o `--start`/`--limit` para paginar manualmente. Para un total sin
descargar, usa `--count`.

### ¿Cómo obtengo solo un conteo?
```bash
alegra invoices list --status open --count
```
Usa el total de metadatos de la API en lugar de descargar las filas.

### ¿Puedo editar, anular o eliminar una factura?
Puedes hacer `update` a un borrador, pero **los documentos fiscales timbrados son
inmutables** por diseño. Para revertir uno:

- **Anular / reabrir**: `alegra invoices void <id>` (anular) o `alegra invoices open <id>`.
  Los pagos tienen lo mismo: `alegra payments void <id>` / `alegra payments open <id>`.
- **Corregir el monto**: emite una nota crédito (`alegra credit-notes create`).

La CLI sigue la disciplina append-only de la contabilidad.

### Un comando devuelve HTTP 402 — ¿está roto?
No — `402` significa que tu **plan de Alegra no incluye** esa función (algunos reportes
y funciones de complementos/marketplace). Es un límite de cuenta/plan, no un bug de la CLI.
Ejecuta `alegra doctor` para ver qué endpoints puede alcanzar tu plan. Mira la
[Referencia de errores](errors.md).

### `currencies get USD` — ¿por qué un código y no un número?
Las monedas se identifican por su **código** ISO, no por un id numérico. `get`/`update`/
`delete` toman el código: `alegra currencies get USD`.

### ¿Cómo cambio entre cuentas de clientes?
Usa [perfiles](../user-guide/configuration.md):
```bash
alegra config use clienteA
alegra --profile clienteB invoices list
```

### La salida es difícil de leer / necesito una hoja de cálculo.
Elige un formato: `-o table` (por defecto), `-o json`, `-o yaml`, `-o csv`, y acota las
columnas con `--columns`:
```bash
alegra contacts list -o csv --columns id,name,email > contacts.csv
```

### ¿Hay un sandbox?
La App API corre contra tu cuenta real. Usa una **cuenta de Alegra de prueba/demo**
dedicada (un perfil aparte) para experimentar, y `--dry-run` para previsualizar
escrituras sin enviarlas.

### ¿Cómo reporto un bug o pido una función?
Abre un issue: <https://github.com/jjuanrivvera/alegra-cli/issues>.
