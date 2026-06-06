package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariantAttributes_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/variant-attributes", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","name":"Color","status":"active","options":[{"id":"1","value":"Rojo"},{"id":"2","value":"Verde"}]},
			{"id":2,"name":"Talla","status":"active","options":[{"id":3,"value":"S"}]}
		]`))
	})

	items, err := c.VariantAttributes().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "Color", items[0].Name)
	require.Len(t, items[0].Options, 2)
	assert.Equal(t, "Rojo", items[0].Options[0].Value)
	// Numeric ids are normalized to strings.
	assert.Equal(t, ID("2"), items[1].ID)
	assert.Equal(t, ID("3"), items[1].Options[0].ID)
}

func TestVariantAttributes_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/variant-attributes/1", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":1,"name":"Color","status":"active","options":[{"id":"1","value":"Rojo"}]}`))
	})
	item, err := c.VariantAttributes().Get(context.Background(), "1")
	require.NoError(t, err)
	assert.Equal(t, ID("1"), item.ID)
	assert.Equal(t, "Color", item.Name)
	assert.Equal(t, "active", item.Status)
	require.Len(t, item.Options, 1)
	assert.Equal(t, "Rojo", item.Options[0].Value)
}
