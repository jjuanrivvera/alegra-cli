package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	assert.Equal(t, [][]string{{"1", "2"}, {"3"}}, chunk([]string{"1", "2", "3"}, 2))
	assert.Equal(t, [][]string{{"a"}}, chunk([]string{"a"}, 10))
	assert.Empty(t, chunk(nil, 10))

	// 23 ids → 3 batches of 10/10/3 (the stamp cap)
	ids := make([]string, 23)
	for i := range ids {
		ids[i] = "x"
	}
	got := chunk(ids, maxStampBatch)
	assert.Len(t, got, 3)
	assert.Len(t, got[0], 10)
	assert.Len(t, got[2], 3)
}

func TestFilterEmitted(t *testing.T) {
	cache := map[string]bool{emitKey("prod", "1"): true}
	todo, skipped := filterEmitted([]string{"1", "2", "3"}, cache, "prod", false)
	assert.Equal(t, []string{"2", "3"}, todo)
	assert.Equal(t, []string{"1"}, skipped)

	// force re-emits everything
	todo, skipped = filterEmitted([]string{"1", "2"}, cache, "prod", true)
	assert.Equal(t, []string{"1", "2"}, todo)
	assert.Empty(t, skipped)

	// cache is profile-scoped
	todo, _ = filterEmitted([]string{"1"}, cache, "sandbox", false)
	assert.Equal(t, []string{"1"}, todo)
}
