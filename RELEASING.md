# Releasing & external setup

Everything in this repo is ready. This checklist covers the steps that must be
done **outside the codebase** (GitHub, Homebrew, secrets) to ship `alegra-cli`
as a polished CLI — installable via `go install`, Homebrew, and Docker, with CI,
releases, and a docs site.

Boxes are ordered. Commands assume you're in the repo root and authenticated
with `gh` (`gh auth login`).

---

## 1. Create the GitHub repository

- [ ] Create the repo (public):
  ```bash
  gh repo create jjuanrivvera/alegra-cli --public \
    --description "Command-line interface for the Alegra accounting API" \
    --source . --remote origin --push
  ```
- [ ] Set topics/homepage (optional):
  ```bash
  gh repo edit jjuanrivvera/alegra-cli \
    --add-topic alegra --add-topic cli --add-topic accounting --add-topic golang \
    --homepage https://jjuanrivvera.github.io/alegra-cli/
  ```
- [ ] Upload the **social preview** image (one-time, UI only — not in the API):
  Settings → General → Social preview → upload `assets/social-preview.png`
  (regenerate it from `.demo/social-card.html`; see the comment in that file).
- [ ] Create the `develop` integration branch (the CI + branch model expect it):
  ```bash
  git checkout -b develop && git push -u origin develop && git checkout main
  ```

## 2. Homebrew tap (required before the first non-prerelease tag)

GoReleaser pushes the formula to a separate tap repo. `.goreleaser.yaml` already
points at `jjuanrivvera/homebrew-alegra-cli`.

- [ ] Create the tap repo:
  ```bash
  gh repo create jjuanrivvera/homebrew-alegra-cli --public \
    --description "Homebrew tap for alegra-cli"
  ```
- [ ] Create a **Personal Access Token** with `repo` scope (covering **both** the
  Homebrew tap **and** the Scoop bucket below), and add it as a secret on the
  `alegra-cli` repo named `HOMEBREW_TAP_TOKEN` (the release workflow reads
  `secrets.HOMEBREW_TAP_TOKEN` for `brews` and `scoops`):
  ```bash
  gh secret set HOMEBREW_TAP_TOKEN --repo jjuanrivvera/alegra-cli
  ```
  > If you want to cut a release **before** setting this up, tag a prerelease
  > (e.g. `v0.1.0-rc.1`) — `brews.skip_upload: auto` skips the tap push on
  > prereleases. The first stable `vX.Y.Z` tag requires the tap + token.

### 2b. Scoop bucket (Windows)

- [ ] Create the bucket repo (the `scoops` config points at it):
  ```bash
  gh repo create jjuanrivvera/scoop-alegra-cli --public \
    --description "Scoop bucket for alegra-cli"
  ```
  The `HOMEBREW_TAP_TOKEN` above must be able to push here too.

### 2c. Container image (GHCR) & supply chain — no setup needed

The release workflow already: logs in to **GHCR** with the built-in
`GITHUB_TOKEN` and pushes `ghcr.io/jjuanrivvera/alegra-cli`; generates an **SBOM**
(syft); and **signs the checksums** keyless with **cosign** (Sigstore, via the
`id-token` permission). No extra secrets required. Verify a release:
```bash
cosign verify-blob \
  --certificate checksums.txt.pem --signature checksums.txt.sig \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  --certificate-identity-regexp 'https://github.com/jjuanrivvera/alegra-cli.*' \
  checksums.txt
```

## 3. Documentation site (GitHub Pages)

The `docs.yml` workflow runs `mkdocs gh-deploy` on pushes to `main` that touch
docs/commands. It publishes to a `gh-pages` branch.

- [ ] First run is easiest manually:
  ```bash
  pip install mkdocs-material mkdocs-git-revision-date-localized-plugin
  make docs-gen
  mkdocs gh-deploy --force        # creates & pushes gh-pages
  ```
- [ ] In **Settings → Pages**, set Source = "Deploy from a branch", Branch =
  `gh-pages` / `/ (root)`. Site will be at
  `https://jjuanrivvera.github.io/alegra-cli/`.
- [ ] (Optional) trigger via CI later: `gh workflow run docs.yml`.

## 4. First release

Releases are automated by `.github/workflows/release.yml` + GoReleaser on any
`v*` tag.

- [ ] Tag and push:
  ```bash
  git checkout main
  git tag -a v0.1.0 -m "alegra-cli v0.1.0"
  git push origin main --tags
  ```
- [ ] Watch it: `gh run watch` — GoReleaser builds linux/darwin/windows
  (amd64/arm64) binaries, creates the GitHub Release with changelog, and updates
  the Homebrew formula.

## 5. Branch protection (optional)

- [ ] Protect `main` (require PRs + green CI):
  ```bash
  gh api -X PUT repos/jjuanrivvera/alegra-cli/branches/main/protection \
    -H "Accept: application/vnd.github+json" \
    -f required_status_checks.strict=true \
    -F 'required_status_checks.contexts[]=test (ubuntu-latest)' \
    -F enforce_admins=true \
    -F required_pull_request_reviews.required_approving_review_count=0 \
    -F restrictions=null
  ```

## 6. Code coverage (Codecov)

CI already uploads `coverage.out` to Codecov from the `race` job in
`.github/workflows/ci.yml`. To light up the README badge:

- [ ] Connect the repo at <https://codecov.io/gh/jjuanrivvera/alegra-cli>
  (public repos usually work tokenless).
- [ ] If uploads are rejected, add a `CODECOV_TOKEN` secret:
  ```bash
  gh secret set CODECOV_TOKEN --repo jjuanrivvera/alegra-cli
  ```

## 7. Verify after publishing

- [ ] `go install github.com/jjuanrivvera/alegra-cli/cmd/alegra@latest` → `alegra version`
- [ ] `brew install jjuanrivvera/alegra-cli/alegra-cli` → `alegra version`
- [ ] `docker build -t alegra-cli . && docker run --rm alegra-cli version`
- [ ] Docs site loads and the command reference is present.
- [ ] CI is green on `main`; release assets attached to the GitHub Release.

---

## Notes / decisions

- **Credentials never go in CI.** Tests are unit tests (`httptest`); they need
  no Alegra token. Do **not** add `ALEGRA_TOKEN` as a repo secret.
- **Go version:** workflows pin Go `1.25`. Bump in `ci.yml`, `release.yml`,
  `docs.yml`, and `Dockerfile` together when upgrading.
- **MCP:** `alegra mcp` exposes the whole CLI to AI agents; document/announce it
  once published (`claude mcp add alegra -- alegra mcp`).
- **`.dev/`** (scratch: downloaded API docs + the generation guide + the live
  test harnesses) is gitignored and intentionally not published.
