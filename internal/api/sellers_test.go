package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSellers_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/sellers", r.URL.Path)
		// Alegra returns identification as a JSON NUMBER, not a string — the
		// flexible ID type must absorb it (verified against the live API).
		_, _ = w.Write([]byte(`[{"id":"7","name":"Ana Vendedora","identification":900111222,"status":"active"}]`))
	})
	items, err := c.Sellers().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("7"), items[0].ID)
	assert.Equal(t, "Ana Vendedora", items[0].Name)
	assert.Equal(t, ID("900111222"), items[0].Identification)
	assert.Equal(t, "active", items[0].Status)
}

func TestSellers_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/sellers/7", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":"7","name":"Ana Vendedora","identification":"V-123","status":"active"}`))
	})
	s, err := c.Sellers().Get(context.Background(), "7")
	require.NoError(t, err)
	require.NotNil(t, s)
	assert.Equal(t, ID("7"), s.ID)
	// A string identification also decodes (ID is string-or-number tolerant).
	assert.Equal(t, ID("V-123"), s.Identification)
}
