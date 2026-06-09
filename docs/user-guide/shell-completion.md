---
title: Shell Completion
---

# Shell Completion

`alegra` ships rich shell completion for **bash**, **zsh**, **fish**, and
**PowerShell**. Beyond completing command and flag names, it talks to the API to
complete *your actual data* — press <kbd>Tab</kbd> after `invoices get` and you
get real invoice IDs, each annotated with its status.

!!! tip "Already installed?"
    The release archives, the `.deb`/`.rpm`/`.apk` packages, and the
    Homebrew/Scoop installs all bundle completion scripts. If you installed that
    way, completion may already work — open a new shell and try
    `alegra <Tab>`. The sections below are for `go install` / from-source
    builds, or to (re)install manually.

## What gets completed

| You type | <kbd>Tab</kbd> offers |
| --- | --- |
| `alegra <Tab>` | every resource and command (`invoices`, `contacts`, `auth`, …) |
| `alegra invoices <Tab>` | `list`, `get`, `create`, `update`, `delete`, `export`, `void`, … |
| `alegra invoices get <Tab>` | **live invoice IDs**, labelled with name/number/status |
| `alegra invoices delete 10<Tab>` | live IDs that start with `10` |
| `alegra invoices void <Tab>` | live IDs (custom actions complete IDs too) |
| `alegra invoices list --status <Tab>` | `open`, `closed`, `draft`, `void` |
| `alegra invoices list --order-field <Tab>` | the resource's sortable fields |
| `alegra contacts get --columns <Tab>` | that resource's columns (comma-aware) |
| `alegra -o <Tab>` | `table`, `json`, `yaml`, `csv` |
| `alegra --profile <Tab>` | your configured profiles, labelled with their email |

The data-aware completers (IDs, columns, profiles, enums) make the CLI's power
discoverable without memorizing flags or reading docs.

!!! note "How live completion behaves"
    ID completion makes a single, short (one-page) API call using your active
    profile. It is **best-effort and never blocks your shell**: with no
    credentials, offline, or on a slow network it simply returns no suggestions
    (a 4-second timeout caps the wait) and you can type the ID by hand. It
    respects `--profile`, so `alegra --profile sandbox invoices get <Tab>`
    completes against your sandbox account.

## Install

=== "bash"

    Completion needs the `bash-completion` package.

    ```bash
    # Current shell only (try it now):
    source <(alegra completion bash)

    # Persist — Linux:
    alegra completion bash | sudo tee /etc/bash_completion.d/alegra >/dev/null

    # Persist — macOS (Homebrew bash-completion@2):
    alegra completion bash > "$(brew --prefix)/etc/bash_completion.d/alegra"
    ```

=== "zsh"

    Ensure completion is enabled once (in `~/.zshrc`):

    ```bash
    echo "autoload -U compinit; compinit" >> ~/.zshrc
    ```

    Then install the script onto your `$fpath`:

    ```bash
    # A common location; create it if needed and ensure it's on your fpath.
    alegra completion zsh > "${fpath[1]}/_alegra"
    ```

    Start a new shell. With Oh My Zsh, write to
    `~/.oh-my-zsh/completions/_alegra` instead.

=== "fish"

    ```bash
    # Current shell only:
    alegra completion fish | source

    # Persist:
    alegra completion fish > ~/.config/fish/completions/alegra.fish
    ```

=== "PowerShell"

    ```powershell
    # Current shell only:
    alegra completion powershell | Out-String | Invoke-Expression

    # Persist by appending to your profile:
    alegra completion powershell >> $PROFILE
    ```

Run `alegra completion <shell> --help` for the authoritative, shell-specific
instructions.

## Verify

Open a **new** shell and try:

```bash
alegra invoices <Tab><Tab>     # subcommands
alegra -o <Tab>                # table / json / yaml / csv
alegra invoices get <Tab>      # live IDs (requires a logged-in profile)
```

If nothing happens, confirm completion is enabled for your shell (bash needs
`bash-completion`; zsh needs `compinit`) and that you started a fresh session.

## Troubleshooting

- **No suggestions after `get`/`delete`** — these are live completions. Check you
  have a working profile (`alegra auth status`); without credentials they yield
  nothing by design.
- **zsh: `command not found: compdef`** — add `autoload -U compinit; compinit` to
  `~/.zshrc` and restart the shell.
- **Stale completions after upgrading** — regenerate the script (the commands
  above) and open a new shell.
