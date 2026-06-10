package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdditionalCharges_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/additional-charges", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","name":"Propina","percentage":10,"status":"active"},
			{"id":2,"name":"Contribución parafiscal","percentage":"0","status":"inactive"}
		]`))
	})

	items, err := c.AdditionalCharges().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, ID("1"), items[0].ID)
	assert.Equal(t, "Propina", items[0].Name)
	assert.Equal(t, Money("10"), items[0].Percentage)
	assert.Equal(t, "active", items[0].Status)
	// Numeric id is normalized to a string; numeric string amount is parsed.
	assert.Equal(t, ID("2"), items[1].ID)
	assert.Equal(t, Money("0"), items[1].Percentage)
}

func TestAdditionalCharges_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/additional-charges/5", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":5,"name":"Contribución parafiscal","code":"1","description":"Pruebas","percentage":0,"categorySalesId":5062,"status":"active"}`))
	})
	ch, err := c.AdditionalCharges().Get(context.Background(), "5")
	require.NoError(t, err)
	require.NotNil(t, ch)
	assert.Equal(t, ID("5"), ch.ID)
	assert.Equal(t, "Contribución parafiscal", ch.Name)
	assert.Equal(t, "Pruebas", ch.Description)
	assert.Equal(t, ID("5062"), ch.CategorySalesID)
	assert.Equal(t, "active", ch.Status)
}
