#!/usr/bin/env bash
# Shared helpers for the alegra-cli worktree CLI.
# Source from each subcommand: source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

# Guard against double-sourcing.
if [[ -n "$_ALEGRA_WT_LIB_LOADED" ]]; then
  return 0
fi
_ALEGRA_WT_LIB_LOADED=1

# Colors (skipped when not a TTY).
if [[ -t 1 ]]; then
  RED=$'\033[0;31m'
  GREEN=$'\033[0;32m'
  YELLOW=$'\033[1;33m'
  BLUE=$'\033[0;34m'
  CYAN=$'\033[0;36m'
  DIM=$'\033[2m'
  NC=$'\033[0m'
else
  RED=''; GREEN=''; YELLOW=''; BLUE=''; CYAN=''; DIM=''; NC=''
fi

# The main checkout root (the worktree where the repo was originally cloned).
# `git rev-parse --git-common-dir` returns the shared .git directory; its
# parent is the main worktree. Works no matter which worktree you call from.
main_repo_root() {
  local common_dir
  common_dir="$(git rev-parse --git-common-dir 2>/dev/null)" || return 1
  if [[ "$common_dir" == .git ]]; then
    pwd
    return 0
  fi
  dirname "$common_dir"
}

# Directory where this skill stores worktrees, always relative to the main repo.
worktrees_dir() {
  printf '%s\n' "$(main_repo_root)/.claude/worktrees"
}

# Resolve a worktree slug → absolute path. Slug sanitizes the branch name
# (slashes → dashes), matching the layout `wt create` uses.
worktree_path() {
  local slug="${1//\//-}"
  printf '%s\n' "$(worktrees_dir)/$slug"
}

slugify_branch() {
  printf '%s\n' "${1//\//-}"
}

# Short HEAD sha + branch name for a worktree.
worktree_head() {
  git -C "$1" rev-parse --short HEAD 2>/dev/null || echo "?"
}

worktree_branch() {
  git -C "$1" rev-parse --abbrev-ref HEAD 2>/dev/null || echo "?"
}

# Path to a worktree's isolated binary (built by `wt build`).
binary_path() {
  printf '%s\n' "$1/bin/alegra"
}

# Human-readable build status for a worktree: "built (<age>)" or "not built".
built_status() {
  local bin
  bin="$(binary_path "$1")"
  if [[ -x "$bin" ]]; then
    printf 'built\n'
  else
    printf 'not built\n'
  fi
}

# Resolve a branch/slug argument to a worktree path, defaulting to the current
# worktree when no argument is given. Prints the path or returns non-zero.
resolve_worktree() {
  local arg="$1"
  if [[ -z "$arg" ]]; then
    git rev-parse --show-toplevel 2>/dev/null || return 1
    return 0
  fi
  local path
  path="$(worktree_path "$arg")"
  if [[ -d "$path" ]]; then
    printf '%s\n' "$path"
    return 0
  fi
  return 1
}
