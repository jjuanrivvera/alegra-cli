---
title: alegra
---

## alegra

Alegra accounting system CLI

### Synopsis

alegra is a command-line interface for the Alegra accounting API
(https://developer.alegra.com).

It manages contacts, invoices, items, payments, taxes, and the full Alegra
resource surface, with table/json/yaml/csv output, profiles, and a dry-run mode.

Authenticate with:
  alegra auth login            # stores your API token in the OS keyring
or set ALEGRA_EMAIL and ALEGRA_TOKEN in the environment.

### Options

```
      --base-url string             Override the API base URL (env: ALEGRA_BASE_URL)
      --columns strings             Comma-separated columns for table/csv output
      --dry-run                     Print the equivalent curl request without sending it
  -h, --help                        help for alegra
  -o, --output string               Output format: table, json, yaml, csv (env: ALEGRA_OUTPUT)
      --profile string              Configuration profile to use (env: ALEGRA_PROFILE)
      --requests-per-second float   Client-side rate limit (default from config)
      --show-token                  In --dry-run, do not redact the Authorization header
  -v, --verbose                     Enable verbose (debug) logging to stderr
```

### SEE ALSO

* [alegra additional-charges](alegra_additional-charges.md)	 - Manage additional charges (tips and parafiscal contributions)
* [alegra auth](alegra_auth.md)	 - Manage Alegra API authentication
* [alegra bank-accounts](alegra_bank-accounts.md)	 - Manage bank accounts (bank, credit card, and cash accounts)
* [alegra bills](alegra_bills.md)	 - Manage provider bills (facturas de proveedor)
* [alegra categories](alegra_categories.md)	 - Manage chart-of-accounts accounts (cuentas contables)
* [alegra company](alegra_company.md)	 - View and update the account's company (empresa)
* [alegra config](alegra_config.md)	 - Manage alegra-cli configuration and profiles
* [alegra contacts](alegra_contacts.md)	 - Manage contacts (clients and providers)
* [alegra cost-centers](alegra_cost-centers.md)	 - Manage cost centers
* [alegra credit-notes](alegra_credit-notes.md)	 - Manage credit notes
* [alegra currencies](alegra_currencies.md)	 - Manage currencies (monedas)
* [alegra custom-fields](alegra_custom-fields.md)	 - Manage custom fields (campos adicionales)
* [alegra debit-notes](alegra_debit-notes.md)	 - Manage debit notes
* [alegra estimates](alegra_estimates.md)	 - Manage estimates (cotizaciones / sales quotes)
* [alegra global-invoices](alegra_global-invoices.md)	 - Manage global invoices (facturas globales)
* [alegra income-debit-notes](alegra_income-debit-notes.md)	 - Manage customer debit notes
* [alegra inventory-adjustment-numerations](alegra_inventory-adjustment-numerations.md)	 - Manage inventory adjustment numerations
* [alegra inventory-adjustments](alegra_inventory-adjustments.md)	 - Manage inventory adjustments (manual stock corrections)
* [alegra invoices](alegra_invoices.md)	 - Manage sales invoices
* [alegra item-categories](alegra_item-categories.md)	 - Manage item categories
* [alegra items](alegra_items.md)	 - Manage items (products and services)
* [alegra journals](alegra_journals.md)	 - Manage accounting journal entries (comprobantes contables)
* [alegra mcp](alegra_mcp.md)	 - MCP server management
* [alegra number-templates](alegra_number-templates.md)	 - Manage document numberings (numeraciones de facturación)
* [alegra payments](alegra_payments.md)	 - Manage payments (incomes and expenses)
* [alegra price-lists](alegra_price-lists.md)	 - Manage price lists
* [alegra purchase-orders](alegra_purchase-orders.md)	 - Manage purchase orders (órdenes de compra)
* [alegra reconciliations](alegra_reconciliations.md)	 - Manage bank reconciliations
* [alegra recurring-invoices](alegra_recurring-invoices.md)	 - Manage recurring invoices (facturas recurrentes)
* [alegra recurring-payments](alegra_recurring-payments.md)	 - View recurring payments (pagos recurrentes)
* [alegra remissions](alegra_remissions.md)	 - Manage remissions (delivery notes)
* [alegra reports](alegra_reports.md)	 - Read-only Alegra sales reports
* [alegra retentions](alegra_retentions.md)	 - Manage retentions (withholdings)
* [alegra taxes](alegra_taxes.md)	 - Manage taxes (e.g. IVA)
* [alegra terms](alegra_terms.md)	 - Manage payment terms (términos de pago)
* [alegra transportation-receipts](alegra_transportation-receipts.md)	 - Manage transportation receipts (documentos de traslado)
* [alegra users](alegra_users.md)	 - Manage account users
* [alegra variant-attributes](alegra_variant-attributes.md)	 - Manage item variant attributes (e.g. Color, Talla)
* [alegra version](alegra_version.md)	 - Print version information
* [alegra warehouse-transfers](alegra_warehouse-transfers.md)	 - Manage inventory transfers between warehouses
* [alegra warehouses](alegra_warehouses.md)	 - Manage inventory warehouses (bodegas)

