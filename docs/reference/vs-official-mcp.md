# alegra-cli vs. the official Alegra MCP server

Alegra publishes a hosted **MCP server** (`mcp.alegra.com`) so AI agents can call the
accounting API over the Model Context Protocol. `alegra-cli` also ships an MCP server
(`alegra mcp`) in addition to its terminal interface. This page compares the two so you
can pick the right tool for your use case.

!!! note "Snapshot"
    This comparison reflects the official MCP server's publicly documented tool set and
    `alegra-cli` **v0.6.1**, verified **2026-06-09**. Alegra may add tools over time;
    treat coverage figures as a point-in-time snapshot, not a permanent claim.

## At a glance

| | Official Alegra MCP | alegra-cli |
| --- | --- | --- |
| Interface | MCP only (hosted, HTTP/SSE) | Terminal **and** MCP (`alegra mcp`, local/stdio) |
| Install | None — hosted by Alegra | Install a binary (Homebrew, Scoop, Docker, …) |
| Auth | Alegra token | OS keyring + profiles + env vars |
| Output for humans | — | table / JSON / YAML / CSV |
| Resource coverage | ~24 resource areas | 37 resources |
| Invoice emission (stamp/void/email) | Not exposed | Yes |
| Offline-friendly | Depends on Alegra's hosted service | Runs locally against the API with your token |

## Resource coverage

Both wrap the same Alegra v1 API, but expose different slices of it.

### Covered by both

Contacts, items (and item categories, variant attributes, warehouses, warehouse
transfers, price lists, inventory adjustments and their numerations, custom fields),
bank accounts, reconciliations, journals, cost centers, taxes, retentions, currencies,
sellers, invoices (read/write), bills, supplier debit notes, purchase orders, payments,
and sales reports.

### Only in alegra-cli

The CLI additionally covers several sales-side and fiscal documents that the official
MCP does not expose:

- Estimates (cotizaciones)
- Credit notes
- Customer debit notes
- Remissions (delivery notes)
- Transportation receipts
- Global invoices (CFDI)
- Recurring invoices and recurring payments
- Document numberings (number-templates)
- Payment terms
- Additional charges
- Webhook subscriptions

### Only in the official MCP

The official MCP exposes a few capabilities the CLI does not (yet):

- **Item stock queries** — `get_item_stock` / `get_item_stock_summary`.
- **Richer bill sub-actions** — applying advances, uploading/deleting attachments,
  editing/deleting comments, and updating perceptions/retentions on a bill.
- **Catalog/reference lookups** — bank lists, product keys, units of measure, and
  per-resource reference enums.
- **Support Center** — help-desk tickets (out of scope for an accounting CLI).

These are tracked for `alegra-cli` where they make sense.

## Capability differences that matter

### Electronic invoicing

The official MCP's invoice tools cover create, read, update, and delete. They do **not**
expose `stamp`, `void`, `email`, `open`, or `preview`, so emitting an electronic invoice
to DIAN/SAT, voiding one, or emailing it is done through `alegra-cli` (`alegra invoices
stamp|void|email …`) or the REST API directly. If electronic emission is part of your
workflow, the CLI covers the full cycle.

### Payments

The official MCP splits payments into incoming (`incomePayments`) and outgoing
(`transaction-out`) tool families. `alegra-cli` unifies them under one `payments`
resource with a `type` flag, plus `stamp`/`void`/`open` actions.

### Reports

The two report sets are complementary:

| Report | Official MCP | alegra-cli |
| --- | --- | --- |
| Sales by client / totals | ✓ | ✓ |
| Sales by seller | ✓ | ✓ |
| Sales by seller totals | ✓ | — |
| General sales documents / totals | ✓ | — |
| Clients with items | ✓ | — |
| Account statement | — | ✓ |
| Income statement | — | ✓ |

## Which should I use?

- **Use the official MCP** if you want a zero-install, Alegra-hosted endpoint for an
  agent and your work lives in items, expenses, and read queries — and you don't need to
  emit electronic invoices.
- **Use `alegra-cli`** if you want one tool for both humans (terminal) and agents (its
  MCP server), broader resource coverage, electronic-invoice emission, dry-run previews,
  CSV import/export, and local control over credentials.

Because `alegra-cli`'s MCP server is generated from its command tree, its agent-facing
coverage matches the CLI's: any resource or action you can run in the terminal is also
available as an MCP tool. See [MCP Server](../user-guide/mcp.md) to enable it.
