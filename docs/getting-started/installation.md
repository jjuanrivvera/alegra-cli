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

## Shell completion

Release archives, the `.deb`/`.rpm`/`.apk` packages, and the Homebrew/Scoop
installs bundle completion scripts, so `alegra <Tab>` may already work. For
`go install` / from-source builds, install it yourself, e.g.:

```bash
source <(alegra completion bash)        # bash, current shell
alegra completion zsh  > "${fpath[1]}/_alegra"   # zsh
```

Completion is **data-aware**: `alegra invoices get <Tab>` suggests real invoice
IDs, `--status <Tab>` lists valid values, and `--profile <Tab>` lists your
profiles. See **[Shell Completion](../user-guide/shell-completion.md)** for the
per-shell setup and the full list of what gets completed.
