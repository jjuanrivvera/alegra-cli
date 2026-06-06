package api

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemCategories_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/item-categories", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","name":"Telas","description":"De la mejor calidad","status":"active"},
			{"id":2,"name":"Ropa para hombre","description":"Ropa deportiva","status":"inactive"}
		]`))
	})

	categories, err := c.ItemCategories().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, categories, 2)
	assert.Equal(t, ID("1"), categories[0].ID)
	assert.Equal(t, "Telas", categories[0].Name)
	assert.Equal(t, "active", categories[0].Status)
	// Numeric id is normalized to a string.
	assert.Equal(t, ID("2"), categories[1].ID)
	assert.Equal(t, "inactive", categories[1].Status)
}

func TestItemCategories_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/item-categories/1", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":1,"name":"Telas","description":"De la mejor calidad","status":"active","image":{"id":"img1","url":"https://x/y.png"}}`))
	})
	category, err := c.ItemCategories().Get(context.Background(), "1")
	require.NoError(t, err)
	assert.Equal(t, ID("1"), category.ID)
	assert.Equal(t, "Telas", category.Name)
	require.NotNil(t, category.Image)
	assert.Equal(t, "https://x/y.png", category.Image.URL)
}

func TestItemCategories_Create(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/item-categories", r.URL.Path)
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "Ropa de mujer", body["name"])
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":99,"name":"Ropa de mujer","status":"active"}`))
	})
	created, err := c.ItemCategories().Create(context.Background(), map[string]any{"name": "Ropa de mujer"})
	require.NoError(t, err)
	assert.Equal(t, ID("99"), created.ID)
}
