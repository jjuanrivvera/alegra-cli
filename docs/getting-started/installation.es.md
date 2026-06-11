---
title: Instalación
---

# Instalación

## Homebrew

```bash
brew install jjuanrivvera/alegra-cli/alegra-cli
```

## go install

```bash
go install github.com/jjuanrivvera/alegra-cli/cmd/alegra@latest
```

## Desde el código fuente

```bash
git clone https://github.com/jjuanrivvera/alegra-cli
cd alegra-cli
make build
./bin/alegra --help
```

## Docker

```bash
docker build -t alegra-cli .
docker run --rm -e ALEGRA_EMAIL -e ALEGRA_TOKEN alegra-cli contacts list
```

Verifica:

```bash
alegra version
```

## Autocompletado del shell

Los archivos comprimidos de cada release, los paquetes `.deb`/`.rpm`/`.apk` y las
instalaciones con Homebrew/Scoop incluyen los scripts de autocompletado, así que
`alegra <Tab>` quizás ya funcione. Para las instalaciones con `go install` o desde
el código fuente, instálalo tú mismo, por ejemplo:

```bash
source <(alegra completion bash)        # bash, current shell
alegra completion zsh  > "${fpath[1]}/_alegra"   # zsh
```

El autocompletado **conoce tus datos**: `alegra invoices get <Tab>` sugiere IDs
de facturas reales, `--status <Tab>` lista los valores válidos y `--profile <Tab>`
lista tus perfiles. Consulta **[Autocompletado del shell](../user-guide/shell-completion.md)**
para la configuración por shell y la lista completa de lo que se autocompleta.
