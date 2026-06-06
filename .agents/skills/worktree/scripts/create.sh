#!/usr/bin/env bash
# Create a worktree at .claude/worktrees/<slug>/ for a branch.
# Go needs no secret/env symlinks (config is global in ~/.alegra-cli and the
# module cache is shared), so creation is just `git worktree add`.
#
# Usage: wt create <branch> [--build]

set -e
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

BRANCH="$1"
DO_BUILD=""
[[ "$2" == "--build" ]] && DO_BUILD=1

if [[ -z "$BRANCH" ]]; then
  echo "${RED}Error: branch name required${NC}"
  echo "Usage: wt create <branch> [--build]"
  exit 1
fi

MAIN_ROOT="$(main_repo_root)" || {
  echo "${RED}Error: not inside the alegra-cli git repo${NC}"
  exit 1
}

SLUG="$(slugify_branch "$BRANCH")"
WT_PATH="$(worktree_path "$SLUG")"

if [[ -e "$WT_PATH" ]]; then
  echo "${RED}Error: worktree already exists at $WT_PATH${NC}"
  exit 1
fi

mkdir -p "$(worktrees_dir)"

echo "${BLUE}Creating worktree...${NC}"
echo "  Branch: ${CYAN}$BRANCH${NC}"
echo "  Path:   ${CYAN}$WT_PATH${NC}"
echo ""

cd "$MAIN_ROOT"

# Pick the right `git worktree add` form: existing local branch, remote
# tracking branch, or a brand-new branch off the current HEAD.
if git show-ref --verify --quiet "refs/heads/$BRANCH"; then
  echo "${YELLOW}Branch exists locally — checking out into worktree${NC}"
  git worktree add "$WT_PATH" "$BRANCH"
elif git show-ref --verify --quiet "refs/remotes/origin/$BRANCH"; then
  echo "${YELLOW}Branch exists on origin — creating local tracking branch${NC}"
  git worktree add --track -b "$BRANCH" "$WT_PATH" "origin/$BRANCH"
else
  echo "${YELLOW}Creating new branch + worktree${NC}"
  git worktree add -b "$BRANCH" "$WT_PATH"
fi

echo ""
echo "${GREEN}✓ Worktree ready${NC}"

if [[ -n "$DO_BUILD" ]]; then
  echo ""
  bash "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/build.sh" "$BRANCH"
else
  echo "  cd $WT_PATH"
  echo "  wt build $BRANCH   ${DIM}# compile this branch's binary in isolation${NC}"
fi
