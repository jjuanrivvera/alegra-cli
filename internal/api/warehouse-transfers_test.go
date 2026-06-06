package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWarehouseTransfers_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/warehouse-transfers", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"id":"3",
				"date":"2020-07-28",
				"observations":"",
				"origin":{"id":"1","name":"Principal","isDefault":true,"status":"active"},
				"destination":{"id":"2","name":"Bodega 1","status":"active"},
				"items":[{"id":"4","name":"Item inventariable","quantity":"10.00"}]
			}
		]`))
	})
	items, err := c.WarehouseTransfers().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("3"), items[0].ID)
	assert.Equal(t, "2020-07-28", items[0].Date)
	require.NotNil(t, items[0].Origin)
	assert.Equal(t, ID("1"), items[0].Origin.ID)
	assert.Equal(t, "Principal", items[0].Origin.Name)
	require.Len(t, items[0].Items, 1)
	assert.Equal(t, Money(10), items[0].Items[0].Quantity)
}

func TestWarehouseTransfers_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/warehouse-transfers/3", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id":"3",
			"date":"2020-07-28",
			"origin":{"id":"1","name":"Principal","costCenter":{"id":"4","name":"Centro 4"}},
			"destination":{"id":"2","name":"Bodega 1"},
			"items":[{"id":"4","name":"Item inventariable","quantity":"10.00"}]
		}`))
	})
	tr, err := c.WarehouseTransfers().Get(context.Background(), "3")
	require.NoError(t, err)
	assert.Equal(t, ID("3"), tr.ID)
	require.NotNil(t, tr.Destination)
	assert.Equal(t, "Bodega 1", tr.Destination.Name)
	require.NotNil(t, tr.Origin)
	require.NotNil(t, tr.Origin.CostCenter)
	assert.Equal(t, ID("4"), tr.Origin.CostCenter.ID)
	assert.Equal(t, "Centro 4", tr.Origin.CostCenter.Name)
}
