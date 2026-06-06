#!/usr/bin/env bash
# List worktrees with branch, HEAD sha, and build status. The dot is green when
# the worktree has a compiled binary, grey otherwise.
#
# Usage: wt list

set -e
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

WT_DIR="$(worktrees_dir)"

print_row() {
  printf '%s  %-26s %-30s %-9s %s\n' "$1" "$2" "$3" "$4" "$5"
}

dot_for() {
  if [[ "$(built_status "$1")" == built ]]; then
    printf '%s●%s' "$GREEN" "$NC"
  else
    printf '%s○%s' "$DIM" "$NC"
  fi
}

print_row "  " "WORKTREE" "BRANCH" "HEAD" "BUILD"

# Main checkout first.
MAIN_ROOT="$(main_repo_root)"
print_row "$(dot_for "$MAIN_ROOT")" "main" \
  "$(worktree_branch "$MAIN_ROOT")" "$(worktree_head "$MAIN_ROOT")" \
  "$(built_status "$MAIN_ROOT")"

# Each worktree under .claude/worktrees/.
if [[ -d "$WT_DIR" ]]; then
  shopt -s nullglob
  for wt in "$WT_DIR"/*/; do
    wt="${wt%/}"
    print_row "$(dot_for "$wt")" "$(basename "$wt")" \
      "$(worktree_branch "$wt")" "$(worktree_head "$wt")" \
      "$(built_status "$wt")"
  done
fi
