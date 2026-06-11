---
title: Inicio rápido
---

# Inicio rápido

```bash
# 1. Configuración guiada: pide email + token, los verifica, detecta tu
#    país y guarda un perfil (o usa `alegra auth login` para solo iniciar sesión)
alegra init

# 2. Explora
alegra --help
alegra contacts --help

# 3. Lista y filtra
alegra contacts list
alegra contacts list --type client --query "acme" --all
alegra invoices list --status open --limit 30
alegra invoices list --status open --since this-month   # rangos de fecha naturales

# 4. Lee un registro
alegra invoices get 12
alegra invoices get 12 -o json | jq '.total'

# 5. Crea
alegra contacts create --set name="Acme S.A.S" --set 'type=["client"]'
alegra invoices create -f invoice.json

# 6. Actualiza / elimina
alegra contacts update 99 --set email="hi@acme.com"
alegra contacts delete 99            # pide confirmación (-y para omitirla)

# 7. Acciones del recurso
alegra invoices void 12
alegra invoices email 12 --set 'emails=["client@acme.com"]'

# 8. Previsualiza la petición sin enviarla
alegra invoices create -f invoice.json --dry-run
```

Pasa JSON por un pipe hacia otras herramientas:

```bash
alegra items list --all -o json | jq '[.[] | {id, name, price}]'
```
