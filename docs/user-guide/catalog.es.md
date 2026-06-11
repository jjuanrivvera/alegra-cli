---
title: Catálogos de referencia
---

# Catálogos de referencia

`alegra catalog` (alias `catalogs`, `reference`) te da los datos de referencia de
Alegra por país — unidades de medida y enums de referencia (tipos de
identificación, tipos de impuesto, métodos de pago, tipos de documento,
regímenes, …) — **sin conexión y sin necesidad de iniciar sesión**. Es la forma
de encontrar los códigos válidos que debes llenar al crear ítems, contactos y
documentos electrónicos.

!!! note "De dónde vienen los datos"
    Alegra sirve estos catálogos desde su propio dataset por país y no expone
    ningún endpoint REST público para ellos, así que la CLI **incrusta** los datos
    — generados a partir de las páginas de parámetros por país que publica
    Alegra. Por eso funciona sin cuenta y sin red.

## Lista las categorías de tu país

```bash
alegra catalog          # cada categoría de tu país detectado/configurado
```

## Consulta una categoría

```bash
alegra catalog units                  # unidades de medida
alegra catalog identification-types   # tipos de ID de contacto válidos
alegra catalog units -o json
```

Los alias amigables — `units`, `identification-types`, `payment-methods`,
`invoice-types`, `regimes` — mapean a las claves de categoría subyacentes (en
español). Ejecuta `alegra catalog` para ver todas las claves disponibles para tu
país.

## Elige el país

El país se detecta automáticamente desde tu cuenta (consulta
[Configuración](configuration.md)). Anúlalo explícitamente — esto no requiere
iniciar sesión:

```bash
alegra catalog units --country mexico
alegra reference identification-types --country peru
```

Cubiertos: **Colombia, México, Costa Rica, Perú, España, Panamá.**

!!! tip
    Combina esto con la creación: busca aquí el `unit`, `identification.type` o
    código de impuesto correcto, y luego úsalo en `alegra items create` /
    `contacts create` / `invoices create`.
