---
name: alegra-cli
description: Manage the Alegra accounting API (https://alegra.com) from the terminal with the `alegra` CLI — contacts/clients, items, sales invoices, payments, bills, credit notes, taxes, reports, and DIAN/SAT/SUNAT/Hacienda electronic invoicing. Use this whenever the user wants to create or list invoices, manage clients or products, record or reconcile payments, run sales reports, bulk-import accounting data, or emit electronic (fiscal) invoices. Alegra is one API whose behavior is localized per country (Colombia, Mexico, Peru, Costa Rica, Spain, Panama, Argentina, Chile, the Dominican Republic and more); detect the connected country/version with `alegra company get` before any country-specific write.
version: 0.4.0
homepage: https://github.com/jjuanrivvera/alegra-cli
license: MIT
allowed-tools: Bash(alegra:*)
metadata: {"openclaw":{"category":"finance","emoji":"🧾","requires":{"bins":["alegra"],"env":["ALEGRA_EMAIL","ALEGRA_TOKEN"]},"install":[{"kind":"brew","formula":"jjuanrivvera/alegra-cli/alegra-cli","bins":["alegra"]},{"kind":"go","package":"github.com/jjuanrivvera/alegra-cli/cmd/alegra@latest","bins":["alegra"]}]}}
---

# Alegra CLI

Drive the [Alegra](https://alegra.com) accounting API through the `alegra`
command-line tool. This skill teaches you how and when to use it.

## Prerequisites

- The `alegra` binary must be on `PATH`. Check with `alegra version`. If missing,
  install it: `brew install jjuanrivvera/alegra-cli/alegra-cli` or
  `go install github.com/jjuanrivvera/alegra-cli/cmd/alegra@latest`.
- Credentials: either run `alegra auth login` once, or set `ALEGRA_EMAIL` and
  `ALEGRA_TOKEN` in the environment. Confirm with `alegra auth status` /
  `alegra doctor`.

## Golden rules (read before acting)

1. **Know the country (version) first.** Alegra is *one* API whose required
   fields, enums, and electronic-invoicing flow change per country. Before any
   country-specific write (contacts, items, invoices, credit/debit notes, payment
   stamping), detect the connected version with
   `alegra company get -o json | jq -r '.applicationVersion'` (or `alegra
   doctor`) and follow the matching rules below. See
   [Country = API version](#country--api-version-detect-before-you-act).
2. **Preview writes.** For any create/update/delete/emit, run it once with
   `--dry-run` first to show the exact request, then run for real after the user
   confirms. Never emit invoices in a blind retry loop — emission is **not
   idempotent**.
3. **Parse with JSON.** Add `-o json` and pipe to `jq` when you need to read
   values; the default `table` output is for humans.
4. **Resolve ids live.** Tax ids, numbering resolutions, bank accounts, regimes
   and product codes are account- and country-specific. Never hardcode them —
   look them up (`alegra taxes list`, `alegra number-templates list`, …).
5. **Invoices are append-only.** Never try to edit or delete a stamped invoice —
   issue a credit note (`alegra credit-notes create`).
6. **Count, don't dump.** Use `--count` when the user only needs a number.
7. **Confirm destructive actions** (`delete`, `void`) with the user; `delete`
   prompts unless `-y` is passed.

## Workflow: auth → detect → discover → act → verify

```bash
alegra doctor                      # 1. verify auth, plan, rate-limit, country
alegra company get -o json | jq -r '.applicationVersion'   # 2. which country/version?
alegra <resource> --help           # 3. discover a resource's actions & filters
alegra <resource> list -o json     #    inspect real data/ids
# 4. act (always --dry-run first for writes)
alegra <resource> create -f body.json --dry-run
alegra <resource> create -f body.json
alegra <resource> get <id> -o json # 5. verify the result
```

## Country = API version (detect before you act)

Alegra is a **single REST API** (`https://api.alegra.com/api/v1`, one base URL,
HTTP Basic `email:token`). There is **no `/v2`, no per-country base URL, and no
per-country version path** — the "version" of Alegra you are talking to is the
**country localization of the connected account**. The same endpoint accepts and
returns different fields, enums, and electronic-invoicing data depending on that
country, and the tax authority behind stamping changes with it.

**One token/profile = one country.** A token belongs to a single Alegra company,
which is registered in exactly one country. If the user has accounts in several
countries they switch with `--profile`/`ALEGRA_PROFILE`; treat each profile as a
**different version** and re-detect after switching.

### How to detect the connected version

The country is exposed as `applicationVersion` on the company singleton:

```bash
alegra company get -o json | jq -r '.applicationVersion'   # e.g. "colombia", "mexico"
alegra doctor                                              # prints "country: <applicationVersion>"
```

Confirmed values from the API: `colombia`, `mexico`, `argentina`, `chile`.
Other accounts return the country slug (camelCase), e.g. `costaRica`, `peru`,
`panama`, `spain`, `republicaDominicana`. **Always read the literal value from
`company get`** rather than assuming — it is the source of truth.

### The platform decides the country — you only read it

You do **not** choose an account's country from the CLI; it is fixed by the
Alegra account and reported as `applicationVersion`. **`company get` is the
source of truth**, and the API enforces the real, country-specific rules
server-side (a missing field comes back as a `400` naming it).

`alegra config set-country <country>` exists only as an *optional, offline hint*
for the CLI's client-side pre-flight validation — today it is essentially
cosmetic (it labels validation errors with a country name; it does **not** change
what the API accepts), so you can skip it. If you do set it, use the lowercased
country name (`colombia`, `costa-rica`); note that is a different string from the
camelCase `applicationVersion` (`costaRica`).

### What changes per country

| Country (applicationVersion) | Tax authority / standard | Contact id | Item taxes | Invoice essentials |
|---|---|---|---|---|
| `colombia` | DIAN · FE 2.1 (CUFE) | NIT/CC + `fiscalResponsibilities` | IVA, INC | `numberTemplate`, `paymentForm`, `operationType`, `invoiceType` |
| `mexico` | SAT · CFDI 4.0 (UUID) | RFC + `regimeObject`, `cfdiUse`, recipient ZIP | IVA, IEPS, ISR + `claveProdServ`/`claveUnidad` | `paymentMethod` PUE/PPD, `paymentForm`, global invoices, complemento de pago |
| `peru` | SUNAT (OSE/PSE) | RUC/DNI | IGV | `numberTemplate`; boletas reported as a daily summary |
| `costaRica` | Hacienda · v4.4 | CF/CJ/DIMEX/NITE | IVA | `saleCondition`, `paymentMethod`; credit sales need a **REP** (`payments stamp`) |
| `spain` | AEAT (Verifactu/SII) | NIF/CIF | IVA, IRPF, recargo equiv. | regime-driven operation types |
| `panama` | DGI (PAC) | RUC + DV | ITBMS | `invoiceType` INTERNAL_OPERATION/EXPORT/FREE_ZONE, `thirdType` |
| `argentina` / `chile` / `republicaDominicana` | AFIP CAE / SII DTE / DGII e-CF | CUIT / RUT / RNC | per-country | confirm against the live account |

Field/enum details and the per-country emission flow are in
`references/alegra-commands.md` → *Country / version*. When in doubt, inspect a
real record (`alegra <resource> list -o json`) and the catalogs at
`https://developer.alegra.com/reference/<country>.md` — they embed the full
per-country OpenAPI with every accepted enum value.

## Core commands

`alegra <resource> {list|get|create|update|delete|export|import}` plus resource
actions. Resources include: `contacts`, `items`, `invoices`, `credit-notes`,
`debit-notes`, `estimates`, `payments`, `bills`, `purchase-orders`, `taxes`,
`retentions`, `terms`, `price-lists`, `bank-accounts`, `number-templates`,
`warehouses`, `journals`, `categories`, `company`, `users`, `reports`, and more
(run `alegra --help` for the full list).

```bash
# Contacts
alegra contacts list --type client --all
alegra contacts create --set name="Acme S.A.S" \
  --set 'identification={"type":"NIT","number":"901123456"}' \
  --set 'type=["client"]' --set kindOfPerson=LEGAL_ENTITY

# Items / taxes (tax id 3 = IVA 19% in Colombia — confirm with `alegra taxes list`)
alegra items create --set name="Consultoría" --set price=500000 --set 'tax=[{"id":3}]'

# Invoices (see Electronic invoicing below for the body)
alegra invoices list --status open --since this-month
alegra invoices get 1234 -o json | jq '{number,status,total,balance}'

# Payments (allocate to invoices via a JSON body)
alegra payments list --type in --since 7d
alegra payments create -f payment.json
```

## Filtering, counting, output

```bash
alegra invoices list --status open --since last-month --until last-month
alegra invoices list --status open --count            # total via API metadata
alegra contacts list --param identification=901123456 # any raw API query param
alegra items list --all -o csv --columns id,name,price > items.csv
```

- `--since/--until`: `today`, `this-month`, `last-month`, `7d`, `3m`, `YYYY-MM-DD`.
- `--param key=value`: escape hatch for any Alegra query parameter.
- `-o table|json|yaml|csv`, `--columns a,b,c`, `--all` (every page), `--count`.

## Writing bodies

create/update accept the body three ways (combine freely):

```bash
alegra invoices create -f invoice.json          # file (best for nested documents)
alegra invoices create -f -                      # stdin
alegra contacts create --set name="X" --set 'type=["client"]'   # flat fields
```

`--set` sets top-level fields (values parsed as JSON when valid); for nested
documents (invoice `items[]`, payment allocations) use `--file`.

## Electronic invoicing (facturación electrónica)

**First detect the country** (above) — the required body and the emission flow
change per version. An invoice becomes fiscal when it has a numbering resolution
(`numberTemplate.id`) and is stamped (`stamp.generateStamp: true`), which sends
it to that country's tax authority and returns a fiscal code (CUFE in CO, UUID in
MX, …). Find resolutions with `alegra number-templates list`.

Safe flow — create a draft, review, then emit:

```bash
alegra invoices create -f invoice.json --draft   # internal draft (validated per country)
alegra invoices get <id>                          # review
alegra invoices emit <id>                         # stamp; or `emit --all` for every draft
```

`invoices emit` auto-chunks into batches of 10 and keeps a local idempotency
guard (skips already-emitted ids unless `--force`). `create` runs light
client-side pre-flight validation; the real country rules are enforced by the API
(`--no-validate` skips the local checks). `invoice.json` (Colombia):

```json
{
  "client": { "id": 12 },
  "numberTemplate": { "id": 7 },
  "date": "2026-06-06", "dueDate": "2026-06-21", "paymentForm": "CASH",
  "items": [ { "id": 5, "price": 50000, "quantity": 1, "tax": [{ "id": 3 }] } ],
  "stamp": { "generateStamp": true }
}
```

### Per-country emission flow

The skeleton above is **Colombia-shaped**. Adapt it to the connected version:

- **Colombia (DIAN):** add `operationType` (`STANDARD`/`AIU_SERVICE`) and
  `invoiceType` (`NATIONAL`/`EXPORT`) when relevant; emission returns a **CUFE**.
  Reverse with a typed credit note (`creditNoteType`, e.g. `VOID_ELECTRONIC_INVOICE`).
- **Mexico (SAT/CFDI 4.0):** the receptor needs `regimeObject`, `cfdiUse`, and a
  **ZIP that matches their Constancia**; set `paymentMethod` **`PUE`** (paid now)
  or **`PPD`** (deferred). A **PPD** invoice is only closed when you register a
  payment and stamp it — `alegra payments stamp <id>` emits the **complemento de
  pago (REP)**. For público-en-general sales, consolidate with
  `alegra global-invoices create` (auto-stamped; needs `periodicity`). Emission
  returns the **UUID**.
- **Costa Rica (Hacienda 4.4):** set `saleCondition` and `paymentMethod`. A
  **credit** sale must later issue a **REP (Recibo Electrónico de Pago)** once
  collected: `alegra payments stamp <id>` (the payment must use electronic
  numbering and be linked to the electronic invoice).
- **Peru (SUNAT):** **boletas** are reported to SUNAT as a **consolidated daily
  summary**, not acknowledged one-by-one — don't poll for an instant per-boleta
  legal status; **facturas** are emitted individually.
- **Panama (DGI):** set `invoiceType` (`INTERNAL_OPERATION`/`EXPORT`/`FREE_ZONE`)
  and the recipient `thirdType` (`FINAL_CONSUMER`/`GOVERNMENT`/`TAXPAYER`/`FOREIGN`).
- **Others (AR/CL/DR/ES):** field shapes vary (AFIP CAE, SII DTE, DGII e-CF,
  AEAT). Inspect a real document and the country catalog before scripting.

When a country needs an enum you don't know (a `cfdiUse`, a `saleCondition`, a
`fiscalResponsibilities` id), read it from the country catalog
(`https://developer.alegra.com/reference/<country>.md`) or from an existing
record — never guess.

## Bulk operations

```bash
alegra contacts import -f clients.csv \
  --map 'Name=name,NIT=identification.number' \
  --set 'identification.type=NIT' --set 'type=["client"]'
alegra invoices export --param status=open > receivables.csv
```

## Reports

```bash
alegra reports sales-by-client --from 2026-01-01 --to 2026-03-31
alegra reports sales-by-seller --from 2026-01-01 --to 2026-12-31 -o csv
```

## Errors

`alegra` prints `Error: alegra: HTTP <code>: …` with a remediation hint. Common:
`401` bad credentials (`alegra auth login`), `402` the feature isn't in the
account's plan, `403` no permission, `429` rate limit (150/min — the CLI backs
off automatically), validation errors name the missing field. On a stamping
error, run `alegra invoices get <id>` to check whether it actually emitted before
retrying.

## More

Full docs and recipes: https://jjuanrivvera.github.io/alegra-cli/ . A condensed
command cheatsheet ships alongside this skill in `references/alegra-commands.md`.
