package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEstimates_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/estimates", r.URL.Path)
		assert.Equal(t, "30", r.URL.Query().Get("limit"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"8","number":"8","date":"2016-05-16","dueDate":"2016-06-15","total":1720400,"client":{"id":"1","name":"Mi contacto"}},
			{"id":9,"number":"15","date":"2016-05-16","total":"610200","seller":{"id":"6","name":"Alejandro"}}
		]`))
	})
	items, err := c.Estimates().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, ID("8"), items[0].ID)
	assert.Equal(t, Money(1720400), items[0].Total)
	require.NotNil(t, items[0].Client)
	assert.Equal(t, "Mi contacto", items[0].Client.Name)
	assert.Equal(t, ID("9"), items[1].ID)
	assert.Equal(t, Money(610200), items[1].Total)
}

func TestEstimates_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/estimates/8", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id":"8",
			"number":"8",
			"date":"2016-05-16",
			"dueDate":"2016-06-15",
			"total":1720400,
			"client":{"id":"1","name":"Mi contacto"},
			"items":[
				{"id":"2","name":"Consultoría","price":500000,"quantity":"3.00","total":1500000,"tax":[]},
				{"id":"1","name":"Servicio","price":95000,"quantity":"2.00","total":220400,"tax":[{"id":"3","name":"IVA"}]}
			]
		}`))
	})
	est, err := c.Estimates().Get(context.Background(), "8")
	require.NoError(t, err)
	require.NotNil(t, est)
	assert.Equal(t, ID("8"), est.ID)
	assert.Equal(t, "8", est.Number)
	assert.Equal(t, Money(1720400), est.Total)
	require.Len(t, est.Items, 2)
	assert.Equal(t, "Consultoría", est.Items[0].Name)
	assert.Equal(t, Money(1500000), est.Items[0].Total)
	require.Len(t, est.Items[1].Tax, 1)
	assert.Equal(t, "IVA", est.Items[1].Tax[0].Name)
}
