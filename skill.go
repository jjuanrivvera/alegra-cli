// Package alegracli embeds the agent-skill assets (SKILL.md + references) into
// the binary so `alegra skills install` can write them into an AI agent's skills
// directory. The same SKILL.md at the repo root is what `npx skills add
// jjuanrivvera/alegra-cli` and the Claude Code plugin manifest consume.
package alegracli

import "embed"

// SkillFS holds the skill files (SKILL.md and references/).
//
//go:embed SKILL.md references
var SkillFS embed.FS

// SkillName is the directory the skill installs into within an agent's skills dir.
const SkillName = "alegra-cli"
