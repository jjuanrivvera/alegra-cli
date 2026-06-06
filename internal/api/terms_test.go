package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerms_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/terms", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","name":"De contado","days":0},
			{"id":2,"name":"8 días","days":8}
		]`))
	})
	items, err := c.Terms().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "De contado", items[0].Name)
	assert.Equal(t, Int(0), items[0].Days)
	assert.Equal(t, ID("2"), items[1].ID)
	assert.Equal(t, Int(8), items[1].Days)
}

func TestTerms_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/terms/1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"1","name":"De contado","days":0}`))
	})
	term, err := c.Terms().Get(context.Background(), "1")
	require.NoError(t, err)
	require.NotNil(t, term)
	assert.Equal(t, ID("1"), term.ID)
	assert.Equal(t, "De contado", term.Name)
	assert.Equal(t, Int(0), term.Days)
}
