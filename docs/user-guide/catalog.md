---
title: Reference Catalogs
---

# Reference Catalogs

`alegra catalog` (aliases `catalogs`, `reference`) gives you Alegra's per-country
reference data — units of measure and reference enums (identification types, tax
types, payment methods, document types, régimenes, …) — **offline, no login
required**. It's how you find the valid codes to fill when creating items,
contacts, and electronic documents.

!!! note "Where the data comes from"
    Alegra serves these catalogs from its own per-country dataset and exposes no
    public REST endpoint for them, so the CLI **embeds** the data — generated from
    Alegra's published country parameter pages. That's why it works with no
    account and no network.

## List the categories for your country

```bash
alegra catalog          # every category for your detected/configured country
```

## Look up a category

```bash
alegra catalog units                  # units of measure
alegra catalog identification-types   # valid contact ID types
alegra catalog units -o json
```

Friendly aliases — `units`, `identification-types`, `payment-methods`,
`invoice-types`, `regimes` — map to the underlying (Spanish) category keys. Run
`alegra catalog` to see every key available for your country.

## Pick the country

The country is auto-detected from your account (see
[Configuration](configuration.md)). Override it explicitly — this needs no login:

```bash
alegra catalog units --country mexico
alegra reference identification-types --country peru
```

Covered: **Colombia, Mexico, Costa Rica, Peru, Spain, Panama.**

!!! tip
    Pair this with creation: look up the right `unit`, `identification.type`, or
    tax code here, then use it in `alegra items create` / `contacts create` /
    `invoices create`.
