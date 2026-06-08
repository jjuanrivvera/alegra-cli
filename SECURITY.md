# Security Policy

`alegra-cli` talks to a financial/accounting API and handles **API tokens and
credentials**, so we take security seriously.

## Supported versions

Only the latest released `vX.Y.Z` receives security fixes. Please reproduce on
the latest release (or `main`) before reporting.

## Reporting a vulnerability

**Do not open a public issue, PR, or discussion for security problems.**

Report privately through GitHub's **[Private vulnerability reporting](https://github.com/jjuanrivvera/alegra-cli/security/advisories/new)**
(repo → **Security** → **Report a vulnerability**). If that is unavailable,
contact the maintainer privately via GitHub ([@jjuanrivvera99](https://github.com/jjuanrivvera99)).

Please include:

- a description of the issue and its impact,
- steps to reproduce (a minimal command sequence or PoC),
- affected version (`alegra version`) and OS,
- any logs — **with tokens/credentials redacted**.

## What to expect

- Acknowledgement within **5 business days**.
- An initial assessment and severity within **10 business days**.
- A fix released as promptly as the severity warrants, with credit in the
  release notes (unless you prefer to remain anonymous). We follow coordinated
  disclosure: please give us a reasonable window before publishing details.

## Handling credentials safely

- Prefer `alegra auth login` — tokens are stored in the OS keyring, never written
  to the config file in plaintext.
- When sharing output or `--dry-run` curl for a bug report, **never** include a
  real `Authorization` header or token. `--dry-run` redacts it by default; do not
  pass `--show-token` in shared logs.
- Treat `ALEGRA_TOKEN` like a password: keep it out of shell history, CI logs,
  and committed files.
