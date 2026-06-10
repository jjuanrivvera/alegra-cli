package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecurringInvoices_List(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/recurring-invoices", r.URL.Path)
		_, _ = w.Write([]byte(`[{"id":"01G77Z97X66GZT7320T9HC6EBA","repeatEvery":1,"total":15000,"client":{"id":"1","name":"Coorporación Alegrate"}}]`))
	})
	items, err := c.RecurringInvoices().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ID("01G77Z97X66GZT7320T9HC6EBA"), items[0].ID)
	assert.Equal(t, Money("15000"), items[0].Total)
	require.NotNil(t, items[0].Client)
	assert.Equal(t, "Coorporación Alegrate", items[0].Client.Name)
}

func TestRecurringInvoices_Get(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/recurring-invoices/01G77Z97X66GZT7320T9HC6EBA", r.URL.Path)
		_, _ = w.Write([]byte(`{
			"id":"01G77Z97X66GZT7320T9HC6EBA",
			"startDate":"2022-07-05",
			"repeatEvery":1,
			"currency":{"code":"USD","symbol":"$"},
			"numberTemplate":{"id":"1","documentType":"invoice"},
			"total":15000,
			"items":[{"id":"01GMDYFTF7XX2SFKSMV32KQ91J","itemId":"1","name":"Zapatos","price":5000,"quantity":1,"total":5000,"taxes":[]}]
		}`))
	})
	ri, err := c.RecurringInvoices().Get(context.Background(), "01G77Z97X66GZT7320T9HC6EBA")
	require.NoError(t, err)
	require.NotNil(t, ri)
	assert.Equal(t, ID("01G77Z97X66GZT7320T9HC6EBA"), ri.ID)
	assert.Equal(t, "2022-07-05", ri.StartDate)
	require.NotNil(t, ri.Currency)
	assert.Equal(t, "USD", ri.Currency.Code)
	require.NotNil(t, ri.NumberTemplate)
	assert.Equal(t, "invoice", ri.NumberTemplate.DocumentType)
	require.Len(t, ri.Items, 1)
	assert.Equal(t, "Zapatos", ri.Items[0].Name)
	assert.Equal(t, Money("5000"), ri.Items[0].Price)
}
