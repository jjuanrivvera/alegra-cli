// Package alegracli embeds the agent-skill assets (SKILL.md + references) into
// the binary so `alegra skills install` can write them into an AI agent's skills
// directory. The same files under skills/alegra-cli/ are what
// `npx skills add jjuanrivvera/alegra-cli` and the Claude Code plugin consume.
package alegracli

import (
	"embed"
	"io/fs"
)

//go:embed skills/alegra-cli
var embedded embed.FS

// SkillFS is rooted at the skill directory, so it contains SKILL.md and
// references/ at its top level.
var SkillFS = mustSub(embedded, "skills/"+SkillName)

// SkillName is the directory the skill installs into within an agent's skills dir.
const SkillName = "alegra-cli"

func mustSub(f embed.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
