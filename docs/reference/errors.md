---
title: Error Reference
---

# Error Reference

When a request fails, `alegra` prints `Error: alegra: HTTP <code>: <message>` and
exits non-zero. Here's what each status means and what to do.

## HTTP status codes

| Code | Meaning | What to do |
|------|---------|------------|
| **400** | Bad request — the body is malformed or missing required fields. | Read the message; it usually names the field. For documents, check country requirements (resolution, identification type, régimen). Re-run with `--dry-run` to inspect what you sent. |
| **401** | Authentication failed. | Wrong email/token. Re-run `alegra auth login`, or check `ALEGRA_EMAIL`/`ALEGRA_TOKEN`. Verify with `alegra auth status`. |
| **402** | Payment required — your **plan doesn't include** this action, or the account is suspended. | Upgrade the Alegra plan, or avoid the endpoint. Common on `reconciliations` and some reports. |
| **403** | Forbidden — the user lacks permission for this action. | Check the user's role/permissions in Alegra (`alegra users self`). |
| **404** | Not found — the id doesn't exist (or the account is suspended). | Verify the id with a `list`. |
| **405** | Method not allowed for this endpoint. | Usually a CLI/usage bug — e.g. an action the resource doesn't support. |
| **429** | Too many requests — you hit **150 req/min**. | The CLI auto-retries with backoff. If scripting, slow down; prefer `--count` over `--all`. |
| **500** | Alegra server error. | Transient; the CLI retries. If it persists, try later. |
| **503** | Service unavailable (maintenance). | Retry later. |

The CLI automatically retries `429` and `5xx` with exponential backoff; `4xx`
(except `429`) are returned immediately since retrying won't help.

## Rate limits

Every response includes:

| Header | Meaning |
|--------|---------|
| `X-Rate-Limit-Limit` | Your ceiling (currently **150**/minute). |
| `X-Rate-Limit-Remaining` | Calls left in the current minute. |
| `X-Rate-Limit-Reset` | Seconds until the window resets. |

Run with `-v/--verbose` to see request-level logging.

## Stamping (electronic emission) failures

When the tax authority rejects a document, the error often carries a provider
code. The most common families (Colombia / DIAN via the e-Provider layer):

| Code(s) | Meaning | Fix |
|---------|---------|-----|
| `AEP1001` / `AEP1005` / `AEP1003` | Certificate invalid / expired / missing | Renew or upload the digital certificate in Alegra settings. |
| `AEP6008` | Not authorized to emit electronic invoices | Complete DIAN enablement (habilitación). |
| `AEP6009` / `AEP6011` | Numbering rejected (duplicate / out of range — "Rule 90") | Use a valid resolution range; don't re-emit a number. |
| `AEP6012`–`AEP6017` | Missing required field (org type, id type, name, regime, email, address) | Add the field to the payload/contact. |
| `AEP2011` | Already sent / processing | **Check status before retrying** — it may have emitted. |
| `EPR500` / `EPR503` | DIAN communication error / timeout | Transient — verify status, then retry. |

> Emission is **not idempotent**. After any stamping error, run
> `alegra invoices get <id>` to see whether it actually emitted **before** you
> retry. See [Electronic Invoicing](../guides/electronic-invoicing.md).

## Getting more detail

- `--dry-run` shows the exact request the CLI would send.
- `-o json` on a `get` shows the full server object (including any `stamp`/error
  details Alegra attached).
- `-v` enables debug logging (requests, retries, rate-limit budget).
