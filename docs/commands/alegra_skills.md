---
title: alegra skills
---

## alegra skills

Install this CLI's AI-agent skill into Claude, Cursor, and other agents

### Synopsis

Install the alegra-cli agent skill so AI coding agents know how to
drive this CLI.

The cross-agent way (recommended) is Vercel's installer, which detects every
agent you have and is always up to date:

  npx skills add jjuanrivvera/alegra-cli

This built-in command writes the bundled skill directly, without Node:

  alegra skills install                 # Claude Code (project ./.claude/skills)
  alegra skills install --global        # Claude Code (~/.claude/skills)
  alegra skills install --agent cursor --global

### Options

```
  -h, --help   help for skills
```

### Options inherited from parent commands

```
      --base-url string             Override the API base URL (env: ALEGRA_BASE_URL)
      --columns strings             Comma-separated columns for table/csv output
      --dry-run                     Print the equivalent curl request without sending it
  -o, --output string               Output format: table, json, yaml, csv (env: ALEGRA_OUTPUT)
      --profile string              Configuration profile to use (env: ALEGRA_PROFILE)
      --requests-per-second float   Client-side rate limit (default from config)
      --show-token                  In --dry-run, do not redact the Authorization header
  -v, --verbose                     Enable verbose (debug) logging to stderr
```

### SEE ALSO

* [alegra](alegra.md)	 - Alegra accounting system CLI
* [alegra skills install](alegra_skills_install.md)	 - Write the bundled skill into an agent's skills directory
* [alegra skills path](alegra_skills_path.md)	 - Print where the skill would be installed
* [alegra skills print](alegra_skills_print.md)	 - Print the bundled SKILL.md to stdout

