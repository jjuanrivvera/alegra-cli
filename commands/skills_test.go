package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveSkillTarget(t *testing.T) {
	home, _ := os.UserHomeDir()

	// project (default)
	p, err := resolveSkillTarget("claude", "", false)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(".claude/skills", "alegra-cli"), p)

	// global
	p, err = resolveSkillTarget("claude", "", true)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(home, ".claude/skills", "alegra-cli"), p)

	// cursor project uses .agents/skills
	p, err = resolveSkillTarget("cursor", "", false)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(".agents/skills", "alegra-cli"), p)

	// explicit dir wins
	p, err = resolveSkillTarget("claude", "/tmp/x", true)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join("/tmp/x", "alegra-cli"), p)

	// unknown agent errors
	_, err = resolveSkillTarget("nope", "", false)
	assert.Error(t, err)
}

func TestSkillFiles(t *testing.T) {
	files, err := skillFiles()
	require.NoError(t, err)
	assert.Contains(t, files, "SKILL.md")
	assert.Contains(t, files, "references/alegra-commands.md")
}
