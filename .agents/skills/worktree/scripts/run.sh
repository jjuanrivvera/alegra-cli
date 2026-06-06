#!/usr/bin/env bash
# Run a worktree's binary, building it first if needed. Everything after `--`
# (or after the branch) is passed through to `alegra`.
#
# Usage:
#   wt run <branch> -- <args...>
#   wt run <branch> contacts list      ( -- is optional)

set -e
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

BRANCH="$1"
shift || true
# Allow an optional `--` separator before the passthrough args.
[[ "$1" == "--" ]] && shift || true

WT_PATH="$(resolve_worktree "$BRANCH")" || {
  echo "${RED}Error: worktree '$BRANCH' not found (see: wt list)${NC}"
  exit 1
}

BIN="$(binary_path "$WT_PATH")"
if [[ ! -x "$BIN" ]]; then
  echo "${DIM}binary missing — building first${NC}"
  bash "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/build.sh" "$BRANCH"
fi

exec "$BIN" "$@"
