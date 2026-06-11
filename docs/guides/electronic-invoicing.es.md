---
title: Facturación electrónica
---

# Facturación electrónica

Emitir un documento fiscal legalmente válido es la capacidad estrella de Alegra.
Esta guía explica el modelo y cómo manejarlo desde el CLI.

!!! abstract "Modelo mental"
    Una factura tiene un **ciclo de vida**: `draft → open → stamped`. *Timbrar*
    (timbrado/emisión) es el paso que envía el documento a la autoridad fiscal
    (DIAN en Colombia, SAT en México, SUNAT en Perú, Hacienda en Costa Rica…),
    que lo valida y devuelve un código fiscal único (CUFE / UUID / etc.).

## Las tres cosas que hacen "electrónica" a una factura

1. **Una resolución de numeración** — se pasa como `numberTemplate: { id }`. Esta
   lleva la autorización que te otorgó la autoridad fiscal (resolución DIAN, serie
   SAT/SUNAT). Encuentra las tuyas:
   ```bash
   alegra number-templates list -o json | jq '.[] | {id, name, documentType, prefix}'
   ```
2. **Un cliente que cumpla** — la `identification` del contacto debe coincidir con
   lo que exige el país (NIT para una empresa colombiana, RFC para México, RUC
   para Perú).
3. **El flag de timbrado** — `stamp: { generateStamp: true }`. Con él, Alegra
   emite el documento fiscal al crearlo. Sin él, obtienes un borrador interno.

## Emitir al crear

```bash
alegra invoices create -f invoice.json
```

`invoice.json` (Colombia):

```json
{
  "client": { "id": 12 },
  "numberTemplate": { "id": 7 },
  "date": "2026-06-06",
  "dueDate": "2026-06-21",
  "paymentForm": "CASH",
  "items": [
    { "id": 5, "price": 50000, "quantity": 1, "tax": [{ "id": 3 }] }
  ],
  "stamp": { "generateStamp": true }
}
```

Verifica el resultado y toma el código fiscal:

```bash
alegra invoices get <id> -o json | jq '{number, status, stamp}'
```

!!! tip "Previsualiza antes de enviar"
    Revisa siempre la petición primero:
    ```bash
    alegra invoices create -f invoice.json --dry-run
    ```

## Crear como borrador y emitir después (recomendado)

El flujo más seguro cuando una persona debe revisar antes: crear como **borrador**,
inspeccionarlo y luego emitir con `invoices emit`.

```bash
# 1. crea como borrador (elimina cualquier instrucción de timbrado)
alegra invoices create -f invoice.json --draft

# 2. revisa
alegra invoices get <id>

# 3. emite — envía a la autoridad fiscal para su timbrado
alegra invoices emit <id>
```

`invoices emit` está hecho para esto:

```bash
alegra invoices emit 1234 1235        # facturas específicas
alegra invoices emit --all            # todas las facturas en borrador
alegra invoices emit --all --dry-run  # previsualiza lo que se emitiría
```

**Divide automáticamente** en lotes de 10 (el tope por llamada de Alegra) y
mantiene una **protección de idempotencia** local: una factura que ya se emitió
desde esta máquina se omite a menos que pases `--force`. Como la emisión no es
idempotente del lado del servidor, esta es tu red de seguridad contra documentos
duplicados.

!!! tip "Validación previa"
    `invoices create` valida el cuerpo para tu país antes de enviarlo y rechaza
    errores evidentes (cliente/ítems faltantes, o `stamp` sin un
    `numberTemplate`). Configura tu país una sola vez para que conozca las reglas:
    ```bash
    alegra config set-country colombia
    ```
    Omite una validación con `--no-validate`.

## Otros documentos electrónicos

