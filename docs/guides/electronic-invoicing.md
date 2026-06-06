---
title: Electronic Invoicing
---

# Electronic Invoicing (facturación electrónica)

Emitting a legally valid fiscal document is Alegra's headline capability. This
guide explains the model and how to drive it from the CLI.

!!! abstract "Mental model"
    An invoice has a **lifecycle**: `draft → open → stamped`. *Stamping*
    (timbrado/emisión) is the step that sends the document to the tax authority
    (DIAN in Colombia, SAT in Mexico, SUNAT in Peru, Hacienda in Costa Rica…),
    which validates it and returns a unique fiscal code (CUFE / UUID / etc.).

## The three things that make an invoice "electronic"

1. **A numbering resolution** — passed as `numberTemplate: { id }`. This carries
   the authorization the tax authority granted you (DIAN resolution, SAT/SUNAT
   series). Find yours:
   ```bash
   alegra number-templates list -o json | jq '.[] | {id, name, documentType, prefix}'
   ```
2. **A compliant client** — the contact's `identification` must match what the
   country requires (NIT for a Colombian company, RFC for Mexico, RUC for Peru).
3. **The stamp flag** — `stamp: { generateStamp: true }`. With it, Alegra emits
   the fiscal document on creation. Without it, you get an internal draft.

## Emit on creation

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

Verify the result and grab the fiscal code:

```bash
alegra invoices get <id> -o json | jq '{number, status, stamp}'
```

!!! tip "Preview before you send"
    Always check the request first:
    ```bash
    alegra invoices create -f invoice.json --dry-run
    ```

## Create as a draft, then emit later (recommended)

The safest flow when a human should review first: create as a **draft**, inspect
it, then emit with `invoices emit`.

```bash
# 1. create as a draft (strips any stamp instruction)
alegra invoices create -f invoice.json --draft

# 2. review
alegra invoices get <id>

# 3. emit — sends to the tax authority for stamping
alegra invoices emit <id>
```

`invoices emit` is built for this:

```bash
alegra invoices emit 1234 1235        # specific invoices
alegra invoices emit --all            # every draft invoice
alegra invoices emit --all --dry-run  # preview what would emit
```

It **auto-chunks** into batches of 10 (Alegra's per-call cap) and keeps a local
**idempotency guard**: an invoice already emitted from this machine is skipped
unless you pass `--force`. Because emission isn't idempotent server-side, this is
your safety net against duplicate documents.

!!! tip "Pre-flight validation"
    `invoices create` validates the body for your country before sending and
    refuses obvious mistakes (missing client/items, or `stamp` without a
    `numberTemplate`). Set your country once so it knows the rules:
    ```bash
    alegra config set-country colombia
    ```
    Bypass a check with `--no-validate`.

## Other electronic documents

| Document | Command | Notes |
|---|---|---|
| Sales invoice | `alegra invoices …` | `stamp`, `void`, `open`, `email`, `stamp` (bulk), `preview` |
| Credit note | `alegra credit-notes …` | electronic in CO/MX; `email` |
| Debit note (vendor) | `alegra debit-notes …` | |
| Customer debit note | `alegra income-debit-notes …` | |
| Global invoice (MX) | `alegra global-invoices …` | consolidates period sales without an individual RFC |
| Payment receipt (REP) | `alegra payments stamp <id>` | Costa Rica / Mexico complemento de pago |
| Support document | `alegra number-templates …` (`documentType=supportDocument`) | documento soporte (CO) |

## Country quick-reference

These are the fields the tax authority requires; Alegra enforces them as
server-side `400`s, so model them in your payload up front.

=== "Colombia (DIAN)"
    - `identification`: `{ type: "NIT" | "CC" | "CE" | "PP", number }` (NIT for companies; DV computed by Alegra).
    - `numberTemplate.id`: a DIAN-authorized resolution.
    - `stamp.generateStamp: true` → emits, returns **CUFE** (verifiable in DIAN MUISCA).
    - Common tax: **IVA 19% = tax id `3`** (confirm with `alegra taxes list`).

=== "Mexico (SAT / CFDI 4.0)"
    - Receptor **RFC**, **régimen fiscal** (issuer + receptor must be compatible), **uso CFDI**, and the receptor **código postal must match** their Constancia de Situación Fiscal.
    - **método de pago** PUE (one payment) vs PPD (deferred). PPD later needs a **complemento de pago (REP)**.
    - Stamping returns the **UUID** (folio fiscal). Use `global-invoices` for público-en-general.

=== "Peru (SUNAT)"
    - Client **RUC**; emission via Alegra's OSE/PSE.
    - **Boletas** are reported to SUNAT as a **consolidated daily summary**, not individually — don't expect an instant per-boleta acknowledgement.

=== "Costa Rica (Hacienda 4.4)"
    - Comprobantes v4.4. Credit sales must later issue a **REP (Recibo Electrónico de Pago)** — `alegra payments stamp <id>`.

=== "Argentina / Dominican Republic"
    - Argentina uses AFIP **CAE**; Dominican Republic uses **NCF / e-CF**. Field
      shapes vary — confirm against your live account before scripting.

## Important: emission is not idempotent

There is no idempotency key on the Alegra API. **Re-running an emit can create a
duplicate document** (and the tax authority may reject duplicate numbering). Until
the CLI ships an emission guard:

- Always `--dry-run` first.
- After a failed-looking call, **check status** (`alegra invoices get <id>`)
  before retrying — it may have emitted.
- Never script a blind retry loop around stamping.

## Invoices are append-only

You don't "fix" a stamped invoice. To reverse or correct one, issue a **credit
note** (`alegra credit-notes create`). The CLI intentionally treats fiscal
documents as immutable history.

## Two Alegra APIs

This CLI uses Alegra's **App API** (`api.alegra.com/api/v1`). If you instead need
to sign and transmit XML **directly** to DIAN (you're a software vendor using
Alegra only as the stamping backend), that's the separate **e-Provider API**
(`api.alegra.com/e-provider/col/v1`, Bearer auth) — not covered here.
