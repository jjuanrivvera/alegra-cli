#!/usr/bin/env bash
# Run the Go test suite inside a worktree. Extra args are passed to `go test`.
#
# Usage:
#   wt test [branch]                 (defaults to the current worktree)
#   wt test <branch> -run TestFoo -v

set -e
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

BRANCH="$1"
shift || true

WT_PATH="$(resolve_worktree "$BRANCH")" || {
  echo "${RED}Error: worktree '$BRANCH' not found (see: wt list)${NC}"
  exit 1
}

cd "$WT_PATH"
echo "${BLUE}Testing${NC} ${CYAN}$(worktree_branch "$WT_PATH")${NC} ${DIM}in $WT_PATH${NC}"
if [[ $# -gt 0 ]]; then
  go test "$@" ./...
else
  go test ./...
fi
