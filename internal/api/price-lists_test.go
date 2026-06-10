package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPriceLists_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/price-lists", r.URL.Path)
		_, _ = w.Write([]byte(`[{"id":"1","name":"General","status":"active","type":"amount"},{"id":"2","name":"Precios mayorista","status":"active","type":"percentage","percentage":"10.00"}]`))
	})
	items, err := c.PriceLists().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "amount", items[0].Type)
	assert.Equal(t, Money("10.00"), items[1].Percentage)
}

func TestPriceLists_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/price-lists/53", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":"53","name":"Inventario Básico","status":"active","type":"percentage","percentage":60}`))
	})
	pl, err := c.PriceLists().Get(context.Background(), "53")
	require.NoError(t, err)
	assert.Equal(t, ID("53"), pl.ID)
	assert.Equal(t, "percentage", pl.Type)
	assert.Equal(t, Money("60"), pl.Percentage)
}
