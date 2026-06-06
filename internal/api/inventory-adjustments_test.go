package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventoryAdjustments_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/inventory-adjustments", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"id":"1",
				"date":"2021-02-27",
				"observations":"Ajuste por pérdida de mercancía",
				"warehouse":{"id":"1","name":"Principal","isDefault":true,"status":"active"},
				"items":[{"id":"1","name":"Camisa","quantity":10,"type":"out","unitCost":12000,"reference":"LOEM9483"}],
				"costCenter":null
			}
		]`))
	})
	items, err := c.InventoryAdjustments().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "2021-02-27", items[0].Date)
	require.NotNil(t, items[0].Warehouse)
	assert.Equal(t, "Principal", items[0].Warehouse.Name)
	require.Len(t, items[0].Items, 1)
	assert.Equal(t, "out", items[0].Items[0].Type)
	assert.Equal(t, Money(10), items[0].Items[0].Quantity)
	assert.Equal(t, Money(12000), items[0].Items[0].UnitCost)
}

func TestInventoryAdjustments_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/inventory-adjustments/5", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id":"5",
			"date":"2014-05-09",
			"observations":"",
			"warehouse":{"id":"1","name":"Principal","status":"active","costCenter":{"id":"4","name":"Centro 4"}},
			"items":[{"id":"1","name":"Pantalón","quantity":1,"type":"out","unitCost":"9800","reference":null}],
			"idResolution":"15",
			"prefix":"ADJ",
			"number":"10"
		}`))
	})
	adj, err := c.InventoryAdjustments().Get(context.Background(), "5")
	require.NoError(t, err)
	assert.Equal(t, ID("5"), adj.ID)
	assert.Equal(t, "ADJ", adj.Prefix)
	assert.Equal(t, "10", adj.Number)
	assert.Equal(t, ID("15"), adj.IDResolution)
	require.NotNil(t, adj.Warehouse)
	require.NotNil(t, adj.Warehouse.CostCenter)
	assert.Equal(t, ID("4"), adj.Warehouse.CostCenter.ID)
	require.Len(t, adj.Items, 1)
	assert.Equal(t, Money(9800), adj.Items[0].UnitCost)
}
