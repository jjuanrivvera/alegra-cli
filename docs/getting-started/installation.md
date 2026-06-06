---
title: Installation
---

# Installation

## Homebrew

```bash
brew install jjuanrivvera/alegra-cli/alegra-cli
```

## go install

```bash
go install github.com/jjuanrivvera/alegra-cli/cmd/alegra@latest
```

## From source

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

Verify:

```bash
alegra version
```
