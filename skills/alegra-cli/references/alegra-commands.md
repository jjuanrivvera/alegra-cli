# Alegra CLI — command cheatsheet

Condensed reference loaded on demand by the `alegra-cli` skill. Authoritative
docs: https://jjuanrivvera.github.io/alegra-cli/

## Global flags (any command)

| Flag | Meaning |
|---|---|
| `-o, --output table\|json\|yaml\|csv` | Output format (default table) |
| `--columns a,b,c` | Columns for table/csv |
| `--profile NAME` | Config profile (env `ALEGRA_PROFILE`) |
| `--base-url URL` | Override API base URL |
| `--dry-run` | Print the request (curl) and send nothing |
| `--show-token` | Don't redact auth in `--dry-run` |
| `-v, --verbose` | Debug logging to stderr |

## List flags (every `<resource> list`)

| Flag | Meaning |
|---|---|
| `--all` | Fetch every page (Alegra caps pages at 30) |
| `--count` | Print only the total (uses API metadata) |
| `--limit N` / `--start N` | Manual pagination |
| `-q, --query TEXT` | Free-text search |
| `--since` / `--until` | Date range: `today`,`this-month`,`last-month`,`7d`,`3m`,`YYYY-MM-DD` |
| `--order-field` / `--order-direction` | Sorting |
| `--param key=value` | Any raw Alegra query parameter (repeatable) |

Transactional resources (invoices, bills, payments, credit-notes, estimates)
also have `--status`, `--client-id`, `--date`, `--date-after/--date-before`,
`--due-after/--due-before`.

## Meta commands

```bash
alegra version
alegra doctor                              # config/auth/company/rate-limit/plan
alegra auth login | status | logout
alegra config view | path | use <profile> | set-country <country>
alegra alias set <name> "<expansion>" | list | remove <name>
alegra mcp                                 # run as an MCP server for agents
alegra skills install [--global] [--agent claude|cursor|…]  # install this skill
```

## CRUD + bulk (every resource)

```bash
alegra <res> list [flags]
alegra <res> get <id>
alegra <res> create -f body.json | --set k=v | -d '<json>'   [--draft] [--no-validate] [--country X]
alegra <res> update <id> -f body.json | --set k=v
alegra <res> delete <id> [-y]
alegra <res> export [--param k=v] [--format csv|json] [--out file]
alegra <res> import -f rows.csv [--map col=field.path] [--set k=v]
```

## Resource actions

| Resource | Actions |
|---|---|
| invoices | `emit [--all] [--force]`, `void <id>`, `open <id>`, `email <id> --set 'emails=[…]'`, `stamp`, `preview` |
| credit-notes / estimates | `email <id>` |
| remissions / transportation-receipts | `void`, `open` |
| payments | `void`, `open`, `stamp` |
| bills | `close`, `comments`, `import-by-cufe` |
| purchase-orders | `void`, `email`, `comments` |
| bank-accounts | `transfer <id>` |
| journals | `balance` |
| categories | `settings`, `set-settings` |
| users | `self` |
| company | `get`, `update` (singleton) |
| reports | `sales-by-client[-totals]`, `sales-by-seller`, `income-statement`, `account-statement` |

## Common bodies

Contact (Colombia):
```json
{ "name": "Acme S.A.S", "identification": {"type":"NIT","number":"901123456"},
  "type": ["client"], "kindOfPerson": "LEGAL_ENTITY" }
```

Invoice (electronic, Colombia):
```json
{ "client": {"id":12}, "numberTemplate": {"id":7},
  "date": "2026-06-06", "dueDate": "2026-06-21", "paymentForm": "CASH",
  "items": [ {"id":5,"price":50000,"quantity":1,"tax":[{"id":3}]} ],
  "stamp": {"generateStamp": true} }
```

Income payment allocated to an invoice:
```json
{ "type": "in", "date": "2026-06-06", "bankAccount": {"id":2},
  "client": {"id":12}, "invoices": [ {"id":1234,"amount":59500} ],
  "paymentMethod": "transfer" }
```

## Notes

- Tax ids are account-specific; resolve with `alegra taxes list` (CO IVA 19% is
  commonly id 3).
- Stamped invoices are immutable — reverse with a credit note.
- `delete` prompts unless `-y`; `void` keeps the record but removes accounting
  effect.
