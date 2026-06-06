---
title: Home
---

# Alegra CLI

A fast, scriptable command-line interface for the
[Alegra](https://www.alegra.com/) accounting API.

Manage contacts, invoices, items, payments, taxes, and the full Alegra resource
surface from your terminal — with `table`/`json`/`yaml`/`csv` output, named
profiles, a dry-run mode, and a built-in [MCP server](user-guide/mcp.md).

!!! note "Unofficial"
    This is a community tool, not affiliated with Alegra. It uses the public API
    at `https://api.alegra.com/api/v1`.

## Quick links

- [Installation](getting-started/installation.md)
- [Authentication](getting-started/authentication.md)
- [Quickstart](getting-started/quickstart.md)
- [Command Reference](commands/index.md)

## At a glance

```bash
alegra auth login
alegra contacts list --type client --all
alegra invoices get 12 -o json
alegra invoices create -f new-invoice.json
alegra invoices void 12
```
