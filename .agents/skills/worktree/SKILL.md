---
name: worktree
description: Manage git worktrees for parallel alegra-cli development — create a worktree per branch and build/run/test that branch's `alegra` binary in isolation, without clobbering the main checkout or the installed CLI. Use when working on multiple branches/features at once, reviewing a branch while iterating another, or delegating a risky refactor to a subagent.
---

# alegra-cli worktrees

Parallel local dev across feature branches without `git stash` or branch
ping-pong. Because alegra-cli is a Go binary (not a server), each worktree:

- Lives at `.claude/worktrees/<branch-slug>/` (gitignored).
- Compiles its own `bin/alegra`, so you can run/test one branch's behavior while
  the main checkout keeps its own build — no collisions, no rebuilds stepping on
  each other.
- Shares credentials and the Go module cache (config lives in `~/.alegra-cli` or
  `ALEGRA_*` env), so there are **no secrets to symlink** and creation is fast.

`wt build` stamps the branch + short sha into the binary, so `alegra version`
from a worktree shows e.g. `dev-fix/invoice-emit (commit 0f2f672, …)` — you
always know which build you're running.

## Installation

Call it directly (no PATH setup):

```bash
bash .agents/skills/worktree/scripts/wt <command>
```

Or symlink it once (named `wt-alegra` to avoid clashing with other repos' `wt`):

```bash
ln -sf "$(git rev-parse --show-toplevel)/.agents/skills/worktree/scripts/wt" ~/.local/bin/wt-alegra
```

The skill is stored canonically at `.agents/skills/` (the cross-agent location
used by Cursor, Codex, Gemini CLI, …); `.claude/skills` is a symlink to it so
Claude Code finds it too.

## Commands

### `wt create <branch> [--build]`

Creates a worktree at `.claude/worktrees/<slug>/` for an existing local/remote
branch (or a new branch if the name doesn't exist). `--build` also compiles it.

```bash
wt create fix/invoice-emit            # new branch
wt create develop                     # existing branch
wt create feat/reports --build        # create + compile
```

### `wt build [branch]`

Compiles `<worktree>/bin/alegra` (defaults to the current worktree). Isolated
from the main checkout's `bin/` and from the Homebrew/`go install` `alegra`.

### `wt run <branch> [--] <args...>`

Builds if needed, then runs that branch's binary with the passed args.

```bash
wt run fix/invoice-emit -- invoices list --dry-run
wt run fix/invoice-emit contacts list      # the -- is optional
```

### `wt test [branch] [go test args...]`

Runs `go test ./...` inside the worktree.

```bash
wt test fix/invoice-emit
wt test fix/invoice-emit -run TestEmit -v
```

### `wt list`

Shows every worktree with branch, HEAD sha, and build status (green dot = a
binary is built).

```
●  main                       develop                       0f2f672   built
○  fix-invoice-emit           fix/invoice-emit              a1b2c3d   not built
```

### `wt remove <branch> [--force]`

Runs `git worktree remove` (git's uncommitted-changes guard applies; `--force`
discards).

### `wt cleanup [base] [--force]`

Removes worktrees whose branch is already merged into `base` (default `develop`,
falling back to `main`).

## When to use

- Reviewing/running one branch in the terminal while iterating another.
- Delegating a multi-step implementation to a subagent without its edits
  stepping on your working tree.
- A risky refactor you want to be able to `wt remove` if it goes sideways.

## When NOT to use

- One-line fixes you'll commit in two minutes — a branch + commit is faster.
- Changes to `go.mod`/`go.sum` you want to land once — the module graph is shared
  across worktrees via the global cache, so test dependency bumps deliberately.

## Gotchas

- **Shared module cache & credentials.** All worktrees use the same Go module
  cache and the same `~/.alegra-cli` config / `ALEGRA_*` env. That's intentional
  (fast, one source of truth). For an isolated config, set a different
  `ALEGRA_PROFILE` (or `ALEGRA_CONFIG`) when running that worktree's binary.
- **Live API writes are not isolated.** `wt run … invoices create` hits the same
  Alegra account as everywhere else. Use `--dry-run`, or a separate profile/token
  per worktree, when testing writes.
- **Don't `git checkout -B <branch>` from the main repo if a worktree has that
  branch checked out.** Cross-worktree branch poisoning is a real footgun: the
  ref moves but the other worktree's files don't, so its `git status` shows a pile
  of phantom "staged" changes. Recovery (when `HEAD == origin/<branch>`):
  `cd <poisoned-worktree> && git diff --stat` (empty) then `git reset --hard HEAD`.
  Always do branch operations from the worktree that owns the branch.

## What this skill does NOT do

- No dev server / port management (alegra-cli is a CLI, not a server).
- No secret symlinking (credentials are global).
- No AI-config symlinking — worktrees are checkouts of the repo, so they already
  contain the committed `.claude/` (this skill included).
