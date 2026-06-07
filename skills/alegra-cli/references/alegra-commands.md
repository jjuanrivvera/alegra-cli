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

## Country / version (applicationVersion)

Alegra is **one API** (`api.alegra.com/api/v1`, Basic `email:token`) — no `/v2`,
no per-country URL. The "version" is the **country localization of the connected
account**, exposed as `applicationVersion`. One token = one country. Detect it
before any country-specific write:

```bash
alegra company get -o json | jq -r '.applicationVersion'   # colombia, mexico, costaRica, …
alegra doctor                                              # prints country + regime
```

You don't set the country — the platform (account) does, and the API enforces it.
`alegra config set-country <country>` is only an *optional, currently-cosmetic*
local validation hint (lowercased `colombia`/`costa-rica`, distinct from the
camelCase `applicationVersion`); skip it and rely on the API's `400`s.

### Per-country field & enum cheat-sheet

Resolve actual ids/enums live (`taxes list`, `number-templates list`, the
catalog at `developer.alegra.com/reference/<country>.md`). Common shapes:

| Version | Contact `identification.type` | Regime field & sample values | Item taxes | Invoice-only fields | Emission → code |
|---|---|---|---|---|---|
| `colombia` | NIT, CC, CE, TI, PP, PEP, FOREIGN_NIT (+ `fiscalResponsibilities:[5,7,12,114]`) | `regime`: SIMPLIFIED_REGIME, COMMON_REGIME | IVA (19%≈id 3), INC | `paymentForm` CASH/CREDIT, `operationType` STANDARD/AIU_SERVICE, `invoiceType` NATIONAL/EXPORT | DIAN → **CUFE** |
| `mexico` | RFC | `regimeObject`: GENERAL_REGIME_OF_MORAL_PEOPLE_LAW, REGIME_OF_TRUST (RESICO), BUSINESS_ACTIVITIES_REGIME… + `cfdiUse` G01/G03/I01-08/D01-10/CP01/S01 | IVA, IEPS, ISR + item `claveProdServ`/`claveUnidad` | `paymentMethod` PUE/PPD, `paymentForm` (01/03/99…), receptor ZIP must match | SAT/CFDI 4.0 → **UUID** |
| `costaRica` | CF, CJ, DIMEX, NITE, PE | `regime`: SIMPLIFIED_REGIME, TRADITIONAL_REGIME, PHYSICAL_PERSON | IVA | `saleCondition` (CASH, CREDIT, CONSIGNATION…), `paymentMethod` (CASH, CARD, CHECK, TRANSFER…) | Hacienda 4.4 → REP via `payments stamp` |
| `peru` | RUC, DNI | — | IGV | `numberTemplate`; boletas = daily summary | SUNAT (OSE/PSE) |
| `spain` | NIF, CIF | `regime`: GENERAL_REGIME, SIMPLIFIED_REGIME (+ special regimes) | IVA, IRPF, recargo equiv. | regime-driven operation types | AEAT (Verifactu/SII) |
| `panama` | RUC (+ DV) | — | ITBMS | `invoiceType` INTERNAL_OPERATION/EXPORT/FREE_ZONE, `thirdType` FINAL_CONSUMER/GOVERNMENT/TAXPAYER/FOREIGN | DGI (PAC) |
| `argentina`,`chile`,`republicaDominicana` | CUIT / RUT / RNC | per-country | per-country | confirm on the live account | AFIP CAE / SII DTE / DGII e-CF |

Per-country emission specials: **MX** PPD invoices and **CR** credit sales need a
stamped payment — `alegra payments stamp <id>` (MX complemento de pago / CR REP);
**MX** público-en-general → `alegra global-invoices create` (auto-stamped,
`periodicity`).

## Common bodies

Contact (Colombia):
```json
{ "name": "Acme S.A.S", "identification": {"type":"NIT","number":"901123456"},
  "type": ["client"], "kindOfPerson": "LEGAL_ENTITY" }
```

Contact (Mexico — RFC + SAT regime/uso):
```json
{ "name": "Cliente SA de CV", "identification": "XAXX010101000",
  "type": ["client"], "regimeObject": "GENERAL_REGIME_OF_MORAL_PEOPLE_LAW",
  "cfdiUse": "G03", "address": {"zipCode": "06000"} }
```

Invoice (electronic, Colombia):
```json
{ "client": {"id":12}, "numberTemplate": {"id":7},
  "date": "2026-06-06", "dueDate": "2026-06-21", "paymentForm": "CASH",
  "items": [ {"id":5,"price":50000,"quantity":1,"tax":[{"id":3}]} ],
  "stamp": {"generateStamp": true} }
```

Invoice (electronic, Mexico CFDI):
```json
{ "client": {"id":12}, "numberTemplate": {"id":7},
  "date": "2026-06-06", "paymentMethod": "PUE", "paymentForm": "03",
  "items": [ {"id":5,"price":500,"quantity":1,"tax":[{"id":1}]} ],
  "stamp": {"generateStamp": true} }
```

Income payment allocated to an invoice:
```json
{ "type": "in", "date": "2026-06-06", "bankAccount": {"id":2},
  "client": {"id":12}, "invoices": [ {"id":1234,"amount":59500} ],
  "paymentMethod": "transfer" }
```

## Notes

- **Detect the country first** (`alegra company get`/`doctor`) — required fields,
  enums, and the emission flow all change per `applicationVersion`. One token =
  one country; re-detect after switching `--profile`.
- Tax ids are account-specific; resolve with `alegra taxes list` (CO IVA 19% is
  commonly id 3).
- Stamped invoices are immutable — reverse with a credit note.
- `delete` prompts unless `-y`; `void` keeps the record but removes accounting
  effect.
