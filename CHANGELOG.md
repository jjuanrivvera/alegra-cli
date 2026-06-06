# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/) and this project adheres to
[Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added
- Initial release of `alegra-cli`.
- Full Alegra v1 resource surface with `list`/`get`/`create`/`update`/`delete`
  plus resource-specific actions (void, open, email, stamp, transfer, comments,
  close, import-by-cufe, …).
- Generic typed API client (`Resource[T]`) with HTTP Basic auth, exponential
  backoff retries, adaptive client-side rate limiting, and offset pagination
  (`start`/`limit`, auto `--all`).
- `table`, `json`, `yaml`, and `csv` output with `--columns` selection.
- Named profiles, environment-variable overrides, and OS keyring token storage
  (`alegra auth login`).
- `--dry-run` mode that prints the equivalent `curl` request.
- Built-in MCP server (`alegra mcp`) exposing the command tree to AI agents.
- MkDocs Material documentation site and GoReleaser-based release pipeline.
