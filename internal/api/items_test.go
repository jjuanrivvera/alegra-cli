package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItems_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/items", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","name":"Cuaderno","reference":"CUA-001","status":"active","type":"simple","price":[{"price":1500}]},
			{"id":2,"name":"Consultoría","status":"active","type":"simple"}
		]`))
	})

	items, err := c.Items().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "Cuaderno", items[0].Name)
	assert.Equal(t, "CUA-001", items[0].Reference)
	require.Len(t, items[0].Price, 1)
	assert.Equal(t, Money("1500"), items[0].Price[0].Price)
	// Numeric id is normalized to a string.
	assert.Equal(t, ID("2"), items[1].ID)
}

func TestItems_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/items/5", r.URL.Path)
		_, _ = w.Write([]byte(`{
			"id":5,
			"name":"Teclado",
			"category":{"id":"10","name":"Tecnología"},
			"inventory":{"unit":"unit","availableQuantity":12,"unitCost":"45000"}
		}`))
	})
	item, err := c.Items().Get(context.Background(), "5")
	require.NoError(t, err)
	assert.Equal(t, ID("5"), item.ID)
	require.NotNil(t, item.Category)
	assert.Equal(t, "Tecnología", item.Category.Name)
	require.NotNil(t, item.Inventory)
	assert.Equal(t, "unit", item.Inventory.Unit)
	assert.Equal(t, Money("45000"), item.Inventory.UnitCost)
}
