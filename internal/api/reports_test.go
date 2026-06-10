package api

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReports_SalesDocuments(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/reports/sales-documents", r.URL.Path)
		assert.Equal(t, "2024-01-01", r.URL.Query().Get("from"))
		assert.Equal(t, "2024-01-31", r.URL.Query().Get("to"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[
			{"id":1,"documentNumber":"FV-001","documentType":"invoice","total":1500,"date":"2024-01-15","status":"closed"}
		]}`))
	})

	q := url.Values{}
	q.Set("from", "2024-01-01")
	q.Set("to", "2024-01-31")
	var out ReportEnvelope
	err := c.GetInto(context.Background(), "reports/sales-documents", q, &out)
	require.NoError(t, err)
	require.Len(t, out.Data, 1)
	assert.Equal(t, ID("1"), out.Data[0].ID)
	assert.Equal(t, "FV-001", out.Data[0].DocumentNumber)
	assert.Equal(t, Money("1500"), out.Data[0].Total)
}

func TestReports_SalesTotals(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/reports/sales-totals", r.URL.Path)
		assert.Equal(t, "month", r.URL.Query().Get("groupBy"))
		_, _ = w.Write([]byte(`{"data":[
			{"period":"2024-01","beforeTaxes":10000,"tax":1900,"total":11900,"discount":500,"creditNote":200}
		]}`))
	})

	q := url.Values{}
	q.Set("groupBy", "month")
	var out ReportEnvelope
	err := c.GetInto(context.Background(), "reports/sales-totals", q, &out)
	require.NoError(t, err)
	require.Len(t, out.Data, 1)
	assert.Equal(t, "2024-01", out.Data[0].Period)
	assert.Equal(t, Money("11900"), out.Data[0].Total)
	assert.Equal(t, Money("1900"), out.Data[0].Tax)
}

func TestReports_SalesByClient(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/reports/sales-by-client", r.URL.Path)
		_, _ = w.Write([]byte(`{"data":[
			{"clientName":"Empresa ABC S.A.S","totalDocuments":8,"subTotal":15000,"total":17850}
		]}`))
	})

	var out ReportEnvelope
	err := c.GetInto(context.Background(), "reports/sales-by-client", url.Values{}, &out)
	require.NoError(t, err)
	require.Len(t, out.Data, 1)
	assert.Equal(t, "Empresa ABC S.A.S", out.Data[0].ClientName)
	assert.Equal(t, Int(8), out.Data[0].TotalDocuments)
	assert.Equal(t, Money("15000"), out.Data[0].SubTotal)
}

func TestReports_SalesBySeller(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/reports/sales-by-seller", r.URL.Path)
		_, _ = w.Write([]byte(`{"data":[
			{"sellerName":"Juan Pérez","totalDocuments":15,"totalPayed":5000,"beforeTaxes":8000,"total":9520}
		]}`))
	})

	var out ReportEnvelope
	err := c.GetInto(context.Background(), "reports/sales-by-seller", url.Values{}, &out)
	require.NoError(t, err)
	require.Len(t, out.Data, 1)
	assert.Equal(t, "Juan Pérez", out.Data[0].SellerName)
	assert.Equal(t, Money("5000"), out.Data[0].TotalPayed)
	assert.Equal(t, Money("8000"), out.Data[0].BeforeTaxes)
}

func TestReports_ResourceAccessor(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/reports", r.URL.Path)
		_, _ = w.Write([]byte(`[{"period":"2024-02","total":2000}]`))
	})

	items, err := c.Reports().List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "2024-02", items[0].Period)
	assert.Equal(t, Money("2000"), items[0].Total)
}
