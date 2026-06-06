package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWarehouses_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/warehouses", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","name":"Bodega Norte","address":"Calle principal #45","isDefault":true,"status":"active"},
			{"id":3,"name":"Bodega Sur","address":null,"isDefault":false,"status":"active"}
		]`))
	})
	items, err := c.Warehouses().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "Bodega Norte", items[0].Name)
	assert.True(t, items[0].IsDefault)
	assert.Equal(t, ID("3"), items[1].ID)
}

func TestWarehouses_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/warehouses/22", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id":"22",
			"name":"Bodega de alimentos",
			"observations":"Bodega para almacenar productos alimenticios",
			"address":"123",
			"isDefault":false,
			"status":"active",
			"costCenter":{"id":"1","code":"1","name":"CC 1","description":"","status":"active"}
		}`))
	})
	wh, err := c.Warehouses().Get(context.Background(), "22")
	require.NoError(t, err)
	require.NotNil(t, wh)
	assert.Equal(t, ID("22"), wh.ID)
	assert.Equal(t, "Bodega de alimentos", wh.Name)
	assert.Equal(t, "active", wh.Status)
	require.NotNil(t, wh.CostCenter)
	assert.Equal(t, ID("1"), wh.CostCenter.ID)
	assert.Equal(t, "CC 1", wh.CostCenter.Name)
}
