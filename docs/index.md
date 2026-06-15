---
title: Home
---

# Alegra CLI

A fast, scriptable command-line interface for the
[Alegra](https://www.alegra.com/) accounting API.

Manage contacts, invoices, items, payments, taxes, and the full Alegra resource
surface from your terminal — with `table`/`json`/`yaml`/`csv` output, named
profiles, a dry-run mode, and a built-in [Alegra MCP server](user-guide/mcp.md).

!!! note "Unofficial"
    This is a community tool, not affiliated with Alegra. It uses the public API
    at `https://api.alegra.com/api/v1`.

## Quick links

- [Installation](getting-started/installation.md) · [Authentication](getting-started/authentication.md) · [Quickstart](getting-started/quickstart.md)
- **[Cookbook](cookbook.md)** — copy-paste recipes for everyday tasks
- Guides: [Invoice → Cash](guides/invoice-to-cash.md) · [Expenses & Purchases](guides/expenses-and-purchases.md) · [Electronic Invoicing](guides/electronic-invoicing.md) · [Reporting & Month-End](guides/reporting-and-month-end.md) · [Automation](guides/automation.md)
- Reference: [Commands](commands/index.md) · [Errors & Rate Limits](reference/errors.md) · [FAQ](reference/faq.md)
- AI agents: [Alegra MCP server](user-guide/mcp.md) · [Alegra MCP vs. the official server](reference/vs-official-mcp.md) · [Agent skill](user-guide/agent-skill.md) · [Agent safety](user-guide/agent-safety.md)

## At a glance

```bash
alegra auth login
alegra contacts list --type client --all
alegra invoices get 12 -o json
alegra invoices create -f new-invoice.json
alegra invoices void 12
```
