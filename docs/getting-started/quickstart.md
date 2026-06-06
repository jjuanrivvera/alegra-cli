---
title: Quickstart
---

# Quickstart

```bash
# 1. Authenticate
alegra auth login

# 2. Explore
alegra --help
alegra contacts --help

# 3. List & filter
alegra contacts list
alegra contacts list --type client --query "acme" --all
alegra invoices list --status open --limit 30

# 4. Read a single record
alegra invoices get 12
alegra invoices get 12 -o json | jq '.total'

# 5. Create
alegra contacts create --set name="Acme S.A.S" --set 'type=["client"]'
alegra invoices create -f invoice.json

# 6. Update / delete
alegra contacts update 99 --set email="hi@acme.com"
alegra contacts delete 99            # prompts for confirmation (-y to skip)

# 7. Resource actions
alegra invoices void 12
alegra invoices email 12 --set 'emails=["client@acme.com"]'

# 8. Preview the request without sending it
alegra invoices create -f invoice.json --dry-run
```

Pipe JSON into other tools:

```bash
alegra items list --all -o json | jq '[.[] | {id, name, price}]'
```
