---
title: Autenticación
---

# Autenticación

Alegra autentica las peticiones a la API con **HTTP Basic auth**: el email de tu
cuenta como nombre de usuario y un **token de API** como contraseña. La CLI los
codifica como `Authorization: Basic base64(email:token)`. (Las apps del
marketplace pueden usar en su lugar un bearer token de OAuth mediante
`ALEGRA_BEARER_TOKEN`).

## Obtén tu token de API

1. Inicia sesión en [Alegra](https://app.alegra.com/).
2. Ve a **Configuración → Integraciones → API**.
3. Copia tu token (y el email asociado a la cuenta).

## Opción 0: `alegra init` (lo más fácil para empezar)

```bash
alegra init
```

Te pide tu email y token, los verifica, **detecta tu país automáticamente**
(que se usa para la validación previa) y guarda un perfil — todo en un solo paso.
Tanto en `init` como en `auth login`, la CLI lee `/company` y cachea el país, de
modo que el comportamiento específico de cada país funciona de forma automática.

## Opción 1: `alegra auth login`

```bash
alegra auth login
# Alegra email: you@example.com
# Alegra API token: ********
```

Tu email se guarda en `~/.alegra-cli/config.yaml`; el **token se almacena en el
keyring del sistema** (macOS Keychain, Linux Secret Service, Windows Credential
Manager) — nunca se escribe en disco en texto plano.

## Opción 2: variables de entorno

Ideal para CI y scripts:

```bash
export ALEGRA_EMAIL="you@example.com"
export ALEGRA_TOKEN="your-api-token"
```

Las variables de entorno tienen prioridad sobre el archivo de configuración y el
keyring.

## Verifica

```bash
alegra auth status
```

Esto llama a `GET /users/self` e imprime el perfil activo, la base URL y el
usuario autenticado.

## Múltiples cuentas

Usa [perfiles](../user-guide/configuration.md) para cambiar entre cuentas de Alegra.
