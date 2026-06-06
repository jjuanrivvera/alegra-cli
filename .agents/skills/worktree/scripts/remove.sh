#!/usr/bin/env bash
# Remove a worktree. Git's own uncommitted-changes guard applies unless --force.
#
# Usage:
#   wt remove <branch>
#   wt remove <branch> --force

set -e
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

BRANCH="$1"
FORCE=""
if [[ "$2" == "--force" || "$2" == "-f" ]]; then
  FORCE="--force"
fi

if [[ -z "$BRANCH" ]]; then
  echo "${RED}Error: branch name required${NC}"
  echo "Usage: wt remove <branch> [--force]"
  exit 1
fi

WT_PATH="$(worktree_path "$BRANCH")"
if [[ ! -d "$WT_PATH" ]]; then
  echo "${RED}Error: worktree '$BRANCH' not found at $WT_PATH${NC}"
  exit 1
fi

echo "${BLUE}Removing worktree...${NC}"
echo "  Path: ${CYAN}$WT_PATH${NC}"

cd "$(main_repo_root)"
# shellcheck disable=SC2086 # FORCE is intentionally word-split (empty or --force)
if ! git worktree remove $FORCE "$WT_PATH"; then
  echo ""
  echo "${YELLOW}Uncommitted changes? Rerun with --force to discard them:${NC}"
  echo "  wt remove $BRANCH --force"
  exit 1
fi

echo "${GREEN}✓ Removed${NC}"
