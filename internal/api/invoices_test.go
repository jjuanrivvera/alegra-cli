package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoices_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/invoices", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"id":"1","date":"2015-09-15","dueDate":"2015-10-15","status":"open","total":609,"balance":600,"client":{"id":"1","name":"Coorporación Alegrate"}},
			{"id":2,"date":"2015-11-15","status":"draft","total":"1250.50"}
		]`))
	})

	invoices, err := c.Invoices().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, invoices, 2)
	assert.Equal(t, ID("1"), invoices[0].ID)
	assert.Equal(t, "open", invoices[0].Status)
	assert.Equal(t, Money("609"), invoices[0].Total)
	require.NotNil(t, invoices[0].Client)
	assert.Equal(t, "Coorporación Alegrate", invoices[0].Client.Name)
	// Numeric id and numeric-string money are normalized.
	assert.Equal(t, ID("2"), invoices[1].ID)
	assert.Equal(t, Money("1250.50"), invoices[1].Total)
}

func TestInvoices_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/invoices/5", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":5,"status":"closed","total":100,"items":[{"id":"1","name":"Billetera","price":80,"quantity":5}]}`))
	})
	invoice, err := c.Invoices().Get(context.Background(), "5")
	require.NoError(t, err)
	assert.Equal(t, ID("5"), invoice.ID)
	assert.Equal(t, "closed", invoice.Status)
	require.Len(t, invoice.Items, 1)
	assert.Equal(t, "Billetera", invoice.Items[0].Name)
	assert.Equal(t, Money("80"), invoice.Items[0].Price)
}

func TestInvoices_Void(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/invoices/7/void", r.URL.Path)
		_, _ = w.Write([]byte(`{"id":7,"status":"void"}`))
	})
	var out Invoice
	require.NoError(t, c.Invoices().Action(context.Background(), "7", "void", nil, &out))
	assert.Equal(t, "void", out.Status)
}
