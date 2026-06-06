package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalInvoices_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/global-invoices", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"125","date":"2018-07-25","status":"open","total":1000,
			 "numberTemplate":{"id":"1","number":"422"},
			 "saleTickets":[{"id":"378","status":"open","total":1000,"balance":1000,
			   "client":{"id":"6","name":"Ventas al público general"}}]},
			{"id":124,"date":"2018-07-25","status":"open","total":"1500"}
		]`))
	})
	items, err := c.GlobalInvoices().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, ID("125"), items[0].ID)
	assert.Equal(t, "open", items[0].Status)
	assert.Equal(t, Money(1000), items[0].Total)
	require.Len(t, items[0].SaleTickets, 1)
	assert.Equal(t, ID("6"), items[0].SaleTickets[0].Client.ID)
	assert.Equal(t, ID("124"), items[1].ID)
	assert.Equal(t, Money(1500), items[1].Total)
}

func TestGlobalInvoices_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/global-invoices/125", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id":"125","date":"2018-07-25","status":"open","total":1000,
			"stamp":{"uuid":"1579E004-A899-47CF-915C-674ED4B54D3E","stampDate":"2018-07-25 17:33:09"}
		}`))
	})
	gi, err := c.GlobalInvoices().Get(context.Background(), "125")
	require.NoError(t, err)
	require.NotNil(t, gi)
	assert.Equal(t, ID("125"), gi.ID)
	assert.Equal(t, Money(1000), gi.Total)
	require.NotNil(t, gi.Stamp)
	assert.Equal(t, "1579E004-A899-47CF-915C-674ED4B54D3E", gi.Stamp.UUID)
}
