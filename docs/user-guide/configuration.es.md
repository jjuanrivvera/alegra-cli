---
title: Configuración y perfiles
---

# Configuración y perfiles

El archivo de configuración vive en `~/.alegra-cli/config.yaml` (puedes
sobrescribirlo con `ALEGRA_CONFIG`). Los tokens de la API **no** se guardan
aquí: viven en el keyring del sistema.

```yaml
defaultProfile: prod
profiles:
  prod:
    email: you@biz.com
    baseUrl: https://api.alegra.com/api/v1
  sandbox:
    email: dev@biz.com
settings:
  defaultOutputFormat: table
  requestsPerSecond: 5
  logLevel: info
```

## Gestión de perfiles

```bash
alegra config set-profile --name prod --email you@biz.com
alegra auth login --profile prod        # store the prod token
alegra config use prod                  # set the default profile
alegra config list-profiles
alegra config view                      # tokens shown as (keyring), never printed
alegra config path
```

## Elegir un perfil por comando

```bash
alegra --profile sandbox invoices list
ALEGRA_PROFILE=sandbox alegra invoices list
```

## Precedencia

Para cada opción, el orden es: **flag de línea de comandos → variable de
entorno → perfil en la configuración → valor por defecto incorporado**. Además,
las credenciales recurren al keyring del sistema cuando no se proporcionan por
entorno ni por configuración.

## Variables de entorno

| Variable | Significado |
| --- | --- |
| `ALEGRA_EMAIL` | Usuario de autenticación Basic |
| `ALEGRA_TOKEN` | Contraseña de autenticación Basic (token de la API) |
| `ALEGRA_BEARER_TOKEN` | Token bearer de OAuth |
| `ALEGRA_BASE_URL` | URL base de la API |
| `ALEGRA_PROFILE` | Perfil activo |
| `ALEGRA_OUTPUT` | Formato de salida por defecto |
| `ALEGRA_LOG_LEVEL` | Nivel de log (debug, info, warn, error) |
| `ALEGRA_CONFIG` | Ruta del archivo de configuración |

## País y validación previa

Alegra es una sola API localizada por país. `alegra init` y `alegra auth login`
**detectan automáticamente** tu país desde `/company` (`applicationVersion`:
colombia, mexico, costaRica, peru, …) y lo guardan en caché en el perfil, así
que el comportamiento específico de cada país funciona sin que hagas nada.

`create` ejecuta una **validación previa** según el país (CO/MX/PE/CR) que
detecta errores comunes antes de enviar. El país se resuelve en este orden:

1. el flag `--country` en el comando,
2. el país autodetectado del perfil,
3. la pista offline configurada con `alegra config set-country <country>`.

```bash
alegra config set-country colombia   # offline hint when no account is detected
alegra invoices create -f inv.json --no-validate   # skip the pre-flight check
```

Los mismos datos de referencia por país (unidades, tipos de identificación,
tipos de impuesto, …) están disponibles offline a través de
[`alegra catalog`](catalog.md).