| Documento | Comando | Notas |
|---|---|---|
| Factura de venta | `alegra invoices …` | `stamp`, `void`, `open`, `email`, `stamp` (masivo), `preview` |
| Nota de crédito | `alegra credit-notes …` | electrónica en CO/MX; `email` |
| Nota de débito (proveedor) | `alegra debit-notes …` | |
| Nota de débito al cliente | `alegra income-debit-notes …` | |
| Factura global (MX) | `alegra global-invoices …` | consolida ventas del período sin un RFC individual |
| Recibo de pago (REP) | `alegra payments stamp <id>` | complemento de pago de Costa Rica / México |
| Documento soporte | `alegra number-templates …` (`documentType=supportDocument`) | documento soporte (CO) |

## Referencia rápida por país

Estos son los campos que exige la autoridad fiscal; Alegra los aplica como
`400` del lado del servidor, así que modélalos en tu payload desde el inicio.

=== "Colombia (DIAN)"
    - `identification`: `{ type: "NIT" | "CC" | "CE" | "PP", number }` (NIT para empresas; Alegra calcula el DV).
    - `numberTemplate.id`: una resolución autorizada por la DIAN.
    - `stamp.generateStamp: true` → emite y devuelve el **CUFE** (verificable en DIAN MUISCA).
    - Impuesto común: **IVA 19% = tax id `3`** (confírmalo con `alegra taxes list`).

=== "México (SAT / CFDI 4.0)"
    - **RFC** del receptor, **régimen fiscal** (el emisor y el receptor deben ser compatibles), **uso CFDI** y el **código postal del receptor debe coincidir** con su Constancia de Situación Fiscal.
    - **método de pago** PUE (un solo pago) vs PPD (diferido). PPD luego requiere un **complemento de pago (REP)**.
    - El timbrado devuelve el **UUID** (folio fiscal). Usa `global-invoices` para público en general.

=== "Perú (SUNAT)"
    - **RUC** del cliente; emisión a través del OSE/PSE de Alegra.
    - Las **boletas** se reportan a la SUNAT como un **resumen diario consolidado**, no de forma individual — no esperes un acuse instantáneo por cada boleta.

=== "Costa Rica (Hacienda 4.4)"
    - Comprobantes v4.4. Las ventas a crédito deben emitir después un **REP (Recibo Electrónico de Pago)** — `alegra payments stamp <id>`.

=== "Argentina / República Dominicana"
    - Argentina usa el **CAE** de AFIP; República Dominicana usa **NCF / e-CF**. La
      forma de los campos varía — confírmalos contra tu cuenta real antes de
      automatizar.

## La emisión no es idempotente del lado del servidor — deja que el CLI la proteja

La API de Alegra no tiene clave de idempotencia, así que reejecutar a ciegas un
timbrado en crudo **puede crear un documento duplicado** (y la autoridad fiscal
puede rechazar numeración duplicada). Por eso justamente deberías emitir a través
de `alegra invoices emit`: mantiene una **protección de idempotencia** local y
omite cualquier factura que ya se haya emitido desde esta máquina a menos que
pases `--force`.

- Prefiere `alegra invoices emit` antes que automatizar reintentos de `create … --set stamp=…` /
  `invoices stamp` en crudo — `emit` es la ruta protegida.
- Usa siempre `--dry-run` primero.
- Si una llamada parece haber fallado, **revisa el estado** (`alegra invoices get <id>`)
  antes de reintentar — puede que sí se haya emitido.
- Nunca automatices un bucle de reintentos a ciegas alrededor del timbrado en
  crudo; deja que `emit` maneje los reintentos (o pasa `--force` una sola vez, de
  forma deliberada).

## Las facturas son de solo anexar

Una factura timbrada no se "corrige". Para revertir o corregir una, emite una
**nota de crédito** (`alegra credit-notes create`). El CLI trata los documentos
fiscales como historia inmutable de forma intencional.

## Dos APIs de Alegra

Este CLI usa la **App API** de Alegra (`api.alegra.com/api/v1`). Si en cambio
necesitas firmar y transmitir XML **directamente** a la DIAN (eres un proveedor de
software que usa Alegra solo como backend de timbrado), esa es la **e-Provider API**
aparte (`api.alegra.com/e-provider/col/v1`, autenticación Bearer) — que no se cubre
aquí.
