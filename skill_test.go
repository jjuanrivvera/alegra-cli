package alegracli

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSkillFS verifies the embedded skill bundle is rooted at the skill
// directory and ships the expected assets (exercises the package init/mustSub).
func TestSkillFS(t *testing.T) {
	require.NotNil(t, SkillFS)
	assert.Equal(t, "alegra-cli", SkillName)

	data, err := fs.ReadFile(SkillFS, "SKILL.md")
	require.NoError(t, err)
	assert.Contains(t, string(data), "name: alegra-cli")

	_, err = fs.Stat(SkillFS, "references/alegra-commands.md")
	require.NoError(t, err)
}
