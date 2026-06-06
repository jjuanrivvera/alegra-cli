#!/usr/bin/env bash
# Remove worktrees whose branch is already merged into the base branch
# (default: develop). Skips dirty worktrees unless --force.
#
# Usage:
#   wt cleanup [base]            (base defaults to develop, then main)
#   wt cleanup [base] --force

set -e
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

BASE="$1"
FORCE=""
for a in "$@"; do
  [[ "$a" == "--force" || "$a" == "-f" ]] && FORCE="--force"
done
[[ "$BASE" == "--force" || "$BASE" == "-f" ]] && BASE=""

MAIN_ROOT="$(main_repo_root)"
cd "$MAIN_ROOT"

if [[ -z "$BASE" ]]; then
  if git show-ref --verify --quiet refs/heads/develop; then
    BASE="develop"
  else
    BASE="main"
  fi
fi

WT_DIR="$(worktrees_dir)"
[[ -d "$WT_DIR" ]] || { echo "${DIM}No worktrees.${NC}"; exit 0; }

echo "${BLUE}Cleaning worktrees merged into ${CYAN}$BASE${NC}${BLUE}...${NC}"
removed=0
shopt -s nullglob
for wt in "$WT_DIR"/*/; do
  wt="${wt%/}"
  branch="$(worktree_branch "$wt")"
  if git merge-base --is-ancestor "$branch" "$BASE" 2>/dev/null; then
    echo "  ${GREEN}merged${NC} $branch — removing"
    # shellcheck disable=SC2086
    if git worktree remove $FORCE "$wt" 2>/dev/null; then
      removed=$((removed + 1))
    else
      echo "    ${YELLOW}skipped (dirty; use --force)${NC}"
    fi
  else
    echo "  ${DIM}keep${NC}   $branch (not merged)"
  fi
done

echo "${GREEN}✓ Removed $removed worktree(s)${NC}"
