#!/usr/bin/env bash
# Build a worktree's alegra binary into <worktree>/bin/alegra, isolated from the
# main checkout and the installed `alegra`. Embeds the branch + short sha so
# `alegra version` tells you which worktree you're running.
#
# Usage: wt build [branch]   (defaults to the current worktree)

set -e
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

WT_PATH="$(resolve_worktree "$1")" || {
  echo "${RED}Error: worktree '$1' not found (see: wt list)${NC}"
  exit 1
}

cd "$WT_PATH"
SHA="$(worktree_head "$WT_PATH")"
BR="$(worktree_branch "$WT_PATH")"
PKG="github.com/jjuanrivvera/alegra-cli/internal/version"
LDFLAGS="-X $PKG.Version=dev-$BR -X $PKG.Commit=$SHA -X $PKG.BuildDate=worktree"

echo "${BLUE}Building${NC} ${CYAN}$BR${NC} ${DIM}($SHA)${NC} → bin/alegra"
go build -ldflags "$LDFLAGS" -o bin/alegra ./cmd/alegra

echo "${GREEN}✓ Built${NC} $(binary_path "$WT_PATH")"
echo "  run:  wt run $BR -- <args>"
