package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCostCenters_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cost-centers", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","code":"CC-01","name":"Ventas","status":"active"},
			{"id":2,"code":"CC-02","name":"Administración","status":"inactive"}
		]`))
	})

	items, err := c.CostCenters().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "CC-01", items[0].Code)
	assert.Equal(t, "Ventas", items[0].Name)
	assert.Equal(t, "active", items[0].Status)
	// Numeric id is normalized to a string.
	assert.Equal(t, ID("2"), items[1].ID)
}

func TestCostCenters_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/cost-centers/7", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":7,"code":"CC-07","name":"Logística","description":"Centro de costos de logística","status":"active"}`))
	})
	cc, err := c.CostCenters().Get(context.Background(), "7")
	require.NoError(t, err)
	require.NotNil(t, cc)
	assert.Equal(t, ID("7"), cc.ID)
	assert.Equal(t, "Logística", cc.Name)
	assert.Equal(t, "Centro de costos de logística", cc.Description)
}
