---
title: Authentication
---

# Authentication

Alegra authenticates API requests with **HTTP Basic auth**: your account email
as the username and an **API token** as the password. The CLI encodes these as
`Authorization: Basic base64(email:token)`.

## Get your API token

1. Sign in to [Alegra](https://app.alegra.com/).
2. Go to **Configuración → Integraciones → API**.
3. Copy your token (and the email associated with the account).

## Option 1: `alegra auth login` (recommended)

```bash
alegra auth login
# Alegra email: you@example.com
# Alegra API token: ********
```

Your email is saved to `~/.alegra-cli/config.yaml`; the **token is stored in
your operating system keyring** (macOS Keychain, Linux Secret Service, Windows
Credential Manager) — never written to disk in plaintext.

## Option 2: environment variables

Ideal for CI and scripts:

```bash
export ALEGRA_EMAIL="you@example.com"
export ALEGRA_TOKEN="your-api-token"
```

Environment variables take precedence over the config file and keyring.

## Verify

```bash
alegra auth status
```

This calls `GET /users/self` and prints the active profile, base URL, and the
authenticated user.

## Multiple accounts

Use [profiles](../user-guide/configuration.md) to switch between Alegra accounts.
